package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

func rconCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		srv ifaces.IGameServer
		reg *CommandRegistrar
		cmd *CommandRegistrant

		rcmd string
		out  string
		msg  string
		err  error
	)

	if !HasNumArgs(a, 1, -1) {
		return "", &ErrInvalidArgument{sprintf(
			`%s was passed the wrong number of arguments`, a[0])}
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	srv = reg.server
	rcmd = strings.Join(a[1:], " ")

	if out, err = srv.RunCommand(rcmd); err != nil {
		logger.LogError(cmd, sprintf("Failed to run \"%s\": %s", rcmd, err.Error()))
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
		return "", nil
	}

	_ = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	if strings.ReplaceAll(out, " ", "") != "" {
		msg = sprintf("**Output: `%s`**\n```\n%s\n```", rcmd, out)
		if utf8.RuneCountInString(out) <= 2000 {
			_, err = s.ChannelMessageSend(m.ChannelID, msg)
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "Output too large for discord")
		}
	}

	return "", err
}
