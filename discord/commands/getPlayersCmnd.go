package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func getPlayersCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var reg = cmd.Registrar()

	if reg.server == nil || !reg.server.IsUp() {
		return "", &ErrCommandError{
			message: "Server has not finished initializing",
			cmd:     cmd}
	}

	players := reg.server.Players()
	if len(players) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No tracked players")
		return "", nil
	}

	msg := "**Tracked Players:**\n```"
	for _, p := range players {
		msg = sprintf("\n%s\n%s", msg, p.Name())
	}
	msg += "```"

	s.ChannelMessageSend(m.ChannelID, msg)
	return "", nil
}
