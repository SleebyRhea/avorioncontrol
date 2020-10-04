package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func reloadConfigCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	c.LoadConfiguration()
	_ = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	msg := "Reloading bot configuration"
	if _, err := s.ChannelMessageSend(m.ChannelID, msg); err != nil {
		return "", err
	}
	return "", nil

}
