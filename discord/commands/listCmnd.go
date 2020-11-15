package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func listCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var (
		reg = cmd.Registrar()
		err error
		msg string
	)

	_, cmnds := reg.AllCommands()
	for _, n := range cmnds {
		cmd, _ := reg.Command(n)
		if c.CommandDisabled(cmd.Name()) {
			continue
		}
		msg = sprintf("%s\n%s - %s", msg, cmd.Name(), cmd.description)
	}

	if msg != "" {
		_, err = s.ChannelMessageSend(m.ChannelID,
			sprintf("**Available Commands:**\n```\n%s\n```", msg))
		if err != nil {
			logger.LogError(cmd, "discordgo: "+err.Error())
		}
		return "", nil
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "No commands available")
	if err != nil {
		logger.LogError(cmd, "discordgo: "+err.Error())
		return "", nil
	}

	return "", nil
}
