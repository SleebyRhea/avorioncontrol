package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func reloadConfigCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var out = newCommandOutput(cmd, "Reload Configuration")
		out.Quoted = true

	c.LoadConfiguration()
	out.AddLine("Reloaded bot configuration")
	out.Construct()
	return out, nil
}
