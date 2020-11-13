package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func getPlayersCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		err error

		players = reg.server.Players()
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if len(players) < 1 {
		s.ChannelMessageSend(m.ChannelID, "No tracked players")
		return "", nil
	}

	msg := "**Tracked Players:**\n```"
	for _, p := range players {
		msg = sprintf("\n%s\n%s", msg, p.Name())
	}
	msg = msg + "```"

	s.ChannelMessageSend(m.ChannelID, msg)
	return "", err
}
