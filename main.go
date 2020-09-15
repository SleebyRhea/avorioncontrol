package main

import (
	"AvorionControl/avorion"
	"AvorionControl/configuration"
	"AvorionControl/discord"
	"AvorionControl/ifaces"
	"AvorionControl/logger"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
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
	flag.StringVar(&token, "t", "", "Bot token")
	flag.StringVar(&prefix, "P", "", "Command prefix")
	flag.BoolVar(&showhelp, "h", false, "Help text")
	flag.IntVar(&loglevel, "l", 0, "Log level")
	flag.Parse()

	config = configuration.New()
	config.SetToken(token)
	config.SetPrefix(prefix)
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
	return c.loglevel
}

// SetLoglevel sets the loglevel of an alliance
func (c *Core) SetLoglevel(l int) {
	c.loglevel = l
}

func main() {
	if showhelp {
		os.Exit(0)
	}

	if config.Token() == "" {
		fmt.Printf("%s. %s:\n1. %s\n2. %s\n3. %s\n",
			"Please supply a token",
			"You can use one of the following methods",
			"Store the token in the environment variable [TOKEN]",
			"Use the -t command switch",
			"Supply a configuration file with said token")
		os.Exit(1)
	}

	sc := make(chan os.Signal, 1)

	core = &Core{loglevel: config.Loglevel()}
	server = avorion.New(config)
	disbot = discord.New(config)

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
			server.SendChat(ifaces.ChatData{Msg: "Shutting down server"})
			if server.IsUp() {
				if err := server.Stop(false); err != nil {
					server.SendChat(ifaces.ChatData{Msg: "Server ran into an issue shutting down"})
					log.Fatal(err)
				}
			}
			server.SendChat(ifaces.ChatData{Msg: "Server is off"})
			os.Exit(0)

		case syscall.SIGUSR1:
			logger.LogInfo(core, "Caught SIGUSR1, performing server restart")
			server.Restart()
		}
	}
}
