package avorion

import (
	gamedb "avorioncontrol/avorion/database"
	"avorioncontrol/avorion/events"
	"avorioncontrol/discord"
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	// This import overrides the copy function which is undesirable
	"golang.org/x/crypto/ssh/terminal"
)

const (
	logUUID            = `AvorionServer`
	errBadIndex        = `invalid index provided (%s)`
	errExecFailed      = `failed to run Avorion binary (%s/bin/%s)`
	errBadDataString   = `failed to parse data string (%s)`
	errEmptyDataString = `got empty data string`
	errFailedRCON      = `failed to run RCON command (%s)`
	errFailToGetData   = `failed to acquire data for %s (%s)`

	warnChatDiscarded = `discarded chat message (time: >5 seconds)`
	warnGameLagging   = `Avorion is lagging, performing restart`

	noticeDBUpate       = `Updating player data DB. Potential lag incoming.`
	regexIntegration    = `^([0-9]+):([0-9]{10})$`
	rconPlayerDiscord   = `linkdiscordacct %s %s`
	rconGetPlayerData   = `getplayerdata -p %s`
	rconGetAllianceData = `getplayerdata -a %s`
	rconGetAllData      = `getplayerdata`
)

var (
	sprintf          = fmt.Sprintf
	regexpDiscordPin = regexp.MustCompile(regexIntegration)
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

	// Game Data
	players   []*Player
	alliances []*Alliance
	sectors   map[int]map[int]*ifaces.Sector
	tracking  *gamedb.TrackingDB

	// Cached values so we don't run loops constantly
	onlineplayers     string
	statusoutput      string
	onlineplayercount int
	playercount       int
	alliancecount     int
	sectorcount       int

	// Config
	configfile string
	config     ifaces.IConfigurator

	// Game State
	isrestarting bool
	isstopping   bool
	isstarting   bool
	iscrashed    bool
	password     string
	version      string
	seed         string
	motd         string
	time         string

	// Discord
	bot      *discord.Bot
	requests map[string]string

	// Close goroutines
	close chan struct{}
	stop  chan struct{}
	exit  chan struct{}
	wg    *sync.WaitGroup
}

/********/
/* Main */
/********/

// New returns a new object of type Server
func New(c ifaces.IConfigurator, wg *sync.WaitGroup, exit chan struct{},
	args ...string) ifaces.IGameServer {

	path := c.InstallPath()
	cmnd := "AvorionServer.exe"
	if runtime.GOOS != "windows" {
		cmnd = "AvorionServer"
	}

	version, err := exec.Command(path+"/bin/"+cmnd,
		"--version").Output()
	if err != nil {
		log.Fatal(sprintf(errExecFailed, path, cmnd))
	}

	_, err = exec.Command(c.RCONBin(), "-h").Output()
	if err != nil {
		log.Fatal(sprintf(`Failed to run %s`, c.RCONBin()))
	}

	s := &Server{
		wg:         wg,
		exit:       exit,
		uuid:       logUUID,
		config:     c,
		serverpath: strings.TrimSuffix(path, "/"),
		executable: cmnd,

		version:  string(version),
		rconpass: c.RCONPass(),
		rconaddr: c.RCONAddr(),
		rconport: c.RCONPort(),
		requests: make(map[string]string),
		sectors:  make(map[int]map[int]*ifaces.Sector), // =0~0=

		// State tracking
		isstarting:   false,
		isstopping:   false,
		isrestarting: false}

	s.SetLoglevel(s.config.Loglevel())
	return s
}

// NotifyServer sends an ingame notification
func (s *Server) NotifyServer(in string) error {
	cmd := sprintf("say [NOTIFICATION] %s", in)
	_, err := s.RunCommand(cmd)
	return err
}

/********************************/
/* IFace ifaces.IGameServer */
/********************************/

