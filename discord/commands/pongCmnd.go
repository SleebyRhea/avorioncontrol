package commands

import (
	"AvorionControl/ifaces"

	"github.com/bwmarrin/discordgo"
)

func pongCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Ping~")
	return "Pong request received", err
}
