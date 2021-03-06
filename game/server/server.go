package server

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"avorioncontrol/pubsub"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	sprintf = fmt.Sprintf
	binname = make(map[string]string, 0)
)

func init() {
	var _ ifaces.IGameServer = (*server)(nil)
	binname["windows"] = "AvorionServer.exe"
	binname["linux"] = "AvorionServer"
	binname["other"] = "AvorionServer"
}

type server struct {
	// Execution variables
	Cmd        *exec.Cmd
	dataname   string
	datapath   string
	executable string
	serverpath string

	// IO
	stdin   io.Writer
	stdout  io.Reader
	output  chan []byte
	chatout chan ifaces.ChatData

	// Logger
	loglevel int

	//RCON support
	rconexec string
	rconpass string
	rconaddr string
	rconport string

	// Game information
	password string
	version  string

	// Close goroutines
	close chan struct{}
	exit  chan struct{}
	wg    *sync.WaitGroup

	// Mutex Locks
	mutex struct {
		command *sync.Mutex
		state   *sync.Mutex
	}
}

// New starts a new server in a goroutine, and provides a function to stop it
// gracefully, and a function to issue it commands.
func New(ctx context.Context, bus pubsub.MessageBus, cfg ifaces.IConfigurator,
	wg *sync.WaitGroup, exit chan struct{}, gc ifaces.IGalaxyCache) (
	ifaces.IGameServer, error) {

	path := strings.TrimSuffix(cfg.InstallPath(), "/") + "/bin/"
	if exec, ok := binname[runtime.GOOS]; ok {
		path = path + exec
	} else {
		path = binname["other"] + exec
	}

	// Make sure that Avorion will execute
	version, err := exec.Command(path, "--version").Output()
	if err != nil {
		return nil, &ErrFailedToStart{}
	}

	// Make sure that RCON will execute
	_, err = exec.Command(cfg.RCONBin(), "-h").Output()
	if err != nil {
		return nil, &ErrRconFailedToStart{}
	}

	s := &server{
		wg:         wg,
		exit:       exit,
		version:    string(version),
		dataname:   cfg.Galaxy(),
		datapath:   cfg.DataPath(),
		loglevel:   cfg.Loglevel(),
		rconpass:   cfg.RCONPass(),
		rconaddr:   cfg.RCONAddr(),
		rconport:   fmt.Sprint(cfg.RCONPort()),
		serverpath: strings.TrimSuffix(cfg.InstallPath(), "/"),
		executable: path,
		mutex: struct {
			command *sync.Mutex
			state   *sync.Mutex
		}{
			command: &sync.Mutex{},
			state:   &sync.Mutex{},
		}}

	go s.start(ctx, cfg, gc, bus)

	return s, nil
}

