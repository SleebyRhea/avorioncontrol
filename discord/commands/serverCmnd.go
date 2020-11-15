package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func restartServerCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {

	reg := cmd.Registrar()
	if err := reg.server.Restart(); err != nil {
		logger.LogError(cmd, "Avorion: "+err.Error())
		return "", &ErrCommandError{
			message: "Error restarting Avorion: " + err.Error(),
			cmd:     cmd}
	}

	return "", nil
}

func stopServerCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {

	reg := cmd.Registrar()
	if reg.server.IsUp() {
		if err := reg.server.Stop(true); err != nil {
			logger.LogError(cmd, "Avorion: "+err.Error())
			return "", &ErrCommandError{
				message: "Error stopping Avorion: " + err.Error(),
				cmd:     cmd}
		}
	}

	return "", nil
}

func startServerCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {

	reg := cmd.Registrar()
	if !reg.server.IsUp() {
		if err := reg.server.Start(true); err != nil {
			s.ChannelMessageSend(m.ChannelID, sprintf(
				"Encountered an error starting the server:\n```%s\n```\n", err.Error()))
			logger.LogError(cmd, "Avorion: "+err.Error())
			return "", &ErrCommandError{
				message: "Error starting Avorion: " + err.Error(),
				cmd:     cmd}
		}
	}

	return "", nil
}
