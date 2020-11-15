package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

func rconCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var (
		reg = cmd.Registrar()
		srv = reg.server

		rcmd string
		out  string
		msg  string
		err  error
		cnt  int
	)

	if !HasNumArgs(a, 1, -1) {
		return "", &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`, cmd.Name()),
			cmd:     cmd}
	}

	rcmd = strings.Join(a[1:], " ")

	if out, err = srv.RunCommand(rcmd); err != nil {
		return "", &ErrCommandError{
			message: sprintf("Failed to run `%s`. Error:\n```%s```", rcmd, err.Error()),
			cmd:     cmd}
	}

	if strings.ReplaceAll(out, " ", "") == "" {
		return "", nil
	}

	if utf8.RuneCountInString(out) <= 1900 {
		msg = sprintf("**Output: `%s`**\n```\n%s\n```", rcmd, out)
		if _, err := s.ChannelMessageSend(m.ChannelID, msg); err != nil {
			logger.LogError(cmd, "discordgo: "+err.Error())
		}
		return "", nil
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
		if err != nil {
			logger.LogError(cmd, "discordgo: "+err.Error())
		}
	}

	return sprintf("%s ran the rcon command: [%s]", m.Author.String(),
		strings.Join(a[1:], " ")), nil
}
