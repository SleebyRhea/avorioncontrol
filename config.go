package main

import (
	"log"
	"os"
)

const (
	defaultWebPort            = 8080
	defaultRconPort           = 27015
	defaultRconAddress        = "127.0.0.1"
	defaultRconPassword       = "123123"
	defaultServerInstallation = "/srv/avorion/server_files/"
)

// Configuration -
type Configuration struct {
	// ip net.IP
	// webport int

	installdir string
	uriprefix  string
	hostname   string
	rconpass   string
	rconaddr   string
	rconport   int
}

// NewConfiguration returns a new object representing a server config
func NewConfiguration() *Configuration {
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal("Failed to get server hostname")
		os.Exit(1)
	}

	return &Configuration{
		hostname:   hostname,
		installdir: defaultServerInstallation,
		rconpass:   defaultRconPassword,
		rconaddr:   defaultRconAddress,
		rconport:   defaultRconPort,
		// webport:    defaultWebPort,
	}
}

// // WebPort returns the port in string form (ex :8080)
// func (c *Configuration) WebPort() string {
// 	p := defaultWebPort
// 	if c.webport != 0 {
// 		p = c.webport
// 	}
// 	return sprintf(":%d", p)
// }

// // SetWebPort sets the configured port
// func (c *Configuration) SetWebPort(p int) error {
// 	if p > 65535 || p < 1 {
// 		return errors.New("Invalid port given")
// 	}

// 	if !isPortAvailable(p) {
// 		return errors.New("Port is already in use")
// 	}

// 	c.webport = p
// 	return nil
// }

// // isPortAvailable determines whether or not a given port is available
// func isPortAvailable(p int) bool {
// 	// https://coolaj86.com/articles/how-to-test-if-a-port-is-available-in-go/
// 	ps := strconv.Itoa(p)
// 	l, err := net.Listen("tcp", ":"+ps)
// 	if err != nil {
// 		return false
// 	}
// 	_ = l.Close()
// 	return true
// }
