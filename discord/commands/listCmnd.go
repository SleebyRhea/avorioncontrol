package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func listCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		out = newCommandOutput(cmd, "Command Listing")
		reg = cmd.Registrar()
	)

	// Setting the description manually here, as helpCmd references this function
	// directly when there are no commands provided.
	out.Description = "List available commands for this user"
	out.Header = "Available Commands"
	out.Quoted = true

	cnt, cmnds := reg.AllCommands()
	if cnt < 1 {
		out.AddLine("No commands available")
		out.Construct()
		return out, nil
	}

	for _, n := range cmnds {
		cmd, _ := reg.Command(n)
		if c.CommandDisabled(cmd.Name()) {
			continue
		}
		out.AddLine(sprintf("**_%s_** - _%s_", cmd.Name(), cmd.description))
	}

	out.Construct()
	return out, nil
}
