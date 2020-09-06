package commands

import (
	"AvorionControl/discord/botconfig"
	"AvorionControl/gameserver"
	"AvorionControl/logger"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

func rconCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c *botconfig.Config) (string, error) {

	var (
		reg *CommandRegistrar
		cmd *CommandRegistrant
		gs  gameserver.Server
		err error
		out string
		msg string

		gscmd string
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command("rcon"); err != nil {
		return "", err
	}

	gs = reg.server

	if !HasNumArgs(a, 1, -1) {
		msg = sprintf("Invalid number of args passed to `%s`", a[0])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	gscmd = strings.Join(a[1:], " ")

	if out, err = gs.RunCommand(gscmd); err != nil {
		logger.LogError(cmd, sprintf("Failed to run \"%s\": %s", gscmd, err.Error()))
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
		return "", nil
	}

	_ = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	if strings.ReplaceAll(out, " ", "") != "" {
		out = sprintf("**Output: `%s`**\n```\n%s\n```", gscmd, out)
		if utf8.RuneCountInString(out) <= 2000 {
			_, err = s.ChannelMessageSend(m.ChannelID, out)
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "Output too large for discord")
		}
	}

	return "", nil
}

func restartServerCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c *botconfig.Config) (string, error) {
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
	c *botconfig.Config) (string, error) {
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

	if reg.server.IsUp() {
		if err = reg.server.Stop(); err != nil {
			s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
			s.ChannelMessageSend(m.ChannelID, "Encountered an error stopping Avorion")
			logger.LogError(cmd, err.Error())
			return "", nil
		}
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return "", nil
}

func startServerCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c *botconfig.Config) (string, error) {
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
				"Encountered an error starting Avorion:\n```%s\n```\n", err.Error()))
			logger.LogError(cmd, err.Error())
			return "", nil
		}
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return "", nil
}
