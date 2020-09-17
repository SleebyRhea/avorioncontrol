package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func restartServerCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		cmd *CommandRegistrant
		err error
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	if err = reg.server.Restart(); err != nil {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
		s.ChannelMessageSend(m.ChannelID, "Encountered an error restarting Avorion")
		logger.LogError(cmd, err.Error())
		return "", nil
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return "", nil
}

func stopServerCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		cmd *CommandRegistrant
		err error
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	if reg.server.IsUp() {
		if err = reg.server.Stop(true); err != nil {
			s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
			s.ChannelMessageSend(m.ChannelID, "Encountered an error stopping the server")
			logger.LogError(cmd, err.Error())
			return "", nil
		}
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return "", nil
}

func startServerCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		cmd *CommandRegistrant
		err error
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command("server"); err != nil {
		return "", err
	}

	if !reg.server.IsUp() {
		if err = reg.server.Start(true); err != nil {
			s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
			s.ChannelMessageSend(m.ChannelID, sprintf(
				"Encountered an error starting the server:\n```%s\n```\n", err.Error()))
			logger.LogError(cmd, err.Error())
			return "", nil
		}
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return "", nil
}
