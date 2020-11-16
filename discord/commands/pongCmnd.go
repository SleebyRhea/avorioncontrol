package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func pongCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	logger.LogInfo(cmd, "Ping request recieved")
	out := newCommandOutput(cmd, "Pong")
		out.AddLine("Ping~")
	out.Construct()
	return out, nil
}
