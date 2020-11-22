package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func getPlayersCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		reg = cmd.Registrar()
		out = newCommandOutput(cmd, "Players")
	)

	if reg.server == nil || !reg.server.IsUp() {
		return nil, &ErrCommandError{
			message: "Server has not finished initializing",
			cmd:     cmd}
	}

	players := reg.server.Players()
	if len(players) == 0 {
		out.AddLine("No tracked players available")
		out.Construct()
		return out, nil
	}

	for _, p := range players {
		out.AddLine(sprintf("%s", p.Name()))
	}

	out.Quoted = true
	out.Construct()
	return out, nil
}
