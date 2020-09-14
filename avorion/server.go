package avorion

import (
	"AvorionControl/avorion/events"
	"AvorionControl/discord"
	"AvorionControl/ifaces"
	"AvorionControl/logger"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	// This import overrides the copy function which is undesirable
	cp "github.com/otiai10/copy"
	"golang.org/x/crypto/ssh/terminal"
)

// Server - Avorion server definition
type Server struct {
	ifaces.IGameServer

	// Execution variables
	Cmd        *exec.Cmd
	executable string
	name       string
	admin      string
	serverpath string
	datapath   string

	// IO
	stdin   io.Writer
	stdout  io.Reader
	output  chan []byte
	chatout chan ifaces.ChatData

	// Logger
	loglevel int
	uuid     string

	//RCON support
	rconpass string
	rconaddr string
	rconport int

	// PlayerInfo
	players   []*Player
	alliances []*Alliance

	// Config
	configfile string
	config     ifaces.IConfigurator

	// Game State
	password string
	version  string
	seed     string
	motd     string
	time     string

	// Discord
	bot      *discord.Bot
	requests map[string]string
	chatpipe chan ifaces.ChatData

	// Close goroutines
	close chan struct{}
	stop  chan struct{}

	isrestarting bool
	isstopping   bool
	isstarting   bool
}

/********/
/* Main */
/********/

// New returns a new object of type Server
func New(c ifaces.IConfigurator, args ...string) ifaces.IGameServer {
	executable := "AvorionServer.exe"
	if runtime.GOOS != "windows" {
		executable = "AvorionServer"
	}

	version, err := exec.Command(c.InstallPath()+"/bin/"+executable,
		"--version").Output()

	if err != nil {
		log.Fatal("Failed to get Avorion version: " + c.InstallPath() + "/bin/" + executable)
		os.Exit(1)
	}

	s := &Server{
		config:     c,
		uuid:       "AvorionServer",
		version:    string(version),
		serverpath: c.InstallPath(),
		executable: executable,
		rconpass:   c.RCONPass(),
		rconaddr:   c.RCONAddr(),
		rconport:   c.RCONPort(),
		requests:   make(map[string]string),

		// State tracking
		isstarting:   false,
		isstopping:   false,
		isrestarting: false}

	s.SetLoglevel(3)
	return s
}

// NotifyServer sends an ingame notification
func (s *Server) NotifyServer(in string) error {
	cmd := "say [NOTIFICATION] " + in
	if _, err := s.RunCommand(cmd); err != nil {
		return err
	}
	return nil
}

/********************************/
/* IFace ifaces.IGameServer */
/********************************/

// Start starts the Avorion server process
func (s *Server) Start() error {
	var err error
	defer func() { s.isstarting = false }()
	s.isstarting = true

	startupmsg := ifaces.ChatData{
		Name: "Avorion",
		Msg:  "Starting server"}

	startupdone := ifaces.ChatData{
		Name: "Avorion",
		Msg:  "Server is online"}

	startupfailed := ifaces.ChatData{
		Name: "Avorion",
		Msg:  "Server startup failed"}

	rconhost := fmt.Sprintf("[%s]\nhostname = %s\nport = %d\npassword = %s\n",
		s.config.Galaxy(), s.config.RCONAddr(), s.config.RCONPort(), s.config.RCONPass())

	ioutil.WriteFile(fmt.Sprintf("%s/rconhost.conf", s.config.DataPath()),
		[]byte(rconhost), 0644)

	s.SendChat(startupmsg)
	logger.LogInfo(s, "Syncing mods to data directory")
	cp.Copy("./mods", s.config.DataPath()+"/mods")

	s.Cmd = exec.Command(
		s.serverpath+"/bin/"+s.executable,
		"--galaxy-name", s.config.Galaxy(),
		"--datapath", s.config.DataPath(),
		"--admin", s.admin,
		"--rcon-ip", s.config.RCONAddr(),
		"--rcon-password", s.config.RCONPass(),
		"--rcon-port", fmt.Sprint(s.config.RCONPort()))

	s.Cmd.Dir = s.serverpath
	s.Cmd.Env = append(os.Environ(),
		"LD_LIBRARY_PATH="+s.serverpath+"/linux64")

	// This prevents ctrl+c from killing the child process as well as the parent
	// on *Nix systems (not an issue on Windows). Unneeded when running as a unit.
	// https://rosettacode.org/wiki/Check_output_device_is_a_terminal#Go
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		if runtime.GOOS != "windows" {
			s.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
		}
	}

	logger.LogDebug(s, "Getting Stdin Pipe")
	if s.stdin, err = s.Cmd.StdinPipe(); err != nil {
		return err
	}

	logger.LogDebug(s, "Getting Stdout Pipe")
	if s.stdout, err = s.Cmd.StdoutPipe(); err != nil {
		return err
	}

	logger.LogInit(s, "Starting Server and waiting till ready")
	ready := make(chan struct{})
	s.stop = make(chan struct{})
	s.close = make(chan struct{})

	go superviseAvorionOut(s, ready, s.close)
	go updateAvorionStatus(s, s.close)

	go func() {
		defer close(s.close)
		defer close(s.stop)

		if err := s.Cmd.Start(); err != nil {
			logger.LogError(s, err.Error())
			close(s.close)
		}

		for {
			select {
			case <-s.stop:
				s.Cmd.Wait()
				return
			}
		}
	}()

	for {
		select {
		case <-ready:
			logger.LogInit(s, "Server is online")
			s.UpdatePlayerDatabase(false)
			s.SendChat(startupdone)
			return nil

		case <-s.close:
			close(ready)
			s.SendChat(startupfailed)
			return errors.New("Avorion died before initialization could complete")

		case <-time.After(5 * time.Minute):
			close(ready)
			s.Cmd.Process.Kill()
			s.SendChat(startupfailed)
			return errors.New("Avorion took over 5 minutes to start")
		}
	}
}

