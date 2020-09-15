package configuration

import (
	"AvorionControl/ifaces"
	"AvorionControl/logger"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	// Discord
	defaultLoglevel    = 1
	defaultBotsAllowed = false
	defaultDiscordLink = "https://discord.gg/b5sqfy"

	// Avorion
	defaultGamePort           = 27000
	defaultRconPort           = 27015
	defaultGamePingPort       = 27020
	defaultRconBin            = "/usr/local/bin/rcon"
	defaultRconAddress        = "127.0.0.1"
	defaultGalaxyName         = "Galaxy"
	defaultDataDirectory      = "/srv/avorion/"
	defaultServerLogDirectory = "/srv/avorion/logs"
	defaultServerInstallation = "/srv/avorion/server_files/"
	defaultTimeDatabaseUpdate = time.Minute * 60
	defaultTimeHangCheck      = time.Minute * 5
	defaultCommandPrefix      = "mention"

	defaultTimeZone = "America/New_York"
)

var sprintf = fmt.Sprintf

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Conf is a struct representing a server configuration
type Conf struct {
	BotMention func() string

	// Logging
	loglevel int
	timezone string

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
	token       string
	prefix      string
	chatchannel string
	discordLink string
	botsallowed bool

	aliasedCommands  map[string][]string
	disabledCommands []string

	// Chat
	chatpipe chan ifaces.ChatData
}

// New returns a new object representing our program configuration
func New() *Conf {
	c := &Conf{
		galaxyname: defaultGalaxyName,

		installdir: defaultServerInstallation,
		datadir:    defaultDataDirectory,
		logdir:     defaultServerLogDirectory,

		rconbin:     defaultRconBin,
		rconpass:    makePass(),
		rconaddr:    defaultRconAddress,
		discordLink: defaultDiscordLink,

		rconport: defaultRconPort,
		gameport: defaultGamePort,
		pingport: defaultGamePingPort,

		timezone: defaultTimeZone,

		aliasedCommands: make(map[string][]string)}
	return c
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

// TimeZone -
func (c *Conf) TimeZone() string {
	return c.timezone
}

// SetTimeZone -
func (c *Conf) SetTimeZone(tz string) error {
	c.timezone = tz
	return nil
}

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

// SetAliasCommand - Alias a command if the alias has net yet been set.
//  @r string    Comand to be aliased
//  @a string    Alias to be configured
func (c *Conf) SetAliasCommand(r string, a string) error {
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

/*************************************/
/* IFace ifaces.IDiscordConfigurator */
/*************************************/

// BotsAllowed returns whether or not bots are permitted to issue commands to this
//	bot
func (c *Conf) BotsAllowed() bool {
	return c.botsallowed
}

// SetBotsAllowed sets BotsAllowed
func (c *Conf) SetBotsAllowed(allowed bool) {
	c.botsallowed = allowed
}

// DiscordLink returns a link to a Discord server
func (c *Conf) DiscordLink() string {
	return c.discordLink
}

// SetDiscordLink sets a link to a Discord server
func (c *Conf) SetDiscordLink(link string) {
	c.discordLink = link
}

// SetPrefix sets the current bot prefix
// 	TODO perform the checks for the prefix configuration *here* rather than in
// 	the *discord.Bot object
func (c *Conf) SetPrefix(prefix string) {
	c.prefix = prefix
}

// Prefix returns the current prefix
func (c *Conf) Prefix() string {
	return c.prefix
}

// SetToken sets the current Token
func (c *Conf) SetToken(t string) {
	c.token = t
}

// Token returns the current Token
func (c *Conf) Token() string {
	return c.token
}

/**********************************/
/* IFace ifaces.IGameConfigurator */
/**********************************/

// InstallPath returns the current installation path for Avorion
func (c *Conf) InstallPath() string {
	return c.installdir
}

// DataPath returns the current datapath for Avorion
func (c *Conf) DataPath() string {
	return c.datadir
}

// Galaxy returns the current Galaxyname for Avorion
func (c *Conf) Galaxy() string {
	return c.galaxyname
}

// SetGalaxy returns the current Galaxyname for Avorion
func (c *Conf) SetGalaxy(name string) {
	c.galaxyname = name
}

// RCONBin returns the current RCON binary in use
// TODO: This is temporary, until the rconlib is implemented
func (c *Conf) RCONBin() string {
	return c.rconbin
}

// RCONPort returns the current RCON port
func (c *Conf) RCONPort() int {
	return c.rconport
}

// RCONAddr returns the current RCON address
func (c *Conf) RCONAddr() string {
	return c.rconaddr
}

// RCONPass returns the current RCON password
func (c *Conf) RCONPass() string {
	return c.rconpass
}

/**********************************/
/* IFace ifaces.IChatConfigurator */
/**********************************/

// SetChatChannel sets the channel that ifaces chat is output to
//	@id string		Channel ID to set
func (c *Conf) SetChatChannel(id string) chan ifaces.ChatData {
	c.chatchannel = id

	// Close the channel if its still listening.
	if c.chatpipe != nil {
		select {
		case <-c.chatpipe:
		case <-time.After(10 * time.Nanosecond):
			logger.LogDebug(c, "Closing old chatpipe")
			close(c.chatpipe)
		}
	}

	c.chatpipe = make(chan ifaces.ChatData)
	return c.chatpipe
}

// ChatChannel returns the current chat channel ID string
func (c *Conf) ChatChannel() string {
	return c.chatchannel
}

// ChatPipe returns a go channel for chat piping
func (c *Conf) ChatPipe() chan ifaces.ChatData {
	return c.chatpipe
}
