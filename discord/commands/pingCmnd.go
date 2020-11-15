package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func pingCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	if _, err := s.ChannelMessageSend(m.ChannelID, "Pong!"); err != nil {
		logger.LogError(cmd, "discordgo: "+err.Error())
	}
	return "Ping request received", nil
}