// Stop gracefully stops the Avorion process
func (s *Server) Stop() error {
	defer func() { s.isstopping = false }()
	s.isstopping = true

	defer s.SendChat(ifaces.ChatData{Msg: "Server is offline", Name: "Avorion"})
	s.SendChat(ifaces.ChatData{Msg: "Stopping server", Name: "Avorion"})

	if s.IsUp() != true {
		logger.LogOutput(s, "Server is already offline")
		return nil
	}

	logger.LogOutput(s, "Stopping Avorion server and waiting for it to exit")

	s.stop <- struct{}{}
	s.RunCommand("save")
	s.RunCommand("stop")
	s.players = nil

	select {
	case <-time.After(60 * time.Second):
		s.Cmd.Process.Kill()
		return errors.New("Avorion took too long to exit, killed")

	case <-s.close:
		logger.LogInfo(s, "Avorion server has been stopped")
		return nil
	}
}

// Restart restarts the Avorion server
func (s *Server) Restart() error {
	defer func() { s.isrestarting = true }()
	s.isrestarting = true

	if err := s.Stop(); err != nil {
		return err
	}

	if err := s.Start(); err != nil {
		return err
	}

	logger.LogInfo(s, "Restarted Avorion")
	return nil
}

// IsUp checks whether or not the game process is running
func (s *Server) IsUp() bool {
	if s.Cmd == nil {
		return false
	}

	if s.Cmd.ProcessState != nil {
		return false
	}

	if s.Cmd.Process != nil {
		return true
	}

	return false
}

// Config returns the server configuration struct
func (s *Server) Config() ifaces.IConfigurator {
	return s.config
}

// UpdatePlayerDatabase updates the Avorion player database with all of the
// players that are known to the game
func (s *Server) UpdatePlayerDatabase(notify bool) error {
	if notify {
		s.NotifyServer("Updating playerinfo DB. Possible lag incoming.")
	}

	out, err := s.RunCommand("getplayerdata")
	if err != nil {
		return err
	}

	for _, info := range strings.Split(out, "\n") {
		var m []string

		if strings.HasPrefix(info, "player:") {
			if m = rePlayerData.FindStringSubmatch(info); m != nil {
				if p := s.Player(m[1]); p != nil {
					var arr [15]string
					copy(arr[:], m)
					if err := p.UpdateFromData(arr); err != nil {
						logger.LogError(s, fmt.Sprintf(
							"Failed to update player from data (%s)", err.Error()))
					}
				} else {
					s.NewPlayer(m[1], m)
				}
			} else {
				logger.LogError(s, fmt.Sprintf("Unable to parse player line: (%s)", info))
				continue
			}

		} else if strings.HasPrefix(info, "alliance:") {
			if m = reAllianceData.FindStringSubmatch(info); m != nil {
				if a := s.Alliance(m[1]); a != nil {
					var arr [13]string
					copy(arr[:], m)
					if err := a.UpdateFromData(arr); err != nil {
						logger.LogError(s, fmt.Sprintf(
							"Failed to update player from data (%s)", err.Error()))
					}
				} else {
					s.NewAlliance(m[1], m)
				}
			} else {
				logger.LogError(s, fmt.Sprintf("Unable to parse player line: (%s)", info))
				continue
			}
		}
	}

	for _, p := range s.players {
		logger.LogDebug(s, fmt.Sprintf("Processed player: %s", p.Name()))
	}

	for _, a := range s.alliances {
		logger.LogDebug(s, fmt.Sprintf("Processed alliance: %s", a.Name()))
	}

	return nil
}

