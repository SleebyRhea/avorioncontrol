package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

func statusCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var (
		reg = cmd.Registrar()
		srv = reg.server

		out string
		msg string
		err error
	)

	rcmd := "status"

	if out, err = srv.RunCommand(rcmd); err != nil {
		return "", &ErrCommandError{
			message: sprintf("Failed to run `%s`. Error:\n```%s```", rcmd, err.Error()),
			cmd:     cmd}
	}

	splitout := strings.Split(out, "\n")[0:10]
	out = strings.Join(splitout, "\n")

	if strings.ReplaceAll(out, " ", "") != "" {
		msg = sprintf("**Output: `%s`**\n```\n%s\n```", rcmd, out)
		if utf8.RuneCountInString(out) <= 2000 {
			_, err = s.ChannelMessageSend(m.ChannelID, msg)
			if err != nil {
				logger.LogError(cmd, "discordgo: "+err.Error())
			}
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "Output too large for discord")
			if err != nil {
				logger.LogError(cmd, "discordgo: "+err.Error())
			}
		}
	}

	return "", nil
}
