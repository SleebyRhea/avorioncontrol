package main

import (
	"avorioncontrol/configuration"
	"avorioncontrol/discord"
	"avorioncontrol/game/galaxy"
	"avorioncontrol/game/server/events"
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"avorioncontrol/pubsub"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
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

func main() {
	var wg sync.WaitGroup

	if showhelp {
		flag.Usage()
		os.Exit(0)
	}

	if err := config.LoadConfiguration(); err != nil {
		os.Exit(1)
	}

	if token != "" {
		config.SetToken(token)
	}

	if config.Token() == "" {
		fmt.Print("Please supply a token (see -h)\n")
		os.Exit(1)
	}

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	exit := make(chan struct{})

	gcache := galaxy.New()
	bus := pubsub.New(exit)
	core = &Core{loglevel: config.Loglevel()}
	disbot = discord.New(config, &wg, exit)

	// We start this early to prevent an errant os.Interrupt from leaving the
	// AvorionServer process running.
	signal.Notify(sc)

	disbot.Start(server, bus, gcache)
	logger.LogInit(core, "Completed init, awaiting termination signal.")
	for sig := range sc {
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			logger.LogInfo(core, "Caught termination signal. Gracefully stopping")
			close(exit)
			wg.Wait()
			config.SaveConfiguration()
			os.Exit(0)

		case syscall.SIGUSR1:
			logger.LogInfo(core, "Caught SIGUSR1, performing server reload+restart")
			config.LoadConfiguration()

		case syscall.SIGUSR2:
			logger.LogInfo(core, "Caught SIGUSR2, stopping Avorion")
			config.LoadConfiguration()
		}
	}
}

// InitializeEvents runs the event initializer
func InitializeEvents(cfg ifaces.IConfigurator, bus pubsub.MessageBus) {
	var (
		regexPlayerFID   = regexp.MustCompile(`^player:(\d+)$`)
		regexAllianceFID = regexp.MustCompile(`^alliance:(\d+)$`)
		regexSectorXY    = regexp.MustCompile(`^sector:(-?\d+:-?\d+)$`)
	)

	events.Initialize()

	for _, ed := range cfg.GetEvents() {
		ge := &events.Event{
			FString: ed.FString,
			Capture: ed.Regex,
			Handler: func(e *events.Event, in string, gc ifaces.IGalaxyCache,
				cfg ifaces.IConfigurator, sendServ, sendChat, sendLog chan interface{}) {

				logger.LogDebug(e, "Got event: "+e.FString)
				m := e.Capture.FindStringSubmatch(in)
				strings := make([]interface{}, 0)

				// Attempt to match against our player/alliance database and set that
				// string to be the name of said object
				for _, v := range m {
					switch {
					case regexPlayerFID.MatchString(v):
						v = regexPlayerFID.FindStringSubmatch(v)[1]
						if p := gc.Players().FromFactionID(v); p != nil {
							v = p.Name() + "/" + p.Steam64ID()
						}

					// TODO: Reimplement alliances
					case regexAllianceFID.MatchString(v):
						v = regexAllianceFID.FindStringSubmatch(v)[1]
						// a := s.Alliance(v)
						// if a != nil {
						// 	v = a.Name()
						// }

					// TODO: Reimplement sectors and finish this feature
					case regexSectorXY.MatchString(v):
						v = regexSectorXY.FindStringSubmatch(v)[1]
					}

					strings = append(strings, v)
				}

				sendChat <- ifaces.ChatData{Msg: fmt.Sprintf(e.FString, strings[1:]...)}
			}}

		ge.SetLoglevel(cfg.Loglevel())

		if err := events.Add(ed.Name, ge); err != nil {
			logger.LogWarning(core, "Failed to register event: "+err.Error())
			continue
		}
	}
}
