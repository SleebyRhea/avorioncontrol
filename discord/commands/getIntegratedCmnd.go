package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func getIntegratedCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {

	var (
		reg = cmd.Registrar()
		srv = reg.server
		out = newCommandOutput(cmd, "Players with Discord Integration")
		cnt = 0
	)

	out.Header = "Players Found"

	for _, p := range srv.Players() {
		if p.DiscordUID() != "" {
			logger.LogDebug(cmd, sprintf("Processing player %s [%s] ", p.Name(),
				p.DiscordUID()))
			if member, err := s.GuildMember(reg.GuildID, p.DiscordUID()); err == nil {
				out.AddLine(sprintf("%s - Linked to %s", p.Name(),
					member.User.String()))
				cnt++
				break
			}
		}
	}

	if cnt < 1 {
		out.AddLine("No integrated users found")
		out.Status = ifaces.CommandFailure
	}

	out.Construct()
	return out, nil
}
