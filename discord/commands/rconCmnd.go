package commands

import (
	"avorioncontrol/ifaces"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

func rconCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		srv ifaces.IGameServer
		reg *CommandRegistrar

		rcmd string
		out  string
		msg  string
		err  error
		cnt  int
	)

	if !HasNumArgs(a, 1, -1) {
		return "", &ErrInvalidArgument{sprintf(
			`%s was passed the wrong number of arguments`, a[0])}
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	srv = reg.server
	rcmd = strings.Join(a[1:], " ")

	if out, err = srv.RunCommand(rcmd); err != nil {
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
		return "", &ErrCommandError{sprintf(
			"Failed to run \"%s\": %s", rcmd, err.Error())}
	}

	_ = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if strings.ReplaceAll(out, " ", "") == "" {
		return "", nil
	}

	if utf8.RuneCountInString(out) <= 1900 {
		msg = sprintf("**Output: `%s`**\n```\n%s\n```", rcmd, out)
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return "", err
	}

	cnt = 1
	msg = ""
	for _, line := range strings.Split(out, "\n") {
		msg += line
		msg += "\n"
		if utf8.RuneCountInString(msg) >= 1900 {
			_, err = s.ChannelMessageSend(m.ChannelID, sprintf(
				"**Output %d: `%s`**\n```\n%s```\n", cnt, rcmd, msg))
			cnt++
			msg = ""
		}
	}

	if msg != "" {
		_, err = s.ChannelMessageSend(m.ChannelID, sprintf(
			"**Output %d: `%s`**\n```\n%s```\n", cnt, rcmd, msg))
	}

	return "", err
}
