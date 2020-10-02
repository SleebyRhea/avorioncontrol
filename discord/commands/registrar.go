package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var registrars map[string]*CommandRegistrar

func init() {
	registrars = make(map[string]*CommandRegistrar)
}

// CommandRegistrar - Guild specific container for commands and authorization
// level settings
//  TODO: Move away from separate structs for role authorization. Consider moving
//  towards an int or int64 based approach
type CommandRegistrar struct {
	GuildID  string
	commands map[string]*CommandRegistrant
	loglevel int
	server   ifaces.IGameServer
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

// NewRegistrar - Create and return a new instance of CommandRegistrar
//  @gid string    ID string of the guild the CommandRegistrar belongs to
func NewRegistrar(gid string, gs ifaces.IGameServer) *CommandRegistrar {
	registrars[gid] = &CommandRegistrar{
		GuildID:  gid,
		commands: make(map[string]*CommandRegistrant, 10),
		server:   gs,
		loglevel: 1}

	return registrars[gid]
}

// Registrar - Return the Registrar that is associated with a specific guild
func Registrar(gid string) (r *CommandRegistrar, err error) {
	if r = registrars[gid]; r == nil {
		return nil, errors.New("no registrar exists for that guild")
	}
	return r, nil
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
		loglevel:    3,
		registrar:   reg}

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
	return true
}

// ProcessCommand - Processes a Discord message that has the configured prefix,
// and runs the correct command given its contents
//  @s *discordgo.Session          Discordgo Session
//  @m *discordgo.MessageCreate    Discordgo message event
//  @c IConfigurator               Bot configuration pointer
func (reg *CommandRegistrar) ProcessCommand(s *discordgo.Session,
	m *discordgo.MessageCreate, c ifaces.IConfigurator) error {
	var (
		err    error
		out    string
		cmd    *CommandRegistrant
		member *discordgo.Member
	)

	args := make(BotArgs, 0)
	input := strings.TrimPrefix(m.Content, c.Prefix())
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

	if c.CommandDisabled(cmd.Name()) {
		_, err = invalidCmd(s, m, args, c)
		return err
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