/************************/
/* IFace logger.ILogger */
/************************/

// UUID returns the UUID of an avorion.Server
func (s *Server) UUID() string {
	return s.uuid
}

// Loglevel returns the loglevel of an avorion.Server
func (s *Server) Loglevel() int {
	return s.loglevel
}

// SetLoglevel sets the loglevel of an avorion.Server
func (s *Server) SetLoglevel(l int) {
	s.loglevel = l
}

/***************/
/* Commandable */
/***************/

// RunCommand runs a command via rcon and returns the output
//	TODO 1: Modify this to use the games rcon websocket interface or an rcon lib
//	TODO 2: Modify this function to make use of permitted command levels
func (s *Server) RunCommand(c string) (string, error) {
	if s.IsUp() {
		logger.LogDebug(s, "Running: "+c)

		// TODO: Make this use an rcon lib
		ret, err := exec.Command(s.config.RCONBin(), "-H", s.rconaddr,
			"-p", fmt.Sprint(s.rconport), "-P", s.rconpass, c).Output()
		out := string(ret)

		if err != nil {
			return "", errors.New("Failed to run the following command: " + c)
		}

		if strings.HasPrefix(out, "Unknown command: ") {
			return out, errors.New("Invalid command provided")
		}

		return strings.TrimSuffix(out, "\n"), nil
	}

	return "", errors.New("Server is not online")
}

/*************************************/
/* IFace ifaces.IVersionedServer */
/*************************************/

// SetVersion - Sets the current version of the Avorion server
func (s *Server) SetVersion(v string) {
	s.version = v
}

// Version - Return the version of the Avorion server
func (s *Server) Version() string {
	return s.version
}

/**********************************/
/* IFace ifaces.ISeededServer */
/**********************************/

// Seed - Return the current game seed
func (s *Server) Seed() string {
	return s.seed
}

// SetSeed - Sets the seed stored in the *Server, *does not* change
// the games seed
func (s *Server) SetSeed(seed string) {
	s.seed = seed
}

/************************************/
/* IFace ifaces.ILockableServer */
/************************************/

// Password - Return the current password
func (s *Server) Password() string {
	return s.password
}

// SetPassword - Set the server password
func (s *Server) SetPassword(p string) {
	s.password = p
}

/********************************/
/* IFace ifaces.IMOTDServer */
/********************************/

// MOTD - Return the current MOTD
func (s *Server) MOTD() string {
	return s.motd
}

// SetMOTD - Set the server MOTD
func (s *Server) SetMOTD(m string) {
	s.motd = m
}

/************************************/
/* IFace ifaces.IPlayableServer */
/************************************/

// Player return a player object that matches the index given
func (s *Server) Player(plrstr string) ifaces.IPlayer {
	// Prefer to check indexes and steamids first as those are faster to check and are more
	// common anyway
	for _, p := range s.players {
		if p.Index() == plrstr {
			return p
		}
	}
	return nil
}

// PlayerFromName return a player object that matches the name given
func (s *Server) PlayerFromName(name string) ifaces.IPlayer {
	for _, p := range s.players {
		logger.LogDebug(s, fmt.Sprintf("Does (%s) == (%s) ?", p.name, name))
		if p.name == name {
			logger.LogDebug(s, "Found player.")
			return p
		}
	}
	return nil
}

// Players returns a slice of all of the  players that are known
func (s *Server) Players() []ifaces.IPlayer {
	v := make([]ifaces.IPlayer, 0)
	for _, t := range s.players {
		v = append(v, t)
	}
	return v
}

// Alliance returns a reference to the given alliance
func (s *Server) Alliance(index string) ifaces.IAlliance {
	for _, a := range s.alliances {
		if a.Index() == index {
			return a
		}
	}
	return nil
}

