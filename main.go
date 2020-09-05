package main

import (
	"AvorionControl/avorion"
	"AvorionControl/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		sc     = make(chan os.Signal, 1)
		config = avorion.NewConfiguration()
		server = avorion.NewServer(nil, config)
	)

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := server.Start(); err != nil {
		log.Output(1, err.Error())
		os.Exit(1)
	}

	logger.LogInit(server, "Completed INIT. Waiting for termination signal")
	signal.Notify(sc)

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
