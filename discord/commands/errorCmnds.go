package commands

import (
	"avorioncontrol/ifaces"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var sprintf = fmt.Sprintf

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
