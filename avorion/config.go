package avorion

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	defaultGamePort     = 27000
	defaultGamePingPort = 27020
	defaultRconPort     = 27015

	defaultRconBin      = "/usr/bin/rcon"
	defaultRconAddress  = "127.0.0.1"
	defaultRconPassword = "123123"

	defaultGalaxyName         = "Galaxy"
	defaultDataDirectory      = "/srv/avorion/"
	defaultServerLogDirectory = "/srv/avorion/logs"
	defaultServerInstallation = "/srv/avorion/server_files/"

	defaultTimeDatabaseUpdate = time.Minute * 60
	defaultTimeHangCheck      = time.Minute * 5
)

// Configuration is a struct representing a server configuration
type Configuration struct {
	galaxyname string

	installdir string
	datadir    string
	logdir     string

	rconbin  string
	rconpass string
	rconaddr string

	rconport int
	gameport int
	pingport int
}

// NewConfiguration returns a new object representing a server config
func NewConfiguration() *Configuration {
	return &Configuration{
		galaxyname: defaultGalaxyName,

		installdir: defaultServerInstallation,
		datadir:    defaultDataDirectory,
		logdir:     defaultServerLogDirectory,

		rconbin:  defaultRconBin,
		rconpass: defaultRconPassword,
		rconaddr: defaultRconAddress,

		rconport: defaultRconPort,
		gameport: defaultGamePort,
		pingport: defaultGamePingPort,
	}
}

// Validate confirms that the configuration object in its current state is a
// working configuration
func (c *Configuration) Validate() error {
	ports := []int{c.gameport, c.rconport, c.pingport}

	if _, err := os.Stat(c.rconbin); err != nil {
		if os.IsNotExist(err) {
			return errors.New("RCON binary does not exist at " + c.rconbin)
		}
	}

	for _, port := range ports {
		if !isPortAvailable(port) {
			return fmt.Errorf("Port %d is not available", port)
		}
	}

	return nil
}

func isPortAvailable(p int) bool {
	// https://coolaj86.com/articles/how-to-test-if-a-port-is-available-in-go/
	ps := strconv.Itoa(p)
	l, err := net.Listen("tcp", ":"+ps)
	if err != nil {
		return false
	}
	_ = l.Close()
	return true
}
