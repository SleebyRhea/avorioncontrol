package avorion

import (
	"errors"
	"os"
	"time"
)

const (
	defaultRconBin            = "/usr/bin/rcon"
	defaultRconPort           = 27015
	defaultRconAddress        = "127.0.0.1"
	defaultRconPassword       = "123123"
	defaultServerInstallation = "/srv/avorion/server_files/"
	defaultServerLogDirectory = "/srv/avorion/logs"

	defaultTimeDatabaseUpdate = time.Hour * 1
	defaultTimeHangCheck      = time.Minute * 5
)

// Configuration is a struct representing a server configuration
type Configuration struct {
	installdir string
	logdir     string

	rconbin  string
	rconpass string
	rconaddr string
	rconport int
}

// NewConfiguration returns a new object representing a server config
func NewConfiguration() *Configuration {
	return &Configuration{
		installdir: defaultServerInstallation,
		logdir:     defaultServerLogDirectory,
		rconbin:    defaultRconBin,
		rconpass:   defaultRconPassword,
		rconaddr:   defaultRconAddress,
		rconport:   defaultRconPort,
	}
}

// Validate confirms that the configuration object in its current state is a
// working configuration
func (c *Configuration) Validate() error {
	if _, err := os.Stat(c.rconbin); err != nil {
		if os.IsNotExist(err) {
			return errors.New("RCON binary does not exist at " + c.rconbin)
		}
	}

	return nil
}
