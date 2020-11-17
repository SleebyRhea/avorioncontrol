package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func showAdminRolesSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		out = newCommandOutput(cmd, "Admin Role List")
		cnt = 0
	)

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		return nil, &ErrCommandError{
			message: "Could not find guild: " + err.Error(),
			cmd:     cmd}
	}

	cnt = 0
	for _, r := range guild.Roles {
		if l := c.GetRoleAuth(r.ID); l > 0 {
			out.AddLine(sprintf("_**%s's**_ authorization Level: **%d**", r.Name, l))
			cnt++
		}
	}

	if cnt == 0 {
		out.AddLine("No roles have received authorization levels")
	}

	out.Quoted = true
	out.Construct()
	return out, nil
}

func showAdminCmndsSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		reg = cmd.Registrar()
		out = newCommandOutput(cmd, "Admin Command List")
		cnt = 0
	)

	_, commands := reg.AllCommands()
	for _, name := range commands {
		l := c.GetCmndAuth(name)
		if l > 0 {
			out.AddLine(sprintf("_**%s**_ authorization level requirement: **%d**", name, l))
			cnt++
		}
	}

	if cnt == 0 {
		out.AddLine("No commands require admin privileges")
	}

	out.Quoted = true
	out.Construct()
	return out, nil
}

func addAdminRoleSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		guild *discordgo.Guild
		level int
		err   error
		role  string
	)

	// Account for the fact that this is a subcommand by passing the HasNumArgs
	//	a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 2, -1) {
		return nil, &ErrInvalidArgument{
			message: sprintf("`%s` was passed the wrong number of arguments", cmd.Name()),
			cmd:     cmd}
	}

	if level, err = strconv.Atoi(a[2]); err != nil {
		return nil, &ErrInvalidArgument{
			message: sprintf("`%s` is not a valid number", a[2]),
			cmd:     cmd}
	}

	// Account for roles with spaces
	role = strings.Join(a[3:], " ")

	if guild, err = s.Guild(m.GuildID); err != nil {
		return nil, &ErrCommandError{
			message: "Could not find guild: " + err.Error(),
			cmd:     cmd}
	}

	// Account for either case where a role is the ID, Name, or Mention
	for _, r := range guild.Roles {
		if r.ID == role || r.Name == role || r.Mention() == role {
			c.AddRoleAuth(r.ID, level)
			logger.LogInfo(cmd, sprintf("%s set the authorization level for %s to %d",
				m.Author.String(), r.Name, level))
			return nil, nil
		}
	}

	return nil, &ErrCommandError{
		message: sprintf("`%s` is not a valid role", role),
		cmd:     cmd}
}

func removeAdminRoleSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		guild *discordgo.Guild
		level int
		err   error
	)

	// Account for the fact that this is a subcommand by passing the
	//	HasNumArgs function a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 1, -1) {
		return nil, &ErrInvalidArgument{
			message: sprintf("`%s` was passed the wrong number of arguments", cmd.Name()),
			cmd:     cmd}
	}

	role := strings.Join(a[2:], " ")

	if guild, err = s.Guild(m.GuildID); err != nil {
		return nil, &ErrCommandError{
			message: "Could not find guild: " + err.Error(),
			cmd:     cmd}
	}

	for _, r := range guild.Roles {
		if r.ID == role || r.Name == role || r.Mention() == role {
			if err = c.RemoveRoleAuth(r.ID); err != nil {
				return nil, &ErrCommandError{
					message: sprintf("%s is not assigned an authorization level", role),
					cmd:     cmd}
			}

			logger.LogInfo(cmd, sprintf("%s removed authorization for %s",
				m.Author.Mention(), r.Name, level))
			return nil, nil
		}
	}

	return nil, &ErrInvalidArgument{
		message: sprintf("`%s` is not a valid role", role),
		cmd:     cmd}
}

func addAdminCmndSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		err error
		lvl int
	)

	// Account for the fact that this is a subcommand by passing the
	//	HasNumArgs function a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 2, 2) {
		return nil, &ErrInvalidArgument{
			message: sprintf("`%s` was passed the wrong number of arguments", cmd.Name()),
			cmd:     cmd}
	}

	name := a[3]

	if lvl, err = strconv.Atoi(a[2]); err != nil {
		return nil, &ErrInvalidArgument{
			message: sprintf("%s is not a valid number", a[2]),
			cmd:     cmd}
	}

	if lvl > 0 {
		c.AddCmndAuth(name, lvl)
	} else {
		return nil, &ErrInvalidArgument{
			message: "Please specify a value higher than 0",
			cmd:     cmd}
	}

	logger.LogInfo(cmd, sprintf("%s set the authorization level for [%s] to %d",
		m.Author.String(), name, lvl))
	return nil, nil
}

func removeAdminCmndSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		reg = cmd.Registrar()
		err error
	)

	// Account for the fact that this is a subcommand by passing the
	//	HasNumArgs function a slice of the args removing the first argument
	if !HasNumArgs(a[1:], 1, 1) {
		return nil, &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`, cmd.Name()),
			cmd:     cmd}
	}

	name := a[2]

	if _, err = reg.Command(name); err != nil {
		return nil, &ErrInvalidArgument{
			message: sprintf("%s is not a valid command", name),
			cmd:     cmd}
	}

	c.RemoveCmndAuth(name)
	logger.LogInfo(cmd, sprintf("%s removed the authorization requirements for %s",
		m.Author.String(), name))
	return nil, nil
}

func showSelfCmndSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		reg     = cmd.Registrar()
		out     = newCommandOutput(cmd, "Display the users command authorization level")
		authlvl = 0
	)

	// Get the users authorization level
	member, _ := s.GuildMember(reg.GuildID, m.Author.ID)
	for _, r := range member.Roles {
		if l := c.GetRoleAuth(r); l > authlvl {
			authlvl = l
		}
	}

	out.Header = "Result"
	out.Quoted = true
	out.AddLine(sprintf("Your authorization level is: **%d**", authlvl))
	out.Construct()
	return out, nil
}