// Alliances returns a slice of all of the alliances that are currently known
func (s *Server) Alliances() []ifaces.IAlliance {
	v := make([]ifaces.IAlliance, 0)
	for _, t := range s.alliances {
		v = append(v, t)
	}
	return v
}

// NewPlayer adds a new player to the list of players if it isn't already present
func (s *Server) NewPlayer(index string, d []string) ifaces.IPlayer {
	if _, err := strconv.Atoi(index); err != nil {
		log.Fatal(errors.New("Invalid player index provided: " + index))
	}

	if len(d) < 15 {
		if data, err := s.RunCommand("getplayerdata -p " + index); err != nil {
			logger.LogError(s, fmt.Sprintf("Failed to get player data: (%s)", err.Error()))
		} else {
			if d = rePlayerData.FindStringSubmatch(data); d != nil {
				logger.LogError(s,
					fmt.Sprintf("Failed to parse player string, gracefully stopping: (%s)", data))
				s.Stop()
				<-s.close
				panic("Failed to parse data string")
			}
		}
	}

	p := &Player{
		index:       index,
		name:        d[14],
		server:      s,
		jumphistory: make([]ifaces.ShipCoordData, 0),
		loglevel:    s.Loglevel()}

	// Convert our string into an array for safety
	var darr [15]string
	copy(darr[:], d)
	p.UpdateFromData(darr)
	s.players = append(s.players, p)
	logger.LogInfo(p, "Registered player")
	return p
}

// RemovePlayer removes a player from the list of online players
// TODO: This function is currently a stub and needs to be made functional once
// more.
func (s *Server) RemovePlayer(n string) {
	return
}

// NewAlliance adds a new alliance to the list of alliances if it isn't already
//	present
func (s *Server) NewAlliance(index string, d []string) ifaces.IAlliance {
	if _, err := strconv.Atoi(index); err != nil {
		log.Fatal(errors.New("Invalid alliance index provided: " + index))
	}

	if len(d) < 13 {
		if data, err := s.RunCommand("getplayerdata -a " + index); err != nil {
			logger.LogError(s, fmt.Sprintf("Failed to get alliance data: (%s)", err.Error()))
		} else {
			if d = rePlayerData.FindStringSubmatch(data); d != nil {
				logger.LogError(s,
					fmt.Sprintf("Failed to parse alliance string, gracefully stopping: (%s)", data))
				s.Stop()
				<-s.close
				panic("Failed to parse alliance data string")
			}
		}
	}

	if p := s.Alliance(index); p != nil {
		p.Update()
		return p
	}

	a := &Alliance{
		index:       index,
		name:        d[12],
		server:      s,
		jumphistory: make([]ifaces.ShipCoordData, 0),
		loglevel:    s.Loglevel()}

	s.alliances = append(s.alliances, a)
	logger.LogInfo(a, "Registered alliance")
	return a
}

/*********************************************/
/* IFace ifaces.IDiscordIntegratedServer */
/*********************************************/

// DCOutput returns the chan that is used to output to Discord
func (s *Server) DCOutput() chan ifaces.ChatData {
	return s.chatout
}

// AddIntegrationRequest registers a request by a player for Discord integration
// TODO: Move this to sqlite3
func (s *Server) AddIntegrationRequest(index, pin string) {
	s.requests[index] = pin
	path := fmt.Sprintf("%s/%s/discordrequests", s.config.DataPath(),
		s.config.Galaxy())

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0)
	}

	// For tracking when the server goes down and needs this rebuilt
	if err := ioutil.WriteFile(path+"/"+index, []byte(pin), 0); err != nil {
		logger.LogError(s, "Failed to write Discord integration request to "+
			path+"/"+index)
	}
}

// ValidateIntegrationPin confirms that a given pin was indeed a valid request
//	and registers the integration
//  TODO: Move this to sqlite3
func (s *Server) ValidateIntegrationPin(in, discordID string) bool {
	m := regexp.MustCompile("^([0-9]+):([0-9]{6})$").FindStringSubmatch(in)
	if val, ok := s.requests[m[1]]; ok {
		path := fmt.Sprintf("%s/%s/discordrequests/%s",
			s.config.DataPath(), s.config.Galaxy(), m[1])
		if val == m[2] {
			if _, err := os.Stat(path); err == nil {
				os.Remove(path)
			} else {
				logger.LogError(s, fmt.Sprintf("Failed to remove request file (%s)",
					err))
			}

			delete(s.requests, m[1])
			s.addIntegration(m[1], discordID)
			return true
		}
	}

	return false
}

