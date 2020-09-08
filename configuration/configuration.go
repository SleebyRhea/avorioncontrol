package configuration

import (
	"AvorionControl/logger"
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	// Discord
	defaultDiscordLink = ""
	defaultBotsAllowed = false
	defaultChatChannel = ""
	defaultLoglevel    = 1

	// Avorion
	defaultGamePort           = 27000
	defaultGamePingPort       = 27020
	defaultRconPort           = 27015
	defaultRconBin            = "/usr/bin/rcon"
	defaultRconAddress        = "127.0.0.1"
	defaultRconPassword       = "123123"
	defaultGalaxyName         = "Galaxy"
	defaultDataDirectory      = "/srv/avorion/"
	defaultServerLogDirectory = "/srv/avorion/logs"
	defaultServerInstallation = "/srv/avorion/server_files/"
	defaultTimeDatabaseUpdate = time.Minute * 60
	defaultTimeHangCheck      = time.Minute * 5
)

var sprintf = fmt.Sprintf

// Conf is a struct representing a server configuration
type Conf struct {
	BotMention func() string

	// Logging
	loglevel int

	// Avorion
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

	// Discord
	Token       string
	Prefix      string
	chatchannel string
	DiscordLink string
	BotsAllowed bool

	aliasedCommands  map[string][]string
	disabledCommands []string
}

// New returns a new object representing our program configuration
func New() *Conf {
	return &Conf{
		galaxyname: defaultGalaxyName,

		installdir: defaultServerInstallation,
		datadir:    defaultDataDirectory,
		logdir:     defaultServerLogDirectory,

		rconbin:  defaultRconBin,
		rconpass: defaultRconPassword,
		rconaddr: defaultRconAddress,

		rconport: defaultRconPort,
		gameport: defaultGamePort,
		pingport: defaultGamePingPort}
}

/************************/
/* IFace: logger.Logger */
/************************/

// UUID - Return the objects name
func (c *Conf) UUID() string {
	return "configuration"
}

// Loglevel - Return the objects name
func (c *Conf) Loglevel() int {
	return c.loglevel
}

// SetLoglevel - Return the objects name
func (c *Conf) SetLoglevel(l int) {
	c.loglevel = l
}

/********/
/* Main */
/********/

// Validate confirms that the configuration object in its current state is a
// working configuration
func (c *Conf) Validate() error {
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

// CommandDisabled - Check if a given command is disabled
//  @cmd string    Commmand to be checked
func (c *Conf) CommandDisabled(cmd string) bool {
	if len(c.disabledCommands) == 0 {
		return false
	}

	for _, dis := range c.disabledCommands {
		if dis == cmd {
			return true
		}
	}

	return false
}

// CommandAliases - Return a slice with the current aliases for a command.
//  @cmd string    Command with aliases that will be returned
func (c *Conf) CommandAliases(cmd string) (bool, []string) {
	if c.aliasedCommands[cmd] == nil {
		return false, nil
	}
	return true, c.aliasedCommands[cmd]
}

// GetAliasedCommand - Locate a command that has the given string as an alias.
//  @s string    Alias to be checked
func (c *Conf) GetAliasedCommand(s string) (bool, string) {
	for cmd, alArr := range c.aliasedCommands {
		for _, a := range alArr {
			if a == s {
				return true, cmd
			}
		}
	}

	return false, ""
}

// DisableCommand - Disable a command if it isn't already disabled.
//  @s string    Command to be disabled
func (c *Conf) DisableCommand(s string) error {
	for _, cmd := range c.disabledCommands {
		if cmd == s {
			return errors.New("Command `" + s + "` is already disabled")
		}
	}
	logger.LogInfo(c, "Disabled the command "+s)
	c.disabledCommands = append(c.disabledCommands, s)
	return nil
}

// AliasCommand - Alias a command if the alias has net yet been set.
//  @r string    Comand to be aliased
//  @a string    Alias to be configured
func (c *Conf) AliasCommand(r string, a string) error {
	for cmdname, cmd := range c.aliasedCommands {
		for _, ac := range cmd {
			if ac == a {
				return errors.New("Alias `" + a + "` is aleady aliased to `" +
					cmdname + "`")
			}
		}
	}

	if c.aliasedCommands[r] == nil {
		c.aliasedCommands[r] = make([]string, 10)
	}

	c.aliasedCommands[r] = append(c.aliasedCommands[r], a)
	logger.LogInfo(c, sprintf("Added command alias: %s -> %s", a, r))
	return nil
}

// SetChatChannel sets the channel that ifaces chat is output to
//	@id string		Channel ID to set
func (c *Conf) SetChatChannel(id string) error {
	c.chatchannel = id
	return nil
}

// ChatChannel returns the current chat channel ID string
func (c *Conf) ChatChannel() string {
	return c.chatchannel
}
