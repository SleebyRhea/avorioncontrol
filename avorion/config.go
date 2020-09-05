package avorion

import (
	"log"
	"os"
)

const (
	defaultRconPort           = 27015
	defaultRconAddress        = "127.0.0.1"
	defaultRconPassword       = "123123"
	defaultServerInstallation = "/srv/avorion/server_files/"
	defaultServerLogDirectory = "/srv/avorion/logs"
)

// Configuration is a struct representing a server configuration
type Configuration struct {
	installdir string
	logdir     string

	uriprefix string
	hostname  string
	rconpass  string
	rconaddr  string
	rconport  int
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
		logdir:     defaultServerLogDirectory,
		rconpass:   defaultRconPassword,
		rconaddr:   defaultRconAddress,
		rconport:   defaultRconPort,
	}
}
