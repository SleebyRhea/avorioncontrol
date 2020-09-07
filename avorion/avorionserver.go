package avorion

import (
	"AvorionControl/discord"
	"AvorionControl/gameserver"
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

	"github.com/otiai10/copy"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	coordsre = "\\(-?[0-9]+:-?[0-9]+\\)"
	ipv4re   = "[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}"
)

var (
	reUsersOnline = regexp.MustCompile(
		"^([0-9]{18}) ([0-9]+) (" + coordsre + ") (.*) (" + ipv4re + "):[0-9]+$")

	reUsersOffline = regexp.MustCompile(
		"^([0-9]{18}) ([0-9]+) (.*)$")
)

// Server - Avorion server definition
type Server struct {
	gameserver.Server

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
	chatout chan gameserver.ChatData

	// Logger
	loglevel int
	uuid     string

	//RCON support
	rconpass string
	rconaddr string
	rconport int

	// PlayerInfo
	players  []*Player
	messages [][2]string

	// Config
	configfile string
	config     *Configuration

	// Game State
	password string
	version  string
	seed     string
	motd     string
	time     string

	// Discord
	bot      *discord.Bot
	requests map[string]string

	// Close goroutines
	close chan struct{}
	stop  chan struct{}
}

/******************/
/* avorion.Server */
/******************/

// NewServer returns a new object of type Server
func NewServer(out chan gameserver.ChatData, c *Configuration, args ...string) *Server {
	executable := "AvorionServer.exe"
	if runtime.GOOS != "windows" {
		executable = "AvorionServer"
	}

	version, err := exec.Command(c.installdir+"/bin/"+executable,
		"--version").Output()

	if err != nil {
		log.Fatal("Failed to get Avorion version: " + c.installdir + "/bin/" + executable)
		os.Exit(1)
	}

	s := &Server{
		uuid:       "AvorionServer",
		chatout:    out,
		version:    string(version),
		serverpath: c.installdir,
		executable: executable,
		config:     c,
		rconpass:   c.rconpass,
		rconaddr:   c.rconaddr,
		rconport:   c.rconport,
		requests:   make(map[string]string)}

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

// SetBot assignes a given Discord bot to the Server
func (s *Server) SetBot(b *discord.Bot) {
	s.bot = b
}

// Bot returns the currently assigned Discord bot
func (s *Server) Bot() *discord.Bot {
	return s.bot
}

/*********************/
/* gameserver.Server */
/*********************/

// Start starts the Avorion server process
func (s *Server) Start() error {
	logger.LogInfo(s, "Syncing mods to data directory")
	copy.Copy("./mods", s.config.datadir+"/mods")
	confpath := s.config.datadir +
		"mods/avocontrol-utilities/data/scripts/config/avocontrol-discord.lua"

	read, err := ioutil.ReadFile(confpath)
	if err != nil {
		return err
	}

	newconfig := strings.ReplaceAll(string(read), "%INVLINK%", s.bot.DiscordLink())
	newconfig = strings.ReplaceAll(newconfig, "%BOTNAME%", s.bot.Mention())

	err = ioutil.WriteFile(confpath, []byte(newconfig), 0)
	if err != nil {
		panic(err)
	}

	s.Cmd = exec.Command(
		s.serverpath+"/bin/"+s.executable,
		"--galaxy-name", s.config.galaxyname,
		"--datapath", s.config.datadir,
		"--admin", s.admin,
		"--rcon-ip", s.config.rconaddr,
		"--rcon-password", s.config.rconpass,
		"--rcon-port", fmt.Sprint(s.config.rconport))

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
			return nil

		case <-s.close:
			close(ready)
			return errors.New("Avorion died before initialization could complete")

		case <-time.After(5 * time.Minute):
			close(ready)
			s.Cmd.Process.Kill()
			return errors.New("Avorion took over 5 minutes to start")
		}
	}
}

// Stop gracefully stops the Avorion process
func (s *Server) Stop() error {
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
	if s.Cmd.ProcessState != nil {
		return false
	}

	if s.Cmd.Process != nil {
		return true
	}

	return false
}

// UpdatePlayerDatabase updates the Avorion player database with all of the
// players that are known to the game
func (s *Server) UpdatePlayerDatabase(notify bool) error {
	if notify {
		s.NotifyServer("Updating playerinfo DB. Possible lag incoming.")
	}

	out, err := s.RunCommand("playerinfo -o -i -s -t")
	if err != nil {
		return err
	}

	for _, info := range strings.Split(out, "\n") {
		var (
			m []string
			o bool
		)

		if m = rePlayerDataOfflineSteamIndex.FindStringSubmatch(info); m != nil {
			o = false
		} else if m = rePlayerDataOnlineSteamIndex.FindStringSubmatch(info); m != nil {
			o = true
		} else {
			logger.LogError(s, fmt.Sprintf("Unable to parse line: (%s)", info))
			continue
		}

		if p := s.Player(m[2]); p != nil {
			p.SetOnline(o)
			p.GetData()
		} else {
			p := s.NewPlayer(m[2], "")
			p.SetOnline(o)
		}
	}

	return nil
}

/************/
/* Logger */
/************/

// UUID -
func (s *Server) UUID() string {
	return s.uuid
}

// Loglevel -
func (s *Server) Loglevel() int {
	return s.loglevel
}

