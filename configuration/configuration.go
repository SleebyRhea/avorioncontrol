package configuration

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// Conf
	defaultFile = "config.yaml"

	// Discord
	defaultLoglevel    = 1
	defaultBotsAllowed = false
	defaultDiscordLink = ""

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

	// Conf
	ConfigFile string

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
		ConfigFile: defaultFile,
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

		timezone:        defaultTimeZone,
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
		c.aliasedCommands[r] = make([]string, 0)
	}

	c.aliasedCommands[r] = append(c.aliasedCommands[r], a)
	logger.LogInfo(c, sprintf("Added command alias: %s -> %s", a, r))
	return nil
}

// LoadConfiguration loads the given configuration file
func (c *Conf) LoadConfiguration() {
	if _, err := os.Stat(c.ConfigFile); err != nil {
		if os.IsNotExist(err) {
			return
		}

		fmt.Printf("Configuration file %s cannot be read (%s)",
			c.ConfigFile, err.Error())
		fmt.Printf("Proceeding with defaults\n")
		return
	}

	in, _ := ioutil.ReadFile(c.ConfigFile)
	out := &yamlData{}
	if err := yaml.Unmarshal(in, out); err != nil {
		fmt.Printf("Configuration file %s is invalid:\n%s\n", c.ConfigFile,
			err.Error())
		os.Exit(1)
	}

	//TODO: Make this not a bunch of if statements
	//TODO: Add configuration validation

	if out.Core.LogDir != "" {
		c.logdir = out.Core.LogDir
	}

	if out.Core.LogLevel != 0 {
		c.SetLoglevel(out.Core.LogLevel)
	}

	if out.Core.TimeZone != "" {
		c.SetTimeZone(out.Core.TimeZone)
	}

	if out.Game.DataDir != "" {
		c.datadir = out.Game.DataDir
	}

	if out.Game.GalaxyName != "" {
		c.galaxyname = out.Game.GalaxyName
	}

	if out.Game.GamePort != 0 {
		c.gameport = out.Game.GamePort
	}

	if out.Game.InstallDir != "" {
		c.installdir = out.Game.InstallDir
	}

	if out.Game.PingPort != 0 {
		c.pingport = out.Game.PingPort
	}

	if out.RCON.Address != "" {
		c.rconaddr = out.RCON.Address
	}

	if out.RCON.Binary != "" {
		c.rconbin = out.RCON.Binary
	}

	if out.Discord.ChatChannel != "" {
		logger.LogInfo(c, sprintf("Setting chat channel to: %s",
			out.Discord.ChatChannel))
		c.SetChatChannel(out.Discord.ChatChannel)
	}

	if len(out.Discord.AliasedCommands) > 0 {
		c.aliasedCommands = out.Discord.AliasedCommands
	}

	if len(out.Discord.DisabledCommands) > 0 {
		c.disabledCommands = out.Discord.DisabledCommands
	}

	if out.Discord.DiscordLink != "" {
		c.discordLink = out.Discord.DiscordLink
	}

	if out.Discord.Prefix != "" {
		c.SetPrefix(out.Discord.Prefix)
	}

	if out.Discord.Token != "" {
		c.SetToken(out.Discord.Token)
	}
}

// SaveConfiguration saves our current configuration to a yaml file
func (c *Conf) SaveConfiguration() {
	y := &yamlData{
		Core: yamlDataCore{
			TimeZone: c.timezone,
			LogLevel: c.loglevel,
			LogDir:   c.logdir},

		Game: yamlDataGame{
			GalaxyName: c.galaxyname,
			InstallDir: c.installdir,
			DataDir:    c.datadir,
			GamePort:   c.gameport,
			PingPort:   c.pingport},

		RCON: yamlDataRCON{
			Address: c.rconaddr,
			Binary:  c.rconbin},

		Discord: yamlDataDiscord{
			ChatChannel:      c.chatchannel,
			BotsAllowed:      c.botsallowed,
			DiscordLink:      c.discordLink,
			Prefix:           c.prefix,
			Token:            c.token,
			AliasedCommands:  c.aliasedCommands,
			DisabledCommands: c.disabledCommands}}

	if strings.HasPrefix(y.Discord.Prefix, "<@!") {
		y.Discord.Prefix = "mention"
	}

	out, err := yaml.Marshal(y)
	if err != nil {
		logger.LogError(c, err.Error())
		os.Exit(1)
	}

	if err := ioutil.WriteFile(c.ConfigFile, out, 0644); err != nil {
		logger.LogError(c, err.Error())
		os.Exit(1)
	}

	logger.LogInfo(c, "Saved configuration")
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
		case _, ok := <-c.chatpipe:
			if ok {
				close(c.chatpipe)
			}
		case <-time.After(100 * time.Nanosecond):
			logger.LogDebug(c, "Closing old chatpipe")
			close(c.chatpipe)
		}
	}

	c.chatpipe = make(chan ifaces.ChatData)
	c.SaveConfiguration()
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

/**************************************/
/* IFace ifaces.ICommandAuthenticator */
/**************************************/

// AddRoleAuth sets the role to have a given command authorization level
func (c *Conf) AddRoleAuth(rID string, l int) {
}

// RemoveRoleAuth removes authorization for a role
func (c *Conf) RemoveRoleAuth(rID string) error {
	return nil
}

// AddCmndAuth sets the authorization level required for a given command
func (c *Conf) AddCmndAuth(cmd string, l int) {
}

// GetCmndAuth gets the roles that are authorized to run the given command
func (c *Conf) GetCmndAuth(rID string, l int) {
}

// RemoveCmndAuth removes a commands authorization requirements
func (c *Conf) RemoveCmndAuth(rID string) error {
	return nil
}
