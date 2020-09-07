package commands

import (
	"AvorionControl/discord/botconfig"
	"AvorionControl/gameserver"
	"AvorionControl/logger"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var registrars map[string]*CommandRegistrar

// BotArgs - Botarguments type (for BotCommand)
type BotArgs []string

// BotCommand - Function signature for a bots primary function
type BotCommand = func(*discordgo.Session, *discordgo.MessageCreate, BotArgs,
	*botconfig.Config) (string, error)

// CommandArgument - Define an argument for a command
//  @0    Argument's invokation
//  @1    A description of it's effect on the command
type CommandArgument [2]string

// CommandRegistrant - Command definition
type CommandRegistrant struct {
	name        string
	description string

	exec         BotCommand
	cmdlets      []*CommandRegistrant
	args         []CommandArgument
	usage        string
	hasauthlevel bool
	registrar    *CommandRegistrar

	// Implements logger.Loggable
	loglevel int
}

// SetLoglevel - Set the current loglevel
func (c *CommandRegistrant) SetLoglevel(l int) {
	c.loglevel = l
	logger.LogInfo(c, sprintf("Setting loglevel to %d", l))
}

// UUID - Return a commands name and guild information
func (c *CommandRegistrant) UUID() string {
	return c.name + "Cmd:" + c.registrar.GuildID
}

// Loglevel - Return the current loglevel for the CommandRegistrant
func (c *CommandRegistrant) Loglevel() int {
	var l int
	if c.loglevel < 0 {
		l = c.registrar.Loglevel()
	} else {
		l = c.loglevel
	}

	return l
}

// Name - Return a commands name.
func (c *CommandRegistrant) Name() string {
	return c.name
}

// Help - Return a commands help text with following markdown formatting
// preapplied
//  TODO: Convert this into an embed
//  Current Format:
//    **Command:**
//    `cmdname` - Command description
//    ```
//    Command usage
//    ```
//    **Arguments:**
//    ```
//    1: Arg1 - Description
//    2: Arg2 - Description
//    ```
//    **Subcommands:**
//    ```
//    1: sub1 - Description
//    2: sub1 - Description
//    ```
func (c *CommandRegistrant) Help() (out string, _ error) {
	out = fmt.Sprintf("**`%s`** - %s\n```\n%s\n```\n", c.name,
		c.description, c.usage)

	if len(c.args) > 0 {
		out = fmt.Sprintf("%s**Arguments:**\n```\n", out)
		for _, a := range c.args {
			out = fmt.Sprintf("%s%s - %s\n", out, a[0], a[1])
		}
		out = fmt.Sprintf("%s```\n", out)
	}

	if len(c.cmdlets) > 0 {
		out = fmt.Sprintf("%s**Subcommands:**\n```\n", out)
		for _, sc := range c.cmdlets {
			out = fmt.Sprintf("%s%s - %s\n", out, sc.Name(),
				sc.description)
		}
		out = fmt.Sprintf("%s```\n", out)
	}

	return out, nil
}

// Subcommands - Return all the count of the subcommands added to a command, and
// their slice.
func (c *CommandRegistrant) Subcommands() (int, []*CommandRegistrant) {
	if c.cmdlets != nil {
		return len(c.cmdlets), c.cmdlets
	}
	return 0, nil
}

// Registrar - Return a commands master CommandRegistrar
func (c *CommandRegistrant) Registrar() *CommandRegistrar {
	return c.registrar
}

// CommandRegistrar - Guild specific container for commands and authorization
// level settings
//  TODO: Move away from separate structs for role authorization. Consider moving
//  towards an int or int64 based approach
type CommandRegistrar struct {
	GuildID string

	commands map[string]*CommandRegistrant

	// Roles that are tied to specific levels of authorization
	adminroles []string
	modroles   []string

	// The commands that are tied to specific levels of authorization
	admincmds []string
	modcmds   []string

	// Implents logger.Loggable
	loglevel int

	// Gameserver
	server gameserver.Server
}

// SetLoglevel - Set the current loglevel
func (reg *CommandRegistrar) SetLoglevel(l int) {
	reg.loglevel = l
	logger.LogInfo(reg, sprintf("Setting loglevel to %d", l))
}

// UUID - Return a commands name.
func (reg *CommandRegistrar) UUID() string {
	return "CommandRegistrar:" + reg.GuildID
}

// Loglevel - Return the current loglevel for the CommandRegistrant
func (reg *CommandRegistrar) Loglevel() int {
	var l int
	if reg.loglevel < 0 {
		l = 0
	} else {
		l = reg.loglevel
	}

	return l
}

// Register - Register a command with the given options
//  @n string               Name
//  @d string               Description text
//  @u string               Usage text
//  @a []CommandArgument    Valid arguments
//  @f BotCommand           Function to execute
//
//  @owner ...string        Registered command to add this subcommand to
func (reg *CommandRegistrar) Register(n, d, u string, a []CommandArgument,
	f BotCommand, owners ...string) error {
	// Does it exist?
	if reg.IsRegistered(n) && len(owners) < 1 {
		log.Fatal(fmt.Errorf("command %s was already defined", n))
	}

	if len(a) == 0 {
		a = nil
	}

	registrant := &CommandRegistrant{
		name:        n,
		description: d,
		exec:        f,
		args:        a,
		usage:       u,
		registrar:   reg,
		loglevel:    3}

	// Add cmdlets to owners if strings containing command names are both
	// supplied and valid
	if len(owners) > 0 {
		for _, owner := range owners {
			if !reg.IsRegistered(owner) {
				log.Fatal("Invalid subcommand owner passed to commands.Register")
			}

			reg.commands[owner].cmdlets = append(reg.commands[owner].cmdlets,
				registrant)
			logger.LogDebug(reg, sprintf("Registered subcommand %s to %s ", n,
				reg.commands[owner].UUID()))
		}
		return nil
	}

	// Otherwise, register the command to the registrar
	reg.commands[n] = registrant

	logger.LogDebug(reg, sprintf("Registered command [%s]", n))
	return nil
}

// SetRoleAuth - Sets a given role to have a given authorization level
//  @id string    The role ID to be set
//  @lvl int      The authorization level to be used
func (reg *CommandRegistrar) SetRoleAuth(id string, lvl int) error {
	var roles *[]string

	switch lvl {
	case 0:
		roles = &reg.adminroles
	case 1:
		roles = &reg.modroles
	default:
		return errors.New("Invalid authorization level supplied: " +
			strconv.Itoa(lvl))
	}

	for _, r := range *roles {
		if r == id {
			return errors.New(
				"Role ID was already added to that authorization level")
		}
	}

	*roles = append(*roles, id)

	return nil
}

// IsRegistered - Return true if the command is registered, false if its not
//  @n string    Name of a command
func (reg *CommandRegistrar) IsRegistered(n string) bool {
	if reg.commands[n] == nil {
		return false
	}

	return true
}

// Command - Return a pointer to a given registered CommandRegistrant
//  @n string    Name of a command
//  TODO: Consider refactoring this to make use of BotConfig.GetAliasesCommand
func (reg *CommandRegistrar) Command(n string) (*CommandRegistrant, error) {
	if !reg.IsRegistered(n) {
		return nil, errors.New(sprintf("Command %s isn't registered!", n))
	}
	return reg.commands[n], nil
}

// AllCommands - Return an int and a string slice. The int is how many commands
// are currently registered, and the slice is list of their names
func (reg *CommandRegistrar) AllCommands() (_ int, cs []string) {
	for _, n := range reg.commands {
		cs = append(cs, n.Name())
	}
	return len(cs), cs
}

// UserAuthorized - Check if a user is authorized to use a given command
//  @cmd string             Name of the command being checked
//  @m *discordgo.Member    Pointer to the guild member that ran the command
func (reg *CommandRegistrar) UserAuthorized(cmd string, m *discordgo.Member) bool {
	if n := (len(reg.admincmds) + len(reg.modcmds)); n < 1 {
		return true
	}
	return true
}

// ProcessCommand - Processes a Discord message that has the configured prefix,
// and runs the correct command given its contents
//  @s *discordgo.Session          Discordgo Session
//  @m *discordgo.MessageCreate    Discordgo message event
//  @c *botconfig.Config           Bot configuration pointer
func (reg *CommandRegistrar) ProcessCommand(s *discordgo.Session,
	m *discordgo.MessageCreate, c *botconfig.Config) error {
	var (
		err    error
		out    string
		cmd    *CommandRegistrant
		member *discordgo.Member
	)

	args := make(BotArgs, 0)
	input := strings.TrimPrefix(m.Content, c.Prefix)
	input = strings.TrimSpace(input)

	// Split our arguments and add them to the args slice
	for _, m := range regexp.MustCompile("[^\\s]+").
		FindAllStringSubmatch(input, -1) {
		args = append(args, m[0])
	}

	// If there was no command given, then just return nothing
	if len(args) == 0 || len(args[0]) == 0 {
		return nil
	}

	name := args[0]

	// If the command doesn't exist, check if what was passed was an alias. If
	// it was, then we just reference the command that was aliased
	// TODO: Refactor this if CommandRegistrar.GetAliasedCommand and
	// BotConfig.GetAliasedCommand are refactored. See the TODO
	// on those methods for details
	if cmd, err = reg.Command(name); err != nil {
		if b, ac := c.GetAliasedCommand(name); b == true {
			if cmd, err = reg.Command(ac); err != nil {
				return errors.New("Configured command alias is invalid: " + ac)
			}
		} else {
			_, err = invalidCmd(s, m, args, c)
			return err
		}
	}

	if cmd.exec == nil {
		logger.LogWarning(cmd, "Can't execute (missing exec field)")
		_, err = invalidCmd(s, m, args, c)
		return err
	}

	if member, err = s.GuildMember(reg.GuildID, m.Author.ID); err != nil {
		return err
	}

	if !reg.UserAuthorized(name, member) {
		out, err = unauthorizedCmd(s, m, args, c)
		logger.LogInfo(reg, out)
		return err
	}

	// Update our arguments with the full command name
	args[0] = cmd.Name()

	// Execute our command and log any string returns
	if out, err = cmd.exec(s, m, args, c); out != "" {
		logger.LogInfo(cmd, out)
	}

	return err
}

// Registrar - Return the Registrar that is associated with a specific guild
func Registrar(gid string) (r *CommandRegistrar, err error) {
	if r = registrars[gid]; r == nil {
		return nil, errors.New("no registrar exists for that guild")
	}
	return r, nil
}

// NewRegistrar - Create and return a new instance of CommandRegistrar
//  @gid string    ID string of the guild the CommandRegistrar belongs to
func NewRegistrar(gid string, gs gameserver.Server) *CommandRegistrar {
	registrars[gid] = &CommandRegistrar{
		GuildID:    gid,
		commands:   make(map[string]*CommandRegistrant, 10),
		adminroles: make([]string, 0),
		modroles:   make([]string, 0),
		admincmds:  make([]string, 0),
		modcmds:    make([]string, 0),
		server:     gs,
		loglevel:   1}

	return registrars[gid]
}

func newArgument(a, b string) CommandArgument {
	return CommandArgument{a, b}
}

/*******************/
/* Other Functions */
/*******************/

// HasNumArgs - Determine if a set of command arguments is between min and max
//  @a BotArgs    Argument set to process
//  @min int      Minimum number of positional arguments
//  @max int      Maximum number of positional arguments
//
//  You can use -1 in place of either min or max (or both) to disable the check
//  for that range.
func HasNumArgs(a BotArgs, min, max int) bool {
	if len(a) == 0 || len(a[0]) == 0 {
		log.Fatal("Empty argument list passed to commands.HasNumArgs")
		return false
	}

	if min == -1 {
		min = 0
	}

	if max == -1 {
		max = len(a[1:]) + 1
	}

	if len(a[1:]) > max || len(a[1:]) < min {
		return false
	}

	return true
}

func init() {
	registrars = make(map[string]*CommandRegistrar)
}
