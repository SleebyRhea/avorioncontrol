package commands

import (
	"AvorionControl/ifaces"

	"github.com/bwmarrin/discordgo"
)

func pingCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	return "Ping request received", err
}