// Start starts the Avorion server process
func (s *Server) Start(sendchat bool) error {
	var (
		sectors []*ifaces.Sector
		err     error
	)

	s.InitializeEvents()

	defer func() { s.isstarting = false }()
	s.isstarting = true
	s.iscrashed = false

	logger.LogInit(s, "Beginning Avorion startup sequence")

	s.name = s.config.Galaxy()
	s.datapath = strings.TrimSuffix(s.config.DataPath(), "/")
	galaxydir := s.datapath + "/" + s.name

	if _, err := os.Stat(galaxydir); os.IsNotExist(err) {
		err := os.Mkdir(galaxydir, 0700)
		if err != nil {
			logger.LogError(s, "os.Mkdir: "+err.Error())
		}
	}

	if err := s.config.BuildModConfig(); err != nil {
		return errors.New("Failed to generate modconfig.lua file")
	}

	s.tracking, err = gamedb.New(sprintf("%s/%s", s.config.DataPath(),
		s.config.DBName()))
	if err != nil {
		return err
	}

	sectors, err = s.tracking.Init()
	if err != nil {
		return errors.New("GameDB: " + err.Error())
	}

	s.sectorcount = 0
	for _, sec := range sectors {
		if _, ok := s.sectors[sec.X]; !ok {
			s.sectors[sec.X] = make(map[int]*ifaces.Sector, 0)
		}
		s.sectors[sec.X][sec.Y] = sec
		s.sectorcount++
	}

	s.tracking.SetLoglevel(s.loglevel)
	logger.LogInfo(s, "Syncing mods to data directory")

	s.Cmd = exec.Command(
		s.serverpath+"/bin/"+s.executable,
		"--galaxy-name", s.name,
		"--datapath", s.datapath,
		"--admin", s.admin,
		"--rcon-ip", s.config.RCONAddr(),
		"--rcon-password", s.config.RCONPass(),
		"--rcon-port", fmt.Sprint(s.config.RCONPort()))

	s.Cmd.Dir = s.serverpath
	s.Cmd.Env = append(os.Environ(),
		"LD_LIBRARY_PATH="+s.serverpath+"/linux64")

	if runtime.GOOS != "windows" {
		// This prevents ctrl+c from killing the child process as well as the parent
		// on *Nix systems (not an issue on Windows). Unneeded when running as a unit.
		// https://rosettacode.org/wiki/Check_output_device_is_a_terminal#Go
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			s.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
		}
	}

	logger.LogDebug(s, "Getting Stdin Pipe")
	if s.stdin, err = s.Cmd.StdinPipe(); err != nil {
		return err
	}

	logger.LogDebug(s, "Getting Stdout Pipe")
	outr, outw := io.Pipe()
	s.stdout = outr
	s.Cmd.Stderr = outw
	s.Cmd.Stdout = outw

	// TODO: Determine what to do with Stderr. Either pipe it into a file, or setup
	// sometime to process it much like Stdout. Preferably keep it out of the Stdout
	// processing pipeline.
	s.Cmd.Stderr = os.Stderr
	ready := make(chan struct{})
	s.stop = make(chan struct{})
	s.close = make(chan struct{})

	go superviseAvorionOut(s, ready, s.close)
	go updateAvorionStatus(s, s.close)

	go func() {
		defer func() {
			downstring := strings.TrimSpace(s.config.PostDownCommand())

			if downstring != "" {
				c := make([]string, 0)
				// Split our arguments and add them to the args slice
				for _, m := range regexp.MustCompile(`[^\s]+`).
					FindAllStringSubmatch(downstring, -1) {
					c = append(c, m[0])
				}

				// Only allow the PostDown command to run for 1 minute
				ctx, downcancel := context.WithTimeout(context.Background(), time.Minute)
				defer downcancel()

				postdown := exec.CommandContext(ctx, c[0], c[1:]...)

				// Set the environment
				postdown.Env = append(os.Environ(),
					"SAVEPATH="+s.datapath+"/"+s.name)

				ret, err := postdown.CombinedOutput()
				if err != nil {
					logger.LogError(s, "PostDown: "+err.Error())
				}

				out := string(ret)

				if out != "" {
					for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
						logger.LogInfo(s, "PostDown: "+line)
					}
				}
			}
		}()

		logger.LogInit(s, "Starting Server and waiting till ready")
		if err := s.Cmd.Start(); err != nil {
			logger.LogError(s, err.Error())
			close(s.close)
		}

		for {
			select {
			case <-s.stop:
				s.Cmd.Wait()
				close(s.close)
				return
			case <-s.close:
				s.Cmd.Wait()
				close(s.stop)
				return
			}
		}
	}()

	for {
		select {
		case <-ready:
			logger.LogInit(s, "Server is online")
			s.config.LoadGameConfig()
			s.UpdatePlayerDatabase(false)
			s.loadSectors()

			// If we have a Post-Up command configured, start that script in a goroutine.
			// We start it there, so that in the event that the script is intende to
			// stay online, it won't block the bot from continuing.
			if upstring := strings.TrimSpace(s.config.PostUpCommand()); upstring != "" {
				go func() {
					s.wg.Add(1)
					defer s.wg.Done()

					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()

					c := make([]string, 0)
					// Split our arguments and add them to the args slice
					for _, m := range regexp.MustCompile(`[^\s]+`).
						FindAllStringSubmatch(upstring, -1) {
						c = append(c, m[0])
					}

					postup := exec.CommandContext(ctx, c[0], c[1:]...)
					postup.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
					postup.Env = append(os.Environ(),
						"SAVEPATH="+s.datapath+"/"+s.name,
						"RCONADDR="+s.rconaddr,
						"RCONPASS="+s.rconpass,
						sprintf("RCONPORT=%d", s.rconport))

					// Merge output with AvorionServer. This allows the bot to filter this
					// output along with Avorions without any extra code
					postup.Stdout = outw

					logger.LogInit(s, "Starting PostUp: "+upstring)
					if err := postup.Start(); err != nil {
						logger.LogError(s, "Failed to start configured PostUp command: "+
							upstring)
						logger.LogError(s, "PostUp: "+err.Error())
						return
					}

					defer func() {
						fin := make(chan struct{})

						go func() {
							select {
							case <-fin:
								return
							case <-time.After(time.Minute):
								syscall.Kill(-postup.Process.Pid, syscall.SIGKILL)
							}
						}()

						postup.Wait()
						close(fin)
					}()

					// Stop the script when we stop the game
					select {
					case <-s.stop:
						logger.LogInfo(s, "Stopping PostUp command")
						syscall.Kill(-postup.Process.Pid, syscall.SIGTERM)
						return

					case <-s.close:
						logger.LogInfo(s, "Stopping PostUp command")
						syscall.Kill(-postup.Process.Pid, syscall.SIGTERM)
						return
					}
				}()
			}

			return nil

		case <-s.close:
			close(ready)
			return errors.New("avorion initialization failed")

		case <-time.After(5 * time.Minute):
			close(ready)
			s.Cmd.Process.Kill()
			return errors.New("avorion took over 5 minutes to start")
		}
	}
}

