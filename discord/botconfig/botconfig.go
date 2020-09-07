package botconfig

import (
	"AvorionControl/logger"
	"errors"
	"fmt"
)

var sprintf = fmt.Sprintf

// Config - Struct that describes the bots configuration backend
// TODO: Modiafy configurations to be per Guild (with a base configuration with
// global settings)
type Config struct {
	disabled []string
	aliases  map[string][]string

	Token       string
	Prefix      string
	BotsAllowed bool

	loglevel    int
	chatchannel string
}

// New - Return a pointer to a newly initializaed Config struct.
func New() *Config {
	return &Config{
		disabled:    make([]string, 10),
		aliases:     make(map[string][]string),
		BotsAllowed: false,
		chatchannel: "",
		loglevel:    1}
}

/******************************/
/* Interface: logger.Loggable */
/******************************/

// UUID - Return the objects name
func (c *Config) UUID() string {
	return "configuration"
}

// Loglevel - Return the objects name
func (c *Config) Loglevel() int {
	return c.loglevel
}

// SetLoglevel - Return the objects name
func (c *Config) SetLoglevel(l int) {
	c.loglevel = l
}

/********/
/* Main */
/********/

// IsCommandDisabled - Check if a given command is disabled
//  @cmd string    Commmand to be checked
func (c *Config) IsCommandDisabled(cmd string) bool {
	if len(c.disabled) == 0 {
		return false
	}

	for _, dis := range c.disabled {
		if dis == cmd {
			return true
		}
	}

	return false
}

// CommandAliases - Return a slice with the current aliases for a command.
//  @cmd string    Command with aliases that will be returned
func (c *Config) CommandAliases(cmd string) (bool, []string) {
	if c.aliases[cmd] == nil {
		return false, nil
	}
	return true, c.aliases[cmd]
}

// GetAliasedCommand - Locate a command that has the given string as an alias.
//  @s string    Alias to be checked
func (c *Config) GetAliasedCommand(s string) (bool, string) {
	for cmd, alArr := range c.aliases {
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
func (c *Config) DisableCommand(s string) error {
	for _, cmd := range c.disabled {
		if cmd == s {
			return errors.New("Command `" + s + "` is already disabled")
		}
	}
	logger.LogInfo(c, "Disabled the command "+s)
	c.disabled = append(c.disabled, s)
	return nil
}

// AliasCommand - Alias a command if the alias has net yet been set.
//  @r string    Comand to be aliased
//  @a string    Alias to be configured
func (c *Config) AliasCommand(r string, a string) error {
	for cmdname, cmd := range c.aliases {
		for _, ac := range cmd {
			if ac == a {
				return errors.New("Alias `" + a + "` is aleady aliased to `" +
					cmdname + "`")
			}
		}
	}

	if c.aliases[r] == nil {
		c.aliases[r] = make([]string, 10)
	}

	c.aliases[r] = append(c.aliases[r], a)
	logger.LogInfo(c, sprintf("Added command alias: %s -> %s", a, r))
	return nil
}

// SetChatChannel sets the channel that gameserver chat is output to
//	@id string		Channel ID to set
func (c *Config) SetChatChannel(id string) {
	c.chatchannel = id
}

// ChatChannel returns the current chat channel ID string
func (c *Config) ChatChannel() string {
	return c.chatchannel
}
