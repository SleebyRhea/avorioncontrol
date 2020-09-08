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

	config *configuration.Conf
	server ifaces.IGameServer
	disbot ifaces.IDiscordBot
)

func init() {
	config = configuration.New()
	flag.StringVar(&config.Token, "t", "", "Bot token")
	flag.StringVar(&config.Prefix, "P", "", "Command prefix")
	flag.BoolVar(&showhelp, "h", false, "Help text")
	flag.IntVar(&loglevel, "l", 0, "Log level")
	flag.Parse()
}

func main() {
	if showhelp {
		os.Exit(0)
	}

	if config.Token == "" {
		fmt.Printf("%s. %s:\n1. %s\n2. %s\n3. %s\n",
			"Please supply a token",
			"You can use one of the following methods",
			"Store the token in the environment variable [TOKEN]",
			"Use the -t command switch",
			"Supply a configuration file with said token")
		os.Exit(1)
	}

	sc := make(chan os.Signal, 1)
	dc := make(chan ifaces.ChatData)

	server = avorion.New(dc, config)
	disbot = discord.New(dc, config)

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	// We start this early to prevent an errant os.Interrupt from leaving the
	// AvorionServer process running.
	signal.Notify(sc)
	disbot.Start()

	if err := server.Start(); err != nil {
		log.Output(1, err.Error())
		os.Exit(1)
	}

	logger.LogInit(server, "Completed init, awaiting termination signal.")
	for sig := range sc {
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			log.Output(1, "Caught termination signal. Gracefully stopping")
			if server.IsUp() {
				if err := server.Stop(); err != nil {
					log.Fatal(err)
				}
			}
			os.Exit(0)

		case syscall.SIGUSR1:
			log.Output(1, "Caught SIGUSR1, performing server restart")
			server.Restart()
		}
	}
}
