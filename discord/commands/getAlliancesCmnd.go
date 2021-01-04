package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func getAlliancesCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		reg = cmd.Registrar()
		out = newCommandOutput(cmd, "Alliances")
	)

	if reg.server == nil || !reg.server.IsUp() {
		return nil, &ErrCommandError{
			message: "Server has not finished initializing",
			cmd:     cmd}
	}

	alliances := reg.server.Alliances()
	if len(alliances) == 0 {
		out.AddLine("No tracked alliances available")
		out.Construct()
		return out, nil
	}

	for _, a := range alliances {
		out.AddLine(sprintf("**%s**: `%s`", a.Index(), a.Name()))
	}

	out.Quoted = true
	out.Construct()
	return out, nil
}
