package main

import (
	"avorioncontrol/avorion"
	"avorioncontrol/configuration"
	"avorioncontrol/discord"
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	showhelp bool
	loglevel int
	token    string
	prefix   string

	config *configuration.Conf
	server ifaces.IGameServer
	disbot ifaces.IDiscordBot
	core   *Core
)

func init() {
	var configFile string
	config = configuration.New()

	flag.IntVar(&loglevel, "l", 0, "Log level")
	flag.BoolVar(&showhelp, "h", false, "Show help text")
	flag.StringVar(&token, "t", "", "Bot token")
	flag.StringVar(&configFile, "c", "", "Configuration file")
	flag.Parse()

	if configFile != "" {
		config.ConfigFile = configFile
	}

	flag.Usage = func() {
		flag.PrintDefaults()
	}
}

// Core only exists for logging purposes, and contains no other state
type Core struct {
	loglevel int
}

/************************/
/* IFace logger.ILogger */
/************************/

// UUID returns the UUID of an alliance
func (c *Core) UUID() string {
	return fmt.Sprintf("Core")
}

// Loglevel returns the loglevel of an alliance
func (c *Core) Loglevel() int {
	return config.Loglevel()
}

// SetLoglevel sets the loglevel of an alliance
func (c *Core) SetLoglevel(l int) {
	config.SetLoglevel(l)
}

func main() {
	var wg sync.WaitGroup

	if showhelp {
		flag.Usage()
		os.Exit(0)
	}

	config.LoadConfiguration()

	if token != "" {
		config.SetToken(token)
	}

	if config.Token() == "" {
		fmt.Print("Please supply a token (see -h)\n")
		os.Exit(1)
	}

	sc := make(chan os.Signal, 1)
	exit := make(chan struct{})

	core = &Core{loglevel: config.Loglevel()}
	server = avorion.New(config, &wg, exit)
	disbot = discord.New(config, &wg, exit)

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	// We start this early to prevent an errant os.Interrupt from leaving the
	// AvorionServer process running.
	signal.Notify(sc)
	disbot.Start(server)

	if err := server.Start(true); err != nil {
		logger.LogError(core, "Avorion: "+err.Error())
		os.Exit(1)
	}

	logger.LogInit(server, "Completed init, awaiting termination signal.")
	for sig := range sc {
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			logger.LogInfo(core, "Caught termination signal. Gracefully stopping")
			close(exit)
			wg.Wait()
			config.SaveConfiguration()
			os.Exit(0)

		case syscall.SIGUSR1:
			logger.LogInfo(core, "Caught SIGUSR1, performing server restart")
			server.Restart()
		}
	}
}
