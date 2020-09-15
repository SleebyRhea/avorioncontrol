package commands

import (
	"AvorionControl/ifaces"

	"github.com/bwmarrin/discordgo"
)

func getPlayersCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		err error
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	msg := "**Tracked Players:**\n```"

	for _, p := range reg.server.Players() {
		msg = sprintf("\n%s\n%s", msg, p.Name())
	}

	msg = msg + "```"

	s.ChannelMessageSend(m.ChannelID, msg)
	return "", err
}
