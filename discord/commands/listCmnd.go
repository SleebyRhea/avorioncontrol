package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func listCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		out     = newCommandOutput(cmd, "Command Listing")
		reg     = cmd.Registrar()
		authlvl = 0
		cnt     = 0
	)

	// Get the users authorization level
	member, _ := s.GuildMember(reg.GuildID, m.Author.ID)
	for _, r := range member.Roles {
		if l := c.GetRoleAuth(r); l > authlvl {
			authlvl = l
		}
	}

	// Setting the description manually here, as helpCmd references this function
	// directly when there are no commands provided.
	out.Description = "List available commands for this user"
	out.Header = "Available Commands"
	out.Quoted = true

	_, cmnds := reg.AllCommands()

	for _, n := range cmnds {
		cmd, _ := reg.Command(n)
		authreq := c.GetCmndAuth(cmd.Name())
		if c.CommandDisabled(cmd.Name()) || authlvl < authreq {
			continue
		}
		out.AddLine(sprintf("**_%s_** - _%s_", cmd.Name(), cmd.description))
		cnt++
	}

	if cnt < 1 {
		out.AddLine("No commands available")
		out.Construct()
		return out, nil
	}

	out.Construct()
	return out, nil
}