// SendChat sends an ifaces.ChatData object to the discord bot if chatting is
//	currently enabled in the configuration
func (s *Server) SendChat(input ifaces.ChatData) {
	if s.config.ChatPipe() != nil {
		if len(input.Msg) >= 2000 {
			logger.LogInfo(s, "Truncated player message for sending")
			input.Msg = input.Msg[0:1900]
			input.Msg += "...(truncated)"
		}

		select {
		case s.Config().ChatPipe() <- input:
			logger.LogDebug(s, "Sent chat data to bot")
		case <-time.After(time.Second * 5):
			logger.LogWarning(s, "Failed to send chat message, took too long (discarded)")
		}
	}
}

// addIntegration is a helper function that registers an integration
func (s *Server) addIntegration(index, discordID string) {
	s.RunCommand(fmt.Sprintf("run Player(%s):setValue(\"discorduserid\", %s)",
		index, discordID))
}

/**************/
/* Goroutines */
/**************/

// updateAvorionStatus is the goroutine responsible for making sure that the
// server is still accessible, and restarting it when needed. In addition, this
// goroutine also updates various server related data values at set intervals
func updateAvorionStatus(s *Server, closech chan struct{}) {
	defer logger.LogInfo(s, "Stopped status supervisor")
	logger.LogInit(s, "Starting status supervisor")
	for {
		// Close the routine gracefully
		select {
		case <-closech:
			return

		// Check the server status every 5 minutes
		case <-time.After(time.Second * 30):
			if s.isrestarting || s.isstopping || s.isstarting {
				continue
			}

			done := make(chan error)

			go func() {
				_, err := s.RunCommand("status")
				done <- err
			}()

			// Wait for 60 seconds and restart the server if Avorion is taking too long
			select {
			case <-time.After(60 * time.Second):
				logger.LogError(s, "Avorion is lagging, restarting...")
				if err := s.Restart(); err == nil {
					s.SendChat(ifaces.ChatData{
						Msg:  "Complete.",
						Name: "Avorion"})
				} else {
					s.SendChat(ifaces.ChatData{
						Msg:  "Failed to restart Avorion",
						Name: "Avorion"})
				}

			// If rcon couldn't connect, restart.
			case err := <-done:
				if err != nil {
					logger.LogError(s, err.Error())
					s.SendChat(ifaces.ChatData{
						Msg:  "Avorion crashed, restarting...",
						Name: "Avorion"})
					if err := s.Restart(); err == nil {
						s.SendChat(ifaces.ChatData{
							Msg:  "Complete.",
							Name: "Avorion"})
					} else {
						s.SendChat(ifaces.ChatData{
							Msg:  "Failed to restart Avorion",
							Name: "Avorion"})
					}
				}

				continue
			}

		// Update our playerinfo db
		case <-time.After(1 * time.Hour):
			s.UpdatePlayerDatabase(true)
		}
	}
}

// superviseAvorionOut watches the output provided by the Avorion process and
// applies the applicable eventHandler for the output recieved. This routine is
// also responsible for sending the stdout of Avorion to the output channel
// to be processed by our websocket handler.
func superviseAvorionOut(s *Server, ready chan struct{},
	closech chan struct{}) {
	defer logger.LogInfo(s, "Stopped Avorion stdout supervisor")
	logger.LogInit(s, "Started Avorion stdout supervisor")
	scanner := bufio.NewScanner(s.stdout)
	pch := make(chan string, 0) // Player Login

	for scanner.Scan() {
		out := scanner.Text()

		select {
		// Exit gracefully
		case <-closech:
			return

		// Once we're ready, start processing logs.
		case <-ready:
			e := events.GetFromString(out)

			if e == nil {
				logger.LogOutput(s, out)
				continue
			}

			switch e.Name() {
			case "EventPlayerInfo":
				e.Handler(s, e, out, pch)
			default:
				e.Handler(s, e, out, nil)
			}

		// Output as INIT until the server is ready
		default:
			switch out {
			case "Server startup complete.":
				logger.LogInit(s, "Avorion server initialization completed")
				close(ready)

			default:
				logger.LogInit(s, out)
			}
		}
	}
}
