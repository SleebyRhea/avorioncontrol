package commands

import (
	"AvorionControl/ifaces"
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var sprintf = fmt.Sprintf

/********************/
/* Utility Commands */
/********************/

// Default command used in cases where a user supplies an invalid command
func invalidCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	msg := sprintf("The command `%s` is invalid", a[0])
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	return "", err
}

// Default command used in cases where a command doesn't have the correct amount
// of arguments passed
func wrongArgsCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		cmd *CommandRegistrant
		ok  bool
		err error
		out string
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		if ok, a[0] = c.GetAliasedCommand(a[0]); ok {
			cmd, _ = reg.Command(a[0])
		} else {
			return "", err
		}
	}

	out = sprintf("`%s` was passed the wrong number of arguments", cmd.Name())
	s.ChannelMessageSend(m.ChannelID, out)
	return "", nil
}

// Default command used in cases where a user does not have the authorization to
// a specific command
func unauthorizedCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator) (string, error) {
	msg := sprintf("You do not have permission to run `%s`", a[0])
	out := sprintf("Unauthorized attempt to run command: ", a[0])
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	return out, err
}

// Command to be used when the command being created is intended to be used with
// subcommands
func proxySubCmnd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		cmd *CommandRegistrant
		msg string
		out string
		err error
	)

	if !HasNumArgs(a, 1, 1) {
		return wrongArgsCmd(s, m, a, c)
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
/* Admin Level Commands */
/************************/

func setprefixCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator) (string, error) {
	var (
		msg string
		out string
		p   string
	)

	if !HasNumArgs(a, 1, 1) {
		return wrongArgsCmd(s, m, a, c)
	}

	//eg: aa!, aa!!, !, !!, or <@!USERID> if mention is used
	r := "^([a-zA-Z0-9]{0,2}[?!;:>%$#~=+-]{1,2}|mention)$"
	author := m.Author.String()

	if !regexp.MustCompile(r).MatchString(a[1]) {
		msg = sprintf("Invalid prefix supplied: `%s`", a[1])
		out = "User " + author + " attempted to set an invalid prefix"
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	if a[1] == "mention" {
		c.SetPrefix(sprintf("<@!%s>", s.State.User.ID))
		msg = sprintf("Setting prefix to %s", p)
	} else {
		c.SetPrefix(a[1])
		msg = sprintf("Setting prefix to `%s`", a[1])
	}

	err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return out, err
}

func setaliasCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		err error
	)

	if !HasNumArgs(a, 2, 2) {
		return wrongArgsCmd(s, m, a, c)
	}

	author := m.Author.String()
	out := ""
	v := "^[a-zA-Z]{1,10}$"

	if reg, err = Registrar(m.GuildID); err != nil {
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

	if err = c.SetAliasCommand(a[1], a[2]); err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Failed to configure Alias!")
		return out, err
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return out, err
}

/*******************************/
/* Globally Available Commands */
/*******************************/

func listCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		err error
		msg string
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	_, cmnds := reg.AllCommands()
	for _, n := range cmnds {
		cmd, _ := reg.Command(n)
		msg = sprintf("%s\n%s - %s", msg, cmd.Name(), cmd.description)
	}

	if msg != "" {
		_, err = s.ChannelMessageSend(m.ChannelID,
			sprintf("**Available Commands:**\n```\n%s\n```", msg))
		return "", err
	}

	_, err = s.ChannelMessageSend(m.ChannelID,
		sprintf("No commands available"))
	return "", err
}

// FIXME: Subcommands still need to be confirmed to work or fixed, once those
// have been implemented.
func helpCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
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

	r.Register("list",
		"Output a list of all available commands",
		"list",
		make([]CommandArgument, 0),
		listCmd)

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
		loglevelCmnd)

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

	// ifaces.Server (Avorion)
	r.Register("rcon",
		"Run a command in Avorion and return its result",
		"rcon <command> ...",
		[]CommandArgument{
			arg("command", "Name of the command to run"),
			arg("...", "The commands arguments")},
		rconCmnd)

	r.Register("setchatchannel",
		"Sets the channel to output server chat into",
		"setchatchannel channelid",
		[]CommandArgument{
			arg("channelid", "UID of the channel to send server chat messages to")},
		setChatChannelCmnd)

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
