package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func pingCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	logger.LogInfo(cmd, "Pong request received")
	out := newCommandOutput(cmd, "Ping")
		out.AddLine("Pong!")
	out.Construct()
	return out, nil
}