// SetLoglevel -
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

		ret, err := exec.Command(s.config.rconbin, "-H", s.rconaddr,
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

/*************/
/* Versioned */
/*************/

// SetVersion - Sets the current version of the Avorion server
func (s *Server) SetVersion(v string) {
	s.version = v
}

// Version - Return the version of the Avorion server
func (s *Server) Version() string {
	return s.version
}

/**********/
/* Seeded */
/**********/

// Seed - Return the current game seed
func (s *Server) Seed() string {
	return s.seed
}

// SetSeed - Sets the seed stored in the *Server, *does not* change
// the games seed
func (s *Server) SetSeed(seed string) {
	s.seed = seed
}

/********************/
/* PasswordLockable */
/********************/

// Password - Return the current password
func (s *Server) Password() string {
	return s.password
}

// SetPassword - Set the server password
func (s *Server) SetPassword(p string) {
	s.password = p
}

/*****************/
/* LoginMessager */
/*****************/

// MOTD - Return the current MOTD
func (s *Server) MOTD() string {
	return s.motd
}

// SetMOTD - Set the server MOTD
func (s *Server) SetMOTD(m string) {
	s.motd = m
}

/***************/
/* Websocketer */
/***************/

// WSOutput returns the chan that is used to output to a websocket
func (s *Server) WSOutput() chan []byte {
	return s.output
}

/***********************/
/* gameserver.Playable */
/***********************/

// Player return a player object that matches the index given
func (s *Server) Player(plrstr string) gameserver.Player {
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
func (s *Server) PlayerFromName(name string) gameserver.Player {
	for _, p := range s.players {
		logger.LogDebug(s, fmt.Sprintf("Does (%s) == (%s) ?", p.Name(), name))
		if p.Name() == name {
			return p
		}
	}
	return nil
}

// Players returns a slice of all of the  players that are currently in-game
func (s *Server) Players() []gameserver.Player {
	v := make([]gameserver.Player, 0)
	for _, t := range s.players {
		v = append(v, *t)
	}
	return v
}

// NewPlayer adds a new player to the list of players if it isn't already present
func (s *Server) NewPlayer(index, in string) gameserver.Player {
	if _, err := strconv.Atoi(index); err != nil {
		log.Fatal(errors.New("Invalid player index provided: " + index))
	}

	if p := s.Player(index); p != nil {
		p.GetData()
		return p
	}

	p := &Player{
		index:     index,
		online:    true,
		server:    s,
		steam64:   0,
		oldcoords: make([][2]int, 0)}

	if err := p.GetData(); err != nil {
		logger.LogError(s, err.Error())
	}

	s.players = append(s.players, p)
	logger.LogDebug(s, "Registering player index "+index)
	return p
}

// RemovePlayer removes a player from the list of online players
// TODO: This function is currently a stub and needs to be made functional once
// more.
func (s *Server) RemovePlayer(n string) bool {
	return false
}

// ChatMessages returns the total number of messages that are logged
func (s *Server) ChatMessages() [][2]string {
	return s.messages
}

// NewChatMessage logs the existence of a new message
func (s *Server) NewChatMessage(msg, name string) {
	s.messages = append(s.messages, [2]string{name, msg})
}

/********************************/
/* gameserver.DiscordIntegrator */
/********************************/

// DCOutput returns the chan that is used to output to Discord
func (s *Server) DCOutput() chan gameserver.ChatData {
	return s.chatout
}

// AddIntegrationRequest registers a request by a player for Discord integration
func (s *Server) AddIntegrationRequest(index, pin string) {
	s.requests[index] = pin
	path := fmt.Sprintf("%s/%s/discordrequests", s.config.datadir, s.config.galaxyname)
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
func (s *Server) ValidateIntegrationPin(in, discordID string) bool {
	m := regexp.MustCompile("^([0-9]+):([0-9]{6})$").FindStringSubmatch(in)
	if val, ok := s.requests[m[1]]; ok {
		path := fmt.Sprintf("%s/%s/discordrequests/%s",
			s.config.datadir, s.config.galaxyname, m[1])
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
	logger.LogInit(s, "Starting status supervisor")
	for {
		// Close the routine gracefully
		select {
		case <-closech:
			logger.LogInfo(s, "Stopping status supervisor")
			return

		// Check the server status every 5 minutes
		case <-time.After(time.Minute * 5):
			done := make(chan error)

			go func() {
				_, err := s.RunCommand("status")
				done <- err
			}()

			// Wait for 60 seconds and restart the server if Avorion is taking too long
			select {
			case <-time.After(60 * time.Second):
				logger.LogError(s, "Avorion is lagging, restarting")
				s.Restart()
			case <-done:
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
	scanner := bufio.NewScanner(s.stdout)
	pch := make(chan string, 0) // Player Login

	logger.LogInit(s, "Starting STDOUT supervisor")
	for scanner.Scan() {
		out := scanner.Text()

		select {
		// Exit gracefully
		case <-closech:
			logger.LogInfo(s, "Stopping STDOUT supervisor")
			return

		// Once we're ready, start processing logs.
		case <-ready:
			e := gameserver.GetEventFromString(out)

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
				logger.LogInit(s, "Avorion server initialization completed", s.WSOutput())
				close(ready) //Close the channel to close this path
			default:
				logger.LogInit(s, out, s.WSOutput())
			}
		}
	}
}
