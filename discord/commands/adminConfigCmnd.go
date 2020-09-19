package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func addAdminSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		cmd   *CommandRegistrant
		reg   *CommandRegistrar
		guild *discordgo.Guild
		level int
		err   error
	)

	// Account for the fact that this is a subcommand by passing the HasNumArgs
	//	a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 2, 1) {
		return wrongArgsCmd(s, m, a, c)
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	role := a[2]

	if level, err = strconv.Atoi(a[3]); err != nil {
		return "", err
	}

	if guild, err = s.Guild(m.GuildID); err != nil {
		return "", err
	}

	for _, r := range guild.Roles {
		if r.ID == role || r.Name == role {
			c.AddRoleAuth(r.ID, level)
			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			logger.LogInfo(cmd, sprintf("%s set the authorization level for %s to %d",
				m.Author.Mention, r.Name, level))
			return "", nil
		}
	}

	msg := sprintf("`%s` is not a valid role")
	s.ChannelMessageSend(m.ChannelID, msg)
	return "", nil
}

func removeAdminSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		cmd   *CommandRegistrant
		reg   *CommandRegistrar
		guild *discordgo.Guild
		level int
		err   error
	)

	// Account for the fact that this is a subcommand by passing the HasNumArgs
	//	a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 2, 1) {
		return wrongArgsCmd(s, m, a, c)
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	role := a[2]

	if level, err = strconv.Atoi(a[3]); err != nil {
		return "", err
	}

	if guild, err = s.Guild(m.GuildID); err != nil {
		return "", err
	}

	for _, r := range guild.Roles {
		if r.ID == role || r.Name == role {
			if err = c.RemoveRoleAuth(r.ID); err != nil {
				msg := sprintf("%s is not assigned an authorization level", role)
				s.ChannelMessageSend(m.ChannelID, msg)
				return "", nil
			}

			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			logger.LogInfo(cmd, sprintf("%s removed authorization for %s",
				m.Author.Mention, r.Name, level))
			return "", nil
		}
	}

	msg := sprintf("`%s` is not a valid role")
	s.ChannelMessageSend(m.ChannelID, msg)
	return "", nil
}

func showAdminSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	return "", nil
}

func setNeedAdminSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	return "", nil
}
