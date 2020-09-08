package commands

import (
	"AvorionControl/ifaces"
	"AvorionControl/logger"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

func rconCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		srv ifaces.IGameServer
		reg *CommandRegistrar
		cmd *CommandRegistrant

		rcmd string
		out  string
		msg  string
		err  error
	)

	if !HasNumArgs(a, 1, -1) {
		return wrongArgsCmd(s, m, a, c)
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	srv = reg.server
	rcmd = strings.Join(a[1:], " ")

	if out, err = srv.RunCommand(rcmd); err != nil {
		logger.LogError(cmd, sprintf("Failed to run \"%s\": %s", rcmd, err.Error()))
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
		return "", nil
	}

	_ = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	if strings.ReplaceAll(out, " ", "") != "" {
		msg = sprintf("**Output: `%s`**\n```\n%s\n```", rcmd, out)
		if utf8.RuneCountInString(out) <= 2000 {
			_, err = s.ChannelMessageSend(m.ChannelID, msg)
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "Output too large for discord")
		}
	}

	return "", err
}

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
		if err = reg.server.Stop(); err != nil {
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
		if err = reg.server.Start(); err != nil {
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

func setChatChannelCmnd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator) (string, error) {
	var (
		channels []*discordgo.Channel

		reg *CommandRegistrar
		cmd *CommandRegistrant
		out string
		msg string
		err error
	)

	if !HasNumArgs(a, 1, 1) {
		return wrongArgsCmd(s, m, a, c)
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return out, err
	}

	if cmd, err = reg.Command("setchatchannel"); err != nil {
		return out, err
	}

	channels, err = s.GuildChannels(m.GuildID)

	for _, dch := range channels {
		logger.LogDebug(cmd, sprintf("Checking channel ID %s against %s", dch.ID, a[1]))
		if dch.ID == a[1] && dch.Type == discordgo.ChannelTypeGuildText {
			c.SetChatChannel(a[1])
			err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			return "", err
		}
	}

	msg = sprintf("Invalid channel ID: `%s`", a[1])
	_, err = s.ChannelMessageSend(m.ChannelID, msg)
	return out, err
}
