package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func reloadConfigCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var out = newCommandOutput(cmd, "Reload Configuration")
	out.Quoted = true

	if err := c.LoadConfiguration(); err != nil {
		return nil, &ErrCommandError{
			message: err.Error(),
			cmd:     cmd}
	}

	out.AddLine("Reloaded bot configuration")

	cmd.Registrar().server.InitializeEvents()
	out.Construct()
	return out, nil
}
