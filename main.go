package main

import (
	"AvorionControl/avorion"
	"AvorionControl/discord"
	"AvorionControl/discord/botconfig"
	"AvorionControl/gameserver"
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

	sc chan os.Signal
	dc chan gameserver.ChatData

	bot          *discord.Bot
	botconf      *botconfig.Config
	server       *avorion.Server
	serverconfig *avorion.Configuration
)

func init() {
	bot = &discord.Bot{}
	botconf = botconfig.New()
	serverconfig = avorion.NewConfiguration()

	flag.StringVar(&botconf.Token, "t", "", "Bot token")
	flag.StringVar(&botconf.Prefix, "P", "", "Command prefix")
	flag.BoolVar(&showhelp, "h", false, "Help text")
	flag.IntVar(&loglevel, "l", 0, "Log level")
	flag.Parse()
}

func main() {
	if botconf.Token == "" {
		fmt.Printf("%s. %s:\n1. %s\n2. %s\n3. %s\n",
			"Please supply a token",
			"You can use one of the following methods",
			"Store the token in the environment variable [TOKEN]",
			"Use the -t command switch",
			"Supply a configuration file with said token")
		os.Exit(1)
	}

	sc = make(chan os.Signal, 1)
	dc = make(chan gameserver.ChatData)

	server = avorion.NewServer(dc, serverconfig)
	server.SetBot(bot)

	if err := serverconfig.Validate(); err != nil {
		log.Fatal(err)
	}

	// We start this early to prevent an errant os.Interrupt from leaving the
	// AvorionServer process running.
	signal.Notify(sc)

	server.SetLoglevel(loglevel)
	bot.SetLoglevel(loglevel)

	discord.Init(bot, botconf, server)
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
