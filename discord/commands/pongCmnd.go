package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func pongCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	if _, err := s.ChannelMessageSend(m.ChannelID, "Ping~"); err != nil {
		logger.LogError(cmd, "discordgo: "+err.Error())
	}
	return "Pong request received", nil
}
