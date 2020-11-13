package commands

import (
	"avorioncontrol/ifaces"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var sprintf = fmt.Sprintf

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
		return "", &ErrInvalidArgument{sprintf(
			`%s was passed the wrong number of arguments`, a[0])}
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