// Stop gracefully stops the Avorion process
func (s *Server) Stop(sendchat bool) error {
	defer func() { s.isstopping = false }()
	s.isstopping = true

	if s.IsUp() != true {
		logger.LogOutput(s, "Server is already offline")
		return nil
	}

	logger.LogOutput(s, "Stopping Avorion server and waiting for it to exit")

	go func() {
		s.stop <- struct{}{}
		_, err := s.RunCommand("save")
		if err == nil {
			s.RunCommand("stop")
			return
		}

		logger.LogError(s, err.Error())
	}()

	s.players = nil

	select {
	case <-time.After(60 * time.Second):
		s.Cmd.Process.Kill()
		return errors.New("Avorion took too long to exit (killed)")

	case <-s.close:
		logger.LogInfo(s, "Avorion server has been stopped")
		return nil
	}
}

// Restart restarts the Avorion server
func (s *Server) Restart() error {
	if err := s.Stop(false); err != nil {
		logger.LogError(s, err.Error())
	}

	defer func() { s.isrestarting = false }()
	s.isrestarting = true

	if err := s.Start(false); err != nil {
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

// UpdatePlayerDatabase updates the Avorion player database with all
// of the players that are known to the game
//
// FIXME: Fix this absolute mess of a method
func (s *Server) UpdatePlayerDatabase(notify bool) error {
	var (
		out string
		err error
		m   []string

		allianceCount = 0
		playerCount   = 0
	)

	if notify {
		s.NotifyServer(noticeDBUpate)
	}

	if out, err = s.RunCommand(rconGetAllData); err != nil {
		logger.LogError(s, err.Error())
		return err
	}

	for _, info := range strings.Split(out, "\n") {
		switch {
		case strings.HasPrefix(info, "player: "):
			playerCount++
			if m = rePlayerData.FindStringSubmatch(info); m != nil {
				if p := s.Player(m[1]); p == nil {
					s.NewPlayer(m[1], m)
				}
			} else {
				logger.LogError(s, "player: "+sprintf(errBadDataString, info))
				continue
			}

		case strings.HasPrefix(info, "alliance: "):
			allianceCount++
			if m = reAllianceData.FindStringSubmatch(info); m != nil {
				if a := s.Alliance(m[1]); a == nil {
					s.NewAlliance(m[1], m)
				}
			} else {
				logger.LogError(s, sprintf(errBadDataString, info))
				continue
			}

		case info == "":
			logger.LogWarning(s, "playerdb: "+errEmptyDataString)

		default:
			logger.LogError(s, sprintf(errBadDataString, info))
		}
	}

	s.playercount = playerCount
	s.alliancecount = allianceCount

	for _, p := range s.players {
		s.tracking.SetDiscordToPlayer(p)
		p.SteamUID()
		logger.LogDebug(s, "Processed player: "+p.Name())
	}

	for _, a := range s.alliances {
		logger.LogDebug(s, "Processed alliance: "+a.Name())
	}

	return nil
}

// Status returns a struct containing the current status of the server
func (s *Server) Status() ifaces.ServerStatus {
	var status = ifaces.ServerOffline

	switch {
	case s.isrestarting:
		status = ifaces.ServerRestarting
	case s.isstopping:
		status = ifaces.ServerStopping
	case s.isstarting:
		status = ifaces.ServerStarting
	case s.IsUp():
		status = ifaces.ServerOnline
	}

	if s.iscrashed {
		status = ifaces.ServerCrashedOffline + status
	}

	name := s.name
	if name == "" {
		name = s.config.Galaxy()
	}

	config, _ := s.config.GameConfig()

	return ifaces.ServerStatus{
		Name:          name,
		Status:        status,
		Players:       s.onlineplayers,
		TotalPlayers:  s.playercount,
		PlayersOnline: s.onlineplayercount,
		Alliances:     s.alliancecount,
		Output:        s.statusoutput,
		Sectors:       s.sectorcount,
		INI:           config}
}

// CompareStatus takes two ifaces.ServerStatus arguments and compares
//	them. If they are equivalent, then return true. Else, false.
func (s *Server) CompareStatus(a, b ifaces.ServerStatus) bool {
	if a.Name == b.Name &&
		a.Status == b.Status &&
		a.Players == b.Players &&
		a.TotalPlayers == b.TotalPlayers &&
		a.PlayersOnline == b.PlayersOnline &&
		a.Alliances == b.Alliances &&
		a.Output == b.Output &&
		a.Sectors == b.Sectors {
		return true
	}
	return false
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

/***********************************/
/* IFace ifaces.ICommandableServer */
/***********************************/

// RunCommand runs a command via rcon and returns the output
//	TODO 1: Modify this to use the games rcon websocket interface or an rcon lib
//	TODO 2: Modify this function to make use of permitted command levels
func (s *Server) RunCommand(c string) (string, error) {
	if s.IsUp() {
		logger.LogDebug(s, "Running: "+c)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		// TODO: Make this use an rcon lib
		ret, err := exec.CommandContext(ctx, s.config.RCONBin(), "-H",
			s.rconaddr, "-p", sprintf("%d", s.rconport),
			"-P", s.rconpass, c).Output()
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

/*********************************/
/* IFace ifaces.IVersionedServer */
/*********************************/

// SetVersion - Sets the current version of the Avorion server
func (s *Server) SetVersion(v string) {
	s.version = v
}

// Version - Return the version of the Avorion server
func (s *Server) Version() string {
	return s.version
}

/******************************/
/* IFace ifaces.ISeededServer */
/******************************/

// Seed - Return the current game seed
func (s *Server) Seed() string {
	return s.seed
}

// SetSeed - Sets the seed stored in the *Server, *does not* change
// the games seed
func (s *Server) SetSeed(seed string) {
	s.seed = seed
}

/********************************/
/* IFace ifaces.ILockableServer */
/********************************/

// Password - Return the current password
func (s *Server) Password() string {
	return s.password
}

// SetPassword - Set the server password
func (s *Server) SetPassword(p string) {
	s.password = p
}

/****************************/
/* IFace ifaces.IMOTDServer */
/****************************/

// MOTD - Return the current MOTD
func (s *Server) MOTD() string {
	return s.motd
}

// SetMOTD - Set the server MOTD
func (s *Server) SetMOTD(m string) {
	s.motd = m
}

/********************************/
/* IFace ifaces.IPlayableServer */
/********************************/

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
		logger.LogDebug(s, sprintf("Does (%s) == (%s) ?", p.name, name))
		if p.name == name {
			logger.LogDebug(s, "Found player.")
			return p
		}
	}
	return nil
}

// PlayerFromDiscord return a player object that has been assigned the given
//	Discord user
//
// TODO: Complete this stub
func (s *Server) PlayerFromDiscord(name string) ifaces.IPlayer {
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

// AllianceFromName returns an alliance object that matches the name given
func (s *Server) AllianceFromName(name string) ifaces.IAlliance {
	for _, a := range s.alliances {
		logger.LogDebug(s, sprintf("Does (%s) == (%s) ?", a.name, name))
		if a.name == name {
			logger.LogDebug(s, "Found player.")
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
		logger.LogError(s, "player: "+sprintf(errBadIndex, index))
		s.Stop(true)
		os.Exit(1)
	}

	cmd := sprintf(rconGetPlayerData, index)

	if len(d) < 15 {
		if data, err := s.RunCommand(cmd); err != nil {
			logger.LogError(s, sprintf(errFailedRCON, err.Error()))
		} else {
			if d = rePlayerData.FindStringSubmatch(data); d == nil {
				logger.LogError(s, sprintf(errBadDataString, data))
				s.Stop(true)
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
	if err := s.tracking.TrackPlayer(p); err != nil {
		logger.LogError(s, err.Error())
	}
	logger.LogInfo(p, "Registered player")
	s.playercount++
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
		logger.LogError(s, "alliance: "+sprintf(errBadIndex, index))
		s.Stop(true)
		os.Exit(1)
	}

	if len(d) < 13 {
		if data, err := s.RunCommand("getplayerdata -a " + index); err != nil {
			logger.LogError(s, sprintf("Failed to get alliance data: (%s)", err.Error()))
		} else {
			if d = rePlayerData.FindStringSubmatch(data); d != nil {
				logger.LogError(s,
					sprintf("alliance: "+errBadDataString, data))
				s.Stop(true)
				<-s.close
				panic("Bad data string given in *Server.NewAlliance")
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
	s.tracking.TrackAlliance(a)
	logger.LogInfo(a, "Registered alliance")
	return a
}

// AddPlayerOnline increments the count of online players
func (s *Server) AddPlayerOnline() {
	s.onlineplayercount++
	s.updateOnlineString()
}

// SubPlayerOnline decrements the count of online players
func (s *Server) SubPlayerOnline() {
	s.onlineplayercount--
	s.updateOnlineString()
}

func (s *Server) updateOnlineString() {
	online := ""
	for _, p := range s.players {
		if p.Online() {
			online = sprintf("%s\n%s", online, p.Name())
		}
	}
	s.onlineplayers = online
	logger.LogDebug(s, "Updated online string: "+s.onlineplayers)
}

/*****************************************/
/* IFace ifaces.IDiscordIntegratedServer */
/*****************************************/

// DCOutput returns the chan that is used to output to Discord
func (s *Server) DCOutput() chan ifaces.ChatData {
	return s.chatout
}

// AddIntegrationRequest registers a request by a player for Discord integration
// TODO: Move this to our sqlite DB
func (s *Server) AddIntegrationRequest(index, pin string) {
	s.requests[index] = pin
}

// ValidateIntegrationPin confirms that a given pin was indeed a valid request
//	and registers the integration
func (s *Server) ValidateIntegrationPin(in, discordID string) bool {
	m := regexpDiscordPin.FindStringSubmatch(in)
	if len(m) < 2 {
		logger.LogError(s, sprintf("Invalid integration request provided: [%s]/[%s]",
			in, discordID))
		return false
	}

	if val, ok := s.requests[m[1]]; ok {
		if val == m[2] {
			s.tracking.AddIntegration(discordID, s.Player(m[1]))
			s.addIntegration(m[1], discordID)
			return true
		}
	}

	return false
}

/******************************/
/* IFace ifaces.IGalaxyServer */
/******************************/

// Sector returns a pointer to a sector object (new or prexisting)
func (s *Server) Sector(x, y int) *ifaces.Sector {
	// Make sure we have an X
	if _, ok := s.sectors[x]; !ok {
		s.sectors[x] = make(map[int]*ifaces.Sector, 0)
	}

	if _, ok := s.sectors[x][y]; !ok {
		s.sectors[x][y] = &ifaces.Sector{
			X: x, Y: y, Jumphistory: make([]*ifaces.JumpInfo, 0)}
		logger.LogInfo(s, sprintf("Tracking new sector: (%d:%d)", x, y))

		// TODO: This performs unnecessarily expensive DB calls here. Granted,
		// that ONLY affects initilization, but it should still be optimized
		s.tracking.TrackSector(s.sectors[x][y])
		s.sectorcount++
	}

	return s.sectors[x][y]
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
			logger.LogWarning(s, warnChatDiscarded)
		}
	}
}

// SendLog sends an ifaces.ChatData object to the discord bot if logging is
//	currently enabled in the configuration
func (s *Server) SendLog(input ifaces.ChatData) {
	if s.config.LogPipe() != nil {
		if len(input.Msg) >= 2000 {
			logger.LogInfo(s, "Truncated log for sending")
			input.Msg = input.Msg[0:1900]
			input.Msg += "...(truncated)"
		}

		select {
		case s.Config().LogPipe() <- input:
			logger.LogDebug(s, "Sent event log to bot")
		case <-time.After(time.Second * 5):
			logger.LogWarning(s, warnChatDiscarded)
		}
	}
}

// addIntegration is a helper function that registers an integration
func (s *Server) addIntegration(index, discordID string) {
	s.RunCommand(sprintf(rconPlayerDiscord, index, discordID))
}

/**************/
/* Goroutines */
/**************/

// updateAvorionStatus is the goroutine responsible for making sure that the
// server is still accessible, and restarting it when needed. In addition, this
// goroutine also updates various server related data values at set intervals
func updateAvorionStatus(s *Server, closech chan struct{}) {
	defer s.wg.Done()
	s.wg.Add(1)

	logger.LogInit(s, "Starting status supervisor")
	for {
		// Close the routine gracefully
		select {
		case <-s.exit:
			if err := s.Stop(false); err != nil {
				logger.LogError(s, err.Error())
			}
			return
		case <-closech:
			return

		// Check the server status every 5 minutes
		case <-time.After(time.Second * 30):
			if s.isrestarting || s.isstopping || s.isstarting {
				continue
			}

			if _, err := s.RunCommand("status"); err != nil {
				s.iscrashed = true
				logger.LogError(s, err.Error())
				if err := s.Restart(); err != nil {
					logger.LogError(s, err.Error())
				} else {
					s.iscrashed = false
				}
			}

		// Update our playerinfo db
		// TODO: Move the time into the configuration object
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

	logger.LogInit(s, "Started Avorion stdout supervisor")
	scanner := bufio.NewScanner(s.stdout)
	pch := make(chan string, 0) // Player Login

	// TODO: Move the scanner.Scan() loop into a goroutine.
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

			case "Server startup FAILED.":
				s.iscrashed = true

			case "Server shutdown successful.":
				logger.LogError(s, "Avorion server failed to initialize")
				close(closech)

			default:
				e := events.GetFromString(out)
				if e == nil {
					continue
				}
				e.Handler(s, e, out, nil)
			}
		}
	}
}

// TODO: Make this less godawful
func (s *Server) loadSectors() {
	for _, x := range s.sectors {
		for _, sec := range x {
			for _, j := range sec.Jumphistory {
				for _, p := range s.players {
					if p.Index() == strconv.FormatInt(int64(j.FID), 10) {
						p.jumphistory = append(p.jumphistory, ifaces.ShipCoordData{
							X: j.X, Y: j.Y, Name: j.Name, Time: j.Time})
					}
				}

				for _, a := range s.alliances {
					if a.Index() == strconv.FormatInt(int64(j.FID), 10) {
						a.jumphistory = append(a.jumphistory, ifaces.ShipCoordData{
							X: j.X, Y: j.Y, Name: j.Name, Time: j.Time})
					}
				}
			}
		}
	}

	for _, p := range s.players {
		sort.Sort(jumpsByTime(p.jumphistory))
	}

	for _, a := range s.alliances {
		sort.Sort(jumpsByTime(a.jumphistory))
	}
}

// InitializeEvents runs the event initializer
func (s *Server) InitializeEvents() {
	// Re-init our events and apply custom logged events
	events.Initialize()

	for _, ed := range s.config.GetEvents() {
		ge := &events.Event{
			FString: ed.FString,
			Capture: ed.Regex,
			Handler: func(srv ifaces.IGameServer, e *events.Event,
				in string, oc chan string) {

				logger.LogOutput(s, in)
				logger.LogDebug(e, "Got event: "+e.FString)
				m := e.Capture.FindStringSubmatch(in)
				s := make([]interface{}, 0)

				for _, v := range m {
					s = append(s, v)
				}

				srv.SendLog(ifaces.ChatData{
					Msg: sprintf(e.FString, s[1:]...)})
			}}

		ge.SetLoglevel(s.Loglevel())

		if err := events.Add(ed.Name, ge); err != nil {
			logger.LogWarning(s, "Failed to register event: "+err.Error())
			continue
		}
	}

	// Handle unmanaged text. We initilialize this last so that all other events
	// are handled first.
	events.New("EventNone", ".*", func(srv ifaces.IGameServer, e *events.Event,
		in string, oc chan string) {
		logger.LogOutput(srv, in)
	})

	logger.LogInit(s, "Completed event registration")
}

// Check if a file exists or is a directory.
func exists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
