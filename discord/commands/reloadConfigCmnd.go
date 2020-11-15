package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func reloadConfigCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	c.LoadConfiguration()
	msg := "Reloaded bot configuration"
	if _, err := s.ChannelMessageSend(m.ChannelID, msg); err != nil {
		logger.LogError(cmd, "discordgo: "+err.Error())
	}
	return sprintf("%s triggered a configuration reload", m.Author.String()), nil
}