// start runs the Avorion server process
func (s *server) start(ctx context.Context, cfg ifaces.IConfigurator,
	gc ifaces.IGalaxyCache, bus pubsub.MessageBus) error {

	var err error
	logger.LogInit(s, "Initializing Avorion startup sequence")

	_, endChatSub := bus.NewSubscription(pubsub.DISCCHATBUSID)     // Chat pubsub
	sendLog, endLogSub := bus.NewSubscription(pubsub.DISCLOGBUSID) // Logs pubsub

	s.Cmd = exec.Command(s.executable,
		`--galaxy-name`, s.dataname,
		`--datapath`, s.datapath,
		`--rcon-ip`, s.rconaddr,
		`--rcon-password`, s.rconpass,
		`--rcon-port`, s.rconport)

	s.Cmd.Dir = s.serverpath
	s.Cmd.Env = append(os.Environ(), "LD_LIBRARY_PATH="+
		s.serverpath+"/linux64")

	if runtime.GOOS != "windows" {
		// This prevents ctrl+c from killing the child process as well as the parent
		// on *Nix systems (not an issue on Windows). Unneeded when running as a unit.
		// https://rosettacode.org/wiki/Check_output_device_is_a_terminal#Go
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			s.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
		}
	}

	// Doing this prevents errors, but is a stub
	logger.LogDebug(s, "Getting Stdin Pipe")
	if s.stdin, err = s.Cmd.StdinPipe(); err != nil {
		return err
	}

	// Set the STDOUT pipe, so that we can reuse that as needed later
	outr, outw := io.Pipe()
	s.Cmd.Stderr = outw
	s.Cmd.Stdout = outw
	s.stdout = outr

	// Make our intercom channels
	ready := make(chan struct{})  // Avorion is fully up
	s.close = make(chan struct{}) // Close all goroutines

	// go superviseAvorionOut(s, ready, s.close)
	// go updateAvorionStatus(s, s.close)
	go func() {
		recv, end := bus.Listen(pubsub.RCONBUSID)
		defer end()

		for {
			select {
			case <-s.exit:
				return

			case in := <-recv:
				if obj, ok := in.(ifaces.RconCommand); ok {
					cmd := strings.TrimSpace(obj.Command)

					// If we get an empty command, don't send it. Instead, return an error
					// if possible. This catches the old bug wherein Avorion will crash
					// immediately upon receiving such a command (in case it returns)
					if regexp.MustCompile(`^\s*.*`).MatchString(cmd) {
						if obj.Return.Err != nil {
							obj.Return.Err <- &ErrCommandInvalid{Cmd: cmd}
							continue
						}
					}

					if len(obj.Arguments) > 0 {
						for _, arg := range obj.Arguments {
							cmd += ` "` + strings.ReplaceAll(arg, `"`, `“`) + `"`
						}
					}

					out, err := s.sendcommand(ctx, cmd)
					if obj.Return.Err != nil {
						obj.Return.Err <- err
						continue
					}

					if obj.Return.Out != nil {
						obj.Return.Out <- out
						continue
					}
				}
			}
		}
	}()

	go func() {
		defer endLogSub()
		defer endChatSub()
		defer func() {
			downstring := strings.TrimSpace(cfg.PostDownCommand())

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

				// Set the environment
				postdown := exec.CommandContext(ctx, c[0], c[1:]...)
				postdown.Env = append(os.Environ(), "SAVEPATH="+s.datapath+"/"+s.dataname)

				// Get the output of the PostDown command
				ret, err := postdown.CombinedOutput()
				if err != nil {
					logger.LogError(s, "PostDown: "+err.Error())
				}

				// Log the output
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
		}

		s.Cmd.Wait()
		logger.LogInfo(s, sprintf("Avorion exited with status code (%d)",
			s.Cmd.ProcessState.ExitCode()))
		code := s.Cmd.ProcessState.ExitCode()
		if code != 0 {
			sendLog <- ifaces.ChatData{Msg: sprintf(
				"**server Error**: Avorion has exited with non-zero status code: `%d`",
				code)}
		}

		close(s.close)
	}()

	select {
	case <-ready:
		logger.LogInit(s, "Server is online")
		cfg.LoadGameConfig()

		// If we have a Post-Up command configured, start that script in a goroutine.
		// We start it there, so that in the event that the script is intende to
		// stay online, it won't block the bot from continuing.
		if upstring := strings.TrimSpace(cfg.PostUpCommand()); upstring != "" {
			go func() {
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
					"SAVEPATH="+s.datapath+"/"+s.dataname,
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
					postup = nil
					return
				}

				defer func() {
					if postup.ProcessState == nil && postup.Process != nil {
						s.wg.Add(1)
						defer s.wg.Done()
						syscall.Kill(-postup.Process.Pid, syscall.SIGTERM)

						fin := make(chan struct{})
						logger.LogInfo(s, "Waiting for PostUp to stop")

						go func() {
							postup.Wait()
							close(fin)
						}()

						select {
						case <-fin:
							logger.LogInfo(s, "PostUp command stopped")
							return
						case <-time.After(time.Minute):
							logger.LogError(s, "Sending kill to PostUp")
							syscall.Kill(-postup.Process.Pid, syscall.SIGKILL)
							return
						}
					}
				}()

				// Stop the script when we stop the game
				select {
				case <-s.close:
					return
				case <-s.exit:
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

// UUID returns the UUID of an avorion.Server
func (s *server) UUID() string {
	return "AvorionServer"
}

// Loglevel returns the loglevel of an avorion.Server
func (s *server) Loglevel() int {
	return s.loglevel
}

// SetLoglevel sets the loglevel of an avorion.Server
func (s *server) SetLoglevel(l int) {
	s.loglevel = l
}

// Version - Return the version of the Avorion server
func (s *server) Version() string {
	return s.version
}

// Stop gracefully stops the Avorion process
func (s *server) Stop(ctx context.Context, cfg ifaces.IConfigurator,
	gc ifaces.IGalaxyCache) error {
	logger.LogDebug(s, "Stop() was called")
	if s.Online() != true {
		logger.LogOutput(s, "Server is already offline")
		return nil
	}

	logger.LogInfo(s, "Stopping Avorion server and waiting for it to exit")
	go func() {
		_, err := s.sendcommand(ctx, `save`)
		if err == nil {
			s.sendcommand(ctx, `stop`)
			return
		}
		logger.LogError(s, err.Error())
	}()

	stopt := time.After(5 * time.Minute)

	// If the process still exists after 5 minutes have passed kill the server
	// We've SIGKILL'ed the game so it *will* close, so we block until its dead
	// and writes have completed
	select {
	case <-stopt:
		s.Cmd.Process.Kill()
		<-s.close
		return errors.New("Avorion took too long to exit and had to be killed")

	// The closer channel will unblock when its closed by Avorions exit, so we can
	// use that to safely detect when this function has completed
	case <-s.close:
		logger.LogInfo(s, "Avorion server has been stopped")
		return nil
	}
}

// Online returns whether or not the game process is running
func (s *server) Online() bool {
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

// sendcommand issues a command to the Avorion server process
// 	TODO: Refactor this to use an rcon library
func (s *server) sendcommand(ctx context.Context, cmd string) (string, error) {
	if !s.Online() {
		return "", &ErrServerOffline{}
	}

	ret, err := exec.CommandContext(ctx,
		s.rconexec, "-H",
		s.rconaddr, "-p",
		s.rconport,
		"-P", s.rconpass, cmd).CombinedOutput()
	out := strings.TrimSuffix(string(ret), "\n")

	if err != nil {
		logger.LogError(s, "rcon: "+err.Error())
		logger.LogError(s, "rcon: "+out)
		return "", &ErrCommandFailedToRun{cmd: cmd, err: err}
	}

	if strings.HasPrefix(out, "Unknown command: ") {
		return out, &ErrCommandInvalid{Cmd: cmd}
	}

	return out, nil
}

// IsUp checks whether or not the game process is running
func (s *server) IsUp() bool {
	logger.LogDebug(s, "IsUp() was called")
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

// Password returns the servers password
func (s *server) Password() string {
	return s.password
}
