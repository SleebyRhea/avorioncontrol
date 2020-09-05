package main

import (
	"AvorionControl/logger"
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	reUsersOffline = regexp.MustCompile("^[0-9]{18} [0-9]+ (.*)$")

	ipv6re = "(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))"
	ipv4re = "[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}"
	ipall  = "(" + ipv6re + "|" + ipv4re + "):[0-9]+"
)

var reUsersOnline = regexp.MustCompile(
	"^[0-9]{18} [0-9]+ (\\(-?[0-9]+:-?[0-9]+\\)) (.*) " + ipall + "$")

// AvorionServer - Avorion server definition
type AvorionServer struct {
	// Execution variables
	Cmd        *exec.Cmd
	executable string
	name       string
	admin      string
	serverpath string
	datapath   string

	// IO
	stdin  io.Writer
	stdout io.Reader
	output chan []byte

	// Logger
	loglevel int
	uuid     string

	//RCON support
	rconpass string
	rconaddr string
	rconport int

	// PlayerInfo
	players  []*AvorionPlayer
	messages [][2]string

	// Config
	worldfile  string
	configfile string

	// Game State
	password string
	version  string
	seed     string
	motd     string
	time     string

	// Close goroutines
	close chan struct{}
}

// NewAvorionServer returns a new object of type AvorionServer
func NewAvorionServer(out chan []byte, c *Configuration, args ...string) *AvorionServer {
	executable := "AvorionServer.exe"
	if runtime.GOOS != "windows" {
		executable = "AvorionServer"
	}

	version, err := exec.Command(c.installdir+"/bin/"+executable,
		"--version").Output()

	if err != nil {
		log.Fatal("Failed to get Avorion version")
		os.Exit(1)
	}

	s := &AvorionServer{
		uuid:       "AvorionServer",
		output:     out,
		version:    string(version),
		serverpath: c.installdir,
		executable: executable,
		rconpass:   c.rconpass,
		rconaddr:   c.rconaddr,
		rconport:   c.rconport,
	}

	s.SetLoglevel(3)

	gameServers = append(gameServers, s)
	return s
}

