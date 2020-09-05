package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	initEvents()
	initEventsDefs()
}

// Very temporary
func main() {
	var (
		sc = make(chan os.Signal, 1)
		// out    = make(chan []byte)
		config = NewConfiguration()
		server = NewAvorionServer(nil, config)
		// hub    = NewConnHub(out)
	)

	// go hub.Start()
	// serveHTTP(hub, server, out)

	// go func() {
	// 	log.Output(1, "Starting webserver")
	// 	log.Fatal(http.ListenAndServe(config.WebPort(), nil))
	// }()

	if err := server.Start(); err != nil {
		log.Output(1, err.Error())
		os.Exit(1)
	}

	log.Output(1, "Completed INIT. Waiting for termination signal")

	signal.Notify(sc)

	for sig := range sc {
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			fmt.Print("\r")
			log.Output(1, "Quitting")
			if server.IsUp() {
				if err := server.Stop(); err != nil {
					log.Fatal(err)
				}
			}
			os.Exit(0)
		default:
			log.Output(1, "Caught signal "+sig.String())
		}
	}
}
