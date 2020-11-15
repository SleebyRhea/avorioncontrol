package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var modURLBase = `https://steamcommunity.com/sharedfiles/filedetails/?id=`

func modAddSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {

	var (
		failed  = make([]string, 0)
		reason  = make([]string, 0)
		success = make([]int64, 0)

		out = ""
	)

	if !HasNumArgs(a[1:], 1, -1) {
		return "", &ErrInvalidArgument{
			message: sprintf("`%s` was passed the wrong number of arguments", cmd.Name()),
			cmd:     cmd}
	}

	for _, arg := range a[2:] {
		if !regexp.MustCompile(`[0-9]{10}`).MatchString(arg) {
			return "", &ErrInvalidArgument{
				message: sprintf("`%s` is not a valid workshop id", arg),
				cmd:     cmd}
		}
	}

	for _, mod := range a[2:] {
		if id, err := strconv.ParseInt(mod, 10, 64); err == nil {
			if err := c.AddServerMod(id); err != nil {
				failed = append(failed, mod)
				reason = append(reason, err.Error())
			} else {
				logger.LogInfo(cmd, sprintf("%s added %d to the mod configuration",
					m.Author.String(), id))
				success = append(success, id)
			}
		} else {
			failed = append(failed, mod)
			reason = append(reason, "Not a valid workshop id")
		}
	}

	if len(success) > 0 {
		for i, id := range success {
			out += sprintf("%d - %s%d\n", i+1, modURLBase, id)
		}

		if len(out) > 1900 {
			msg := ""
			cnt := 1
			for _, line := range strings.Split(out, "\n") {
				msg += line + "\n"
				if len(msg) > 1900 {
					_, err := s.ChannelMessageSend(m.ChannelID, sprintf(
						"**Mods Added (%d):**\n```%s```", cnt, msg))
					if err != nil {
						logger.LogError(cmd, "discordgo: "+err.Error())
					}
					msg = ""
					cnt++
				}
			}

			if msg != "" {
				s.ChannelMessageSend(m.ChannelID, sprintf(
					"**Mods Added (%d):**\n```%s```", cnt, msg))
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, sprintf(
				"**Mods Added:**\n```%s```", out))
			if err != nil {
				logger.LogError(cmd, "discordgo: "+err.Error())
			}
		}
	}

	out = ""

	if len(failed) > 0 {
		for i, id := range failed {
			if len(reason) >= i {
				out += sprintf("%s: %s\n", id, reason[i])
			} else {
				out += sprintf("%s: %s\n", id, "Unspecified error")
			}
		}

		if len(out) > 1900 {
			msg := ""
			cnt := 1
			for _, line := range strings.Split(out, "\n") {
				msg += line + "\n"
				if len(msg) > 1900 {
					_, err := s.ChannelMessageSend(m.ChannelID, sprintf(
						"**Failed to Add (%d):**\n```%s```", cnt, msg))
					if err != nil {
						logger.LogError(cmd, "discordgo: "+err.Error())
					}
					msg = ""
					cnt++
				}
			}

			if msg != "" {
				s.ChannelMessageSend(m.ChannelID, sprintf(
					"**Failed to Add (%d):**\n```%s```", cnt, msg))
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, sprintf(
				"**Failed to Add:**\n```%s```", out))
			if err != nil {
				logger.LogError(cmd, "discordgo: "+err.Error())
			}
		}

		return "", &ErrCommandError{
			message: sprintf("Failed to add %d mod[s]", len(failed)),
			cmd:     cmd}
	}

	return "", nil
}

func modRemoveSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var (
		failed  = make([]string, 0)
		reason  = make([]string, 0)
		success = make([]int64, 0)

		out = ""
	)

	if !HasNumArgs(a[1:], 1, -1) {
		return "", &ErrInvalidArgument{
			message: sprintf("`%s` was passed the wrong number of arguments", cmd.Name()),
			cmd:     cmd}
	}

	for _, arg := range a[2:] {
		if !regexp.MustCompile(`[0-9]{10}`).MatchString(arg) {
			return "", &ErrInvalidArgument{
				message: sprintf("`%s` is not a valid workshop id", arg),
				cmd:     cmd}
		}
	}

	for _, mod := range a[2:] {
		if id, err := strconv.ParseInt(mod, 10, 64); err == nil {
			if err := c.RemoveServerMod(id); err != nil {
				failed = append(failed, mod)
				reason = append(reason, err.Error())
			} else {
				success = append(success, id)
				logger.LogInfo(cmd, sprintf("%s removed %d from the mod configuration",
					m.Author.String(), id))
			}
		} else {
			failed = append(failed, mod)
			reason = append(reason, "Not a valid workshop id")
		}
	}

	if len(success) > 0 {
		for i, id := range success {
			out += sprintf("%d - %s%d\n", i+1, modURLBase, id)
		}

		if len(out) > 1900 {
			msg := ""
			cnt := 1
			for _, line := range strings.Split(out, "\n") {
				msg += line + "\n"
				if len(msg) > 1900 {
					_, err := s.ChannelMessageSend(m.ChannelID, sprintf(
						"**Mods Removed (%d):**\n```%s```", cnt, msg))
					if err != nil {
						logger.LogError(cmd, "discordgo: "+err.Error())
					}
					msg = ""
					cnt++
				}
			}

			if msg != "" {
				s.ChannelMessageSend(m.ChannelID, sprintf(
					"**Mods Removed (%d):**\n```%s```", cnt, msg))
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, sprintf(
				"**Mods Removed:**\n```%s```", out))
			if err != nil {
				logger.LogError(cmd, "discordgo: "+err.Error())
			}
		}
	}

	out = ""

	if len(failed) > 0 {
		for i, id := range failed {
			if len(reason) >= i {
				out += sprintf("%s: %s\n", id, reason[i])
			} else {
				out += sprintf("%s: %s\n", id, "Unspecified error")
			}
		}

		if len(out) > 1900 {
			msg := ""
			cnt := 1
			for _, line := range strings.Split(out, "\n") {
				msg += line + "\n"
				if len(msg) > 1900 {
					_, err := s.ChannelMessageSend(m.ChannelID, sprintf(
						"**Failed to Remove (%d):**\n```%s```", cnt, msg))
					if err != nil {
						logger.LogError(cmd, "discordgo: "+err.Error())
					}
					msg = ""
					cnt++
				}
			}

			if msg != "" {
				s.ChannelMessageSend(m.ChannelID, sprintf(
					"**Failed to Remove (%d):**\n```%s```", cnt, msg))
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, sprintf(
				"**Failed to Remove:**\n```%s```", out))
			if err != nil {
				logger.LogError(cmd, "discordgo: "+err.Error())
			}
		}

		return "", &ErrCommandError{
			message: sprintf("Failed to remove %d mod[s]", len(failed)),
			cmd:     cmd}
	}

	return "", nil
}

func listModsSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var (
		mods = c.ListServerMods()
		list = ""
		msg  = ""
	)

	for i, id := range mods {
		list += sprintf("%d: %s%d\n", i+1, modURLBase, id)
	}

	if list == "" {
		s.ChannelMessageSend(m.ChannelID, "No mods currently configured")
		return "", nil
	}

	if len(list) > 1900 {
		cnt := 1
		for _, line := range strings.Split(list, "\n") {
			msg += line + "\n"
			if len(msg) > 1900 {
				s.ChannelMessageSend(m.ChannelID, sprintf(
					"**Mods Installed (%d):**\n```%s```", cnt, msg))
				msg = ""
				cnt++
			}
		}

		if msg != "" {
			s.ChannelMessageSend(m.ChannelID, sprintf(
				"**Mods Installed (%d):**\n```%s```", cnt, msg))
		}

		return "", nil
	}

	s.ChannelMessageSend(m.ChannelID, sprintf(
		"**Mods Installed:**\n```%s```", list))
	return "", nil
}