// Start starts the Avorion server process
func (s *AvorionServer) Start() error {
	var err error

	s.Cmd = exec.Command(
		s.serverpath+"/bin/"+s.executable,
		"--galaxy-name", s.name,
		"--admin", s.admin,
	)

	s.Cmd.Dir = s.serverpath
	s.Cmd.Env = append(os.Environ(),
		"LD_LIBRARY_PATH="+s.serverpath+"/linux64")

	if runtime.GOOS != "windows" {
		s.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	logger.LogDebug(s, "Getting Stdin Pipe")
	if s.stdin, err = s.Cmd.StdinPipe(); err != nil {
		return err
	}

	logger.LogDebug(s, "Getting Stdout Pipe")
	if s.stdout, err = s.Cmd.StdoutPipe(); err != nil {
		return err
	}

	s.close = make(chan struct{})
	ready := make(chan struct{})

	logger.LogInit(s, "Starting supervisor goroutines")
	go superviseAvorionOut(s, ready, s.close)
	go updateAvorionStatus(s, s.close)

	logger.LogInit(s, "Starting AvorionServer and waiting till ready")
	if err = s.Cmd.Start(); err != nil {
		return err
	}

	// Wait until the process is ready and then continue
	<-ready
	logger.LogInit(s, "AvorionServer is online")
	return nil
}

// Stop gracefully stops the Avorion process
func (s *AvorionServer) Stop() error {
	if s.IsUp() != true {
		logger.LogOutput(s, "Server is already offline")
		return nil
	}

	done := make(chan error)
	s.players = nil

	logger.LogOutput(s, "Stopping Avorion server and waiting for it to exit")
	s.RunCommand("save")
	s.RunCommand("stop")

	go func() { done <- s.Cmd.Wait() }()

	select {
	case <-time.After(60 * time.Second):
		s.Cmd.Process.Kill()
		close(s.close)
		return errors.New("Avorion took too long to exit, killed")

	case err := <-done:
		logger.LogInfo(s, "Avorion server has been stopped")
		close(s.close)
		if err != nil {
			return err
		}
		return nil
	}
}

// Restart restarts the Avorion server
func (s *AvorionServer) Restart() error {
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
func (s *AvorionServer) IsUp() bool {
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
func (s *AvorionServer) UpdatePlayerDatabase() error {
	out, err := s.RunCommand("playerinfo -o -i -s")
	if err != nil {
		return err
	}

	strings.TrimSuffix(out, "\n")

	for _, playerinfo := range strings.Split(out, "\n") {
		m := reUsersOffline.FindStringSubmatch(playerinfo)
		plr := s.Player(m[2])
		if plr == nil {
		}
	}

	return nil
}

/************/
/* Logger */
/************/

// UUID -
func (s *AvorionServer) UUID() string {
	return s.uuid
}

// Loglevel -
func (s *AvorionServer) Loglevel() int {
	return s.loglevel
}

// SetLoglevel -
func (s *AvorionServer) SetLoglevel(l int) {
	s.loglevel = l
}

/***************/
/* Commandable */
/***************/

// RunCommand runs a command via rcon and returns the output
//	TODO 1: Modify this to use the games rcon websocket interface
//	TODO 2: Modify this function to make use of permitted command levels
func (s *AvorionServer) RunCommand(c string) (string, error) {
	if s.IsUp() {
		logger.LogDebug(s, "Running: "+c)

		ret, err := exec.Command("/usr/bin/rcon", "-H", s.rconaddr,
			"-p", fmt.Sprint(s.rconport), "-P", s.rconpass, c).Output()
		out := string(ret)

		if err != nil {
			return "", errors.New("Failed to run the following command: " + c)
		}

		if strings.HasPrefix(out, "Unknown command: ") {
			return out, errors.New("Invalid command provided")
		}

		return out, nil
	}

	return "", errors.New("Server is not online")
}

/*************/
/* Versioned */
/*************/

// SetVersion - Sets the current version of the Avorion server
func (s *AvorionServer) SetVersion(v string) {
	s.version = v
}

// Version - Return the version of the Avorion server
func (s *AvorionServer) Version() string {
	return s.version
}

/**********/
/* Seeded */
/**********/

// Seed - Return the current game seed
func (s *AvorionServer) Seed() string {
	return s.seed
}

// SetSeed - Sets the seed stored in the *AvorionServer, *does not* change
// the games seed
func (s *AvorionServer) SetSeed(seed string) {
	s.seed = seed
}

/********************/
/* PasswordLockable */
/********************/

// Password - Return the current password
func (s *AvorionServer) Password() string {
	return s.password
}

// SetPassword - Set the server password
func (s *AvorionServer) SetPassword(p string) {
	s.password = p
}

/*****************/
/* LoginMessager */
/*****************/

// MOTD - Return the current MOTD
func (s *AvorionServer) MOTD() string {
	return s.motd
}

// SetMOTD - Set the server MOTD
func (s *AvorionServer) SetMOTD(m string) {
	s.motd = m
}

/***************/
/* Websocketer */
/***************/

// WSOutput returns the chan that is used to output to a websocket
func (s *AvorionServer) WSOutput() chan []byte {
	return s.output
}

/********/
/* Main */
/********/

// Player return a player object that matches the string given
func (s *AvorionServer) Player(n string) Player {
	for _, p := range s.players {
		if p.Name() == n {
			return p
		}
	}

	return nil
}

// Players returns a slice of all of the  players that are currently in-game
func (s *AvorionServer) Players() []Player {
	v := make([]Player, 0)
	for _, t := range s.players {
		v = append(v, *t)
	}
	return v
}

// NewPlayer adds a new player to the list of players if it isn't already present
func (s *AvorionServer) NewPlayer(n, ips string) Player {
	if p := s.Player(n); p != nil {
		p.SetIP(ips)
		return p
	}

	plr := &AvorionPlayer{name: n, server: s}
	s.players = append(s.players, plr)

	plr.ip = net.ParseIP(ips)
	logger.LogInfo(s, "New player logged: "+plr.Name())
	return plr
}

// RemovePlayer removes a player from the list of online players
func (s *AvorionServer) RemovePlayer(n string) bool {
	for i, p := range s.players {
		if p.Name() == n {
			logger.LogInfo(s, "Removing "+p.Name())
			s.players = append(s.players[:i], s.players[i+1:]...)
			return true
		}
	}
	return false
}

// ChatMessages returns the total number of messages that are logged
func (s *AvorionServer) ChatMessages() [][2]string {
	return s.messages
}

// NewChatMessage logs the existence of a new message
func (s *AvorionServer) NewChatMessage(msg, name string) {
	s.messages = append(s.messages, [2]string{name, msg})
}

/**************/
/* Goroutines */
/**************/

// updateAvorionStatus is the goroutine responsible for making sure that the
// server is still accessible, and restarting it when needed. In addition, this
// goroutine also updates various server related data values at set intervals
func updateAvorionStatus(s *AvorionServer, closech chan struct{}) {
	for {
		// Close the routine gracefully
		select {
		case <-closech:
			break

		// Check the server status every 10 minutes
		case <-time.After(time.Minute * 5):
			time.Sleep(time.Second * 1)
			if out, err := s.RunCommand("status"); err != nil {
				s.Restart()
			} else {
				for _, o := range strings.Split(out, "\n") {
					logger.LogInfo(s, o)
				}
			}

		// Update our playerdata db
		case <-time.After(time.Hour * 3):
		}
	}
}

// superviseAvorionOut watches the output provided by the Avorion process and
// applies the applicable eventHandler for the output recieved. This routine is
// also responsible for sending the stdout of Avorion to the output channel
// to be processed by our websocket handler.
func superviseAvorionOut(s *AvorionServer, ready chan struct{},
	closech chan struct{}) {
	logger.LogDebug(s, "Started Avorion supervisor")
	scanner := bufio.NewScanner(s.stdout)

	pch := make(chan string, 0) // Player Login

	for scanner.Scan() {
		out := scanner.Text()

		select {
		// Exit gracefully
		case <-closech:
			logger.LogInfo(s, "Closed output supervision routine")
			return

		// Once we're ready, start processing logs.
		case <-ready:
			e := GetEventFromString(out)
			switch e.name {
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
