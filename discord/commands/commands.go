package commands

import (
	"AvorionControl/logger"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"AvorionControl/discord/botconfig"

	"github.com/bwmarrin/discordgo"
)

var sprintf = fmt.Sprintf

/********************/
/* Utility Commands */
/********************/

// Default command used in cases where a user supplies an invalid command
func invalidCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c *botconfig.Config) (string, error) {
	msg := sprintf("The command `%s` is invalid", a[0])
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	return "", err
}

// Default command used in cases where a user does not have the authorization to
// a specific command
func unauthorizedCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c *botconfig.Config) (string, error) {
	msg := sprintf("You do not have permission to run `%s`", a[0])
	out := sprintf("Unauthorized attempt to run command: ", a[0])
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	return out, err
}

func proxySubCmnd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c *botconfig.Config) (string, error) {
	var (
		reg *CommandRegistrar
		cmd *CommandRegistrant
		msg string
		out string
		err error
	)

	if !HasNumArgs(a, 1, 1) {
		msg := sprintf("Invalid number of args passed to `%s`", a[0])
		s.ChannelMessageSend(m.ChannelID, msg)
		return out, nil
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return out, err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return out, err
	}

	_, cmdlets := cmd.Subcommands()
	for _, cmdlet := range cmdlets {
		if a[1] == cmdlet.Name() {
			return cmdlet.exec(s, m, a, c)
		}
	}

	msg = sprintf("Invalid subcommand: `%s`", a[1])
	s.ChannelMessageSend(m.ChannelID, msg)

	return out, nil
}

/************************/
/* Debug Level Commands */
/************************/

func pingCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c *botconfig.Config) (string, error) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	return "Ping request received", err
}

func pongCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c *botconfig.Config) (string, error) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Ping~")
	return "Pong request received", err
}

func loglevelCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c *botconfig.Config) (out string, err error) {
	var (
		reg    *CommandRegistrar
		cmd    *CommandRegistrant
		cmdobj *CommandRegistrant

		l  int
		ac string
		ok bool
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command("loglevel"); err != nil {
		return "", err
	}

	if !HasNumArgs(a, 2, -1) {
		out = sprintf("`%s` was not passed enough arguments", cmd.Name())
		s.ChannelMessageSend(m.ChannelID, out)
		return "", nil
	}

	// Validate our loglevel
	if l, err = strconv.Atoi(a[1]); err != nil || l > 3 || l < 0 {
		out = sprintf("`%s` is not a valid loglevel. Valid levels include:\n",
			a[1])
		out = out + "```\n0 - Only output service level info\n" +
			"1 - Show warnings (default)\n" +
			"2 - Show informational output\n" +
			"3 - Debug mode\n```"

		s.ChannelMessageSend(m.ChannelID, out)
		return "", nil
	}

	logger.LogDebug(cmd, sprintf("Using loglevel %d", l))

	for _, obj := range a[2:] {
		switch obj {
		case "registrar", "guild":
			out = sprintf("Default loglevel is now: _**%d**_", l)
			reg.SetLoglevel(l)
			s.ChannelMessageSend(m.ChannelID, out)
			continue

		// REVIEW: Is user/role specific logging something I want to implement?
		case "user", "role":
			out = "User/Role specific logging is not yet implemented"
			s.ChannelMessageSend(m.ChannelID, out)
			return

		// TODO: Implement core/config level changing
		case "core":
			out = "Changing core logging via command is not yet implemented"
			s.ChannelMessageSend(m.ChannelID, out)
			return

		// Process commands/aliases here
		// TODO: Refactor this if CommandRegistrar.GetAliasedCommand and
		// BotConfig.GetAliasedCommand are refactored. See the TODO
		// on those methods for details
		default:
			out = sprintf("Command `%s` isn't registered, nor is it an alias", obj)
			logger.LogDebug(cmd, sprintf("Checking for command %s", obj))

			if cmdobj, err = reg.Command(obj); err != nil {
				logger.LogDebug(cmd,
					sprintf("Command %s isn't registered, checking for aliases",
						obj))

				if ok, ac = c.GetAliasedCommand(obj); ok == false {
					s.ChannelMessageSend(m.ChannelID, out)
					return "", nil
				}

				if cmdobj, err = reg.Command(ac); err != nil {
					s.ChannelMessageSend(m.ChannelID, out)
					errOut := sprintf("Configured command alias is invalid: `%s`",
						ac)
					return "", errors.New(errOut)
				}
			}

			// Set that commands loglevel
			out = sprintf("Level for `%s` is now: _**%d**_", cmdobj.Name(), l)
			s.ChannelMessageSend(m.ChannelID, out)
			cmdobj.SetLoglevel(l)
			continue
		}
	}

	return "", nil
}

/************************/
/* Admin Level Commands */
/************************/

func setprefixCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c *botconfig.Config) (string, error) {
	var (
		msg string
		out string
		p   string
	)

	//eg: aa!, aa!!, !, !!, or <@!USERID> if mention is used
	r := "^([a-zA-Z0-9]{0,2}[?!;:>%$#~=+-]{1,2}|mention)$"
	author := m.Author.String()

	if !HasNumArgs(a, 1, 1) {
		msg = sprintf("Invalid number of args passed to `%s`", a[0])
		out = "User" + author + " supplied the wrong number of arguments"
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	if !regexp.MustCompile(r).MatchString(a[1]) {
		msg = sprintf("Invalid prefix supplied: `%s`", a[1])
		out = "User " + author + " attempted to set an invalid prefix"
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	if a[1] == "mention" {
		c.Prefix = "<@!" + s.State.User.ID + ">"
		msg = sprintf("Setting prefix to %s", p)
	} else {
		c.Prefix = a[1]
		msg = sprintf("Setting prefix to `%s`", a[1])
	}

	out = sprintf("Set the command prefix to %s", c.Prefix)
	_, err := s.ChannelMessageSend(m.ChannelID, msg)

	return out, err
}

func setaliasCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c *botconfig.Config) (string, error) {
	var (
		reg *CommandRegistrar
		err error
	)

	author := m.Author.String()
	out := ""
	v := "^[a-zA-Z]{1,10}$"

	if reg, err = Registrar(m.GuildID); err != nil {
		return out, err
	}

	if !HasNumArgs(a, 2, 2) {
		msg := sprintf("Invalid number of args passed to `%s`", a[0])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	if !regexp.MustCompile(v).MatchString(a[2]) {
		out = "User " + author + " attempted to set an improper alias"
		msg := sprintf("Invalid alias supplied: `%s`", a[2])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	if reg.IsRegistered(a[1]) == false {
		msg := sprintf("Command supplied is not valid: `%s`", a[1])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	if err = c.AliasCommand(a[1], a[2]); err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Failed to configure Alias!")
		return out, err
	}

	msg := sprintf("Aliased `%s` to `%s`", a[2], a[1])
	_, err = s.ChannelMessageSend(m.ChannelID, msg)

	return out, err
}

/*******************************/
/* Globally Available Commands */
/*******************************/

// FIXME: Subcommands still need to be confirmed to work or fixed, once those
// have been implemented.
func helpCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c *botconfig.Config) (string, error) {
	var (
		maincmd *CommandRegistrant //Primary command being checked
		reg     *CommandRegistrar
		err     error
		out     string
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if len(a[1:]) < 1 {
		_, err = s.ChannelMessageSend(m.ChannelID, "Please provide a command")
		return "", err
	}

	if maincmd, err = reg.Command(a[1]); err != nil {
		msg := sprintf("Command `%s` doesn't exist or isn't registered", a[1])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return "", err
	}

	if c, cmdlets := maincmd.Subcommands(); c > 0 {
		for _, cmd := range a[2:] {
		cmdletloop:
			for _, sub := range cmdlets {
				if string(cmd[0]) == sub.Name() {
					maincmd = sub
					break cmdletloop
				}
			}

			msg := sprintf("Subcommand `%s` doesn't exist under `%s`", cmd[0],
				maincmd.Name())
			_, err := s.ChannelMessageSend(m.ChannelID, msg)
			return "", err
		}
	}

	if out, err = maincmd.Help(); err != nil {
		return "", err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, out)

	return "", err
}

// InitializeCommandRegistry Register the commands for commands/everyonecmds
//	@r *CommandRegistrar    The command resistrar that we are initializing
func InitializeCommandRegistry(r *CommandRegistrar) {
	arg := newArgument

	// Global Commands
	r.Register("help",
		"Output help text for a command",
		"help <command> (subcommand) ...",
		[]CommandArgument{
			arg("command",
				"The base command that you would like help with"),
			arg("subcommand",
				"The subcommand that you would like help with")},
		helpCmd)

	// Debug Commands
	r.Register("ping",
		"Get a \"Pong!\" respons",
		"ping",
		make([]CommandArgument, 0),
		pingCmd)

	r.Register("pong",
		"Get a \"Ping~\" respons",
		"pong",
		make([]CommandArgument, 0),
		pongCmd)

	// Admin Commands
	r.Register("loglevel",
		"Set the log level for a given command(s) or object(s)",
		"loglevel <number> <object> <object> <object> ...",
		[]CommandArgument{
			arg("number",
				"The level of logging being set. Must be a number between 0 and 3"),
			arg("object",
				"Object or command that is being configured. Ex: ping, or guild")},
		loglevelCmd)

	r.Register("setprefix",
		"Set the command prefix",
		"setprefix <prefix>",
		[]CommandArgument{
			arg("prefix",
				"The prefix that is to be applied.")},
		setprefixCmd)

	r.Register("setalias",
		"Set a command alias",
		"setalias <command> <alias>",
		[]CommandArgument{
			arg("command", "Name of the command the new alias will apply to"),
			arg("alias", "Name of the alias that is being created")},
		setaliasCmd)

	// gameserver.Server (Avorion)
	r.Register("rcon",
		"Run a command in Avorion and return its result",
		"rcon <command> ...",
		[]CommandArgument{
			arg("command", "Name of the command to run"),
			arg("...", "The commands arguments")},
		rconCmd)

	r.Register("server",
		"Control the state of the Avorion server",
		"server <start|stop|restart>",
		make([]CommandArgument, 0),
		proxySubCmnd)

	r.Register("stop",
		"Stop the Avorion server (if its up)",
		"stop",
		make([]CommandArgument, 0),
		stopServerCmnd, "server")

	r.Register("start",
		"Start the Avorion server (if its down)",
		"start",
		make([]CommandArgument, 0),
		startServerCmnd, "server")

	r.Register("restart",
		"Restart the Avorion server",
		"restart",
		make([]CommandArgument, 0),
		restartServerCmnd, "server")
}
