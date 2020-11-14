package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func showAdminRolesSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		err   error
		out   string
		cnt   int
		guild *discordgo.Guild
	)

	if guild, err = s.Guild(m.GuildID); err != nil {
		return "", err
	}

	cnt = 0
	out = ""
	for _, r := range guild.Roles {
		if l := c.GetRoleAuth(r.ID); l > 0 {
			out += sprintf("%s - Auth level: %d\n", r.Name, l)
			cnt++
		}
	}

	if cnt > 0 {
		s.ChannelMessageSend(m.ChannelID, sprintf(
			"**_Role Authorizations:_**```\n%s```", out))
		return "", nil
	}

	s.ChannelMessageSend(m.ChannelID,
		"No roles have received authorization levels")
	return "", nil
}

func showAdminCmndsSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		err error
		out string

		cnt = 0
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	_, commands := reg.AllCommands()
	for _, name := range commands {
		l := c.GetCmndAuth(name)
		if l > 0 {
			out += sprintf("%s - Auth level requirement: %d\n", name, l)
			cnt++
		}
	}

	if cnt > 0 {
		s.ChannelMessageSend(m.ChannelID, sprintf(
			"**_Command Auth Levels:_**\n```%s```\n", out))
		return "", nil
	}

	s.ChannelMessageSend(m.ChannelID, "No commands require admin privileges")
	return "", nil
}

func addAdminRoleSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		cmd   *CommandRegistrant
		reg   *CommandRegistrar
		guild *discordgo.Guild
		level int
		role  string
		err   error
	)

	// Account for the fact that this is a subcommand by passing the HasNumArgs
	//	a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 2, -1) {
		return "", &ErrInvalidArgument{sprintf(
			`%s was passed the wrong number of arguments`, a[0])}
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	if level, err = strconv.Atoi(a[2]); err != nil {
		return "", &ErrInvalidArgument{sprintf(
			"`%s` is not a valid number", a[2])}
	}

	// Account for roles with spaces
	role = strings.Join(a[3:], " ")

	if guild, err = s.Guild(m.GuildID); err != nil {
		return "", err
	}

	// Account for either case where a role is the ID, Name, or Mention
	for _, r := range guild.Roles {
		if r.ID == role || r.Name == role || r.Mention() == role {
			c.AddRoleAuth(r.ID, level)
			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			logger.LogInfo(cmd, sprintf("%s set the authorization level for %s to %d",
				m.Author.String(), r.Name, level))
			return "", nil
		}
	}

	msg := sprintf("`%s` is not a valid role")
	return "", &ErrCommandError{msg}
}

func removeAdminRoleSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		cmd   *CommandRegistrant
		reg   *CommandRegistrar
		guild *discordgo.Guild
		level int
		err   error
	)

	// Account for the fact that this is a subcommand by passing the
	//	HasNumArgs function a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 1, -1) {
		return "", &ErrInvalidArgument{sprintf(
			`%s was passed the wrong number of arguments`, a[0])}
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	role := strings.Join(a[2:], " ")

	if guild, err = s.Guild(m.GuildID); err != nil {
		return "", err
	}

	for _, r := range guild.Roles {
		if r.ID == role || r.Name == role || r.Mention() == role {
			if err = c.RemoveRoleAuth(r.ID); err != nil {
				msg := sprintf("%s is not assigned an authorization level", role)
				s.ChannelMessageSend(m.ChannelID, msg)
				return "", nil
			}

			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			logger.LogInfo(cmd, sprintf("%s removed authorization for %s",
				m.Author.Mention(), r.Name, level))
			return "", nil
		}
	}

	msg := sprintf("`%s` is not a valid role", role)
	s.ChannelMessageSend(m.ChannelID, msg)
	return "", nil
}

func addAdminCmndSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		err error
		lvl int
	)

	// Account for the fact that this is a subcommand by passing the
	//	HasNumArgs function a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 2, 2) {
		return "", &ErrInvalidArgument{sprintf(
			"`%s:%s` was passed the wrong number of arguments (%d)", a[0], a[1], len(a[1:]))}
	}

	name := a[3]

	if lvl, err = strconv.Atoi(a[2]); err != nil {
		return "", &ErrInvalidArgument{sprintf("%s is not a valid number",
			a[3])}
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if _, err = reg.Command(name); err != nil {
		return "", &ErrCommandError{sprintf("%s is not a valid command",
			name)}
	}

	if lvl > 0 {
		c.AddCmndAuth(name, lvl)
	} else {
		c.RemoveCmndAuth(name)
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return "", nil
}

func removeAdminCmndSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		err error
	)

	// Account for the fact that this is a subcommand by passing the
	//	HasNumArgs function a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 1, 1) {
		return "", &ErrInvalidArgument{sprintf(
			`%s was passed the wrong number of arguments`, a[0])}
	}

	name := a[2]

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if _, err = reg.Command(name); err != nil {
		return "", &ErrCommandError{sprintf("%s is not a valid command",
			name)}
	}

	c.RemoveCmndAuth(name)
	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return "", nil
}
