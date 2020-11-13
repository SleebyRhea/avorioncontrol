package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"errors"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func loglevelCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (out string, err error) {
	var (
		reg    *CommandRegistrar
		cmd    *CommandRegistrant
		cmdobj *CommandRegistrant

		l  int
		ac string
		ok bool
	)

	if !HasNumArgs(a, 2, -1) {
		return "", &ErrInvalidArgument{sprintf(
			`%s was passed the wrong number of arguments`, a[0])}
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command("loglevel"); err != nil {
		return "", err
	}

	// Validate our loglevel
	if l, err = strconv.Atoi(a[1]); err != nil || l > 3 || l < 0 {
		out = sprintf("`%s` is not a valid loglevel. Valid levels include:\n",
			a[1])
		out = out + "```\n0 - Only output service level info\n" +
			"1 - Show warnings (default)\n" +
			"2 - Show informational output\n" +
			"3 - Debug mode\n```"
		return "", &ErrCommandError{out}
	}

	logger.LogDebug(cmd, sprintf("Using loglevel %d", l))

	for _, obj := range a[2:] {
		switch obj {
		case "registrar", "guild":
			out = sprintf("Default loglevel is now: _**%d**_", l)
			reg.SetLoglevel(l)
			s.ChannelMessageSend(m.ChannelID, out)
			continue

		// REVIEW: Is user/role specific logging something I want to implement?
		case "user", "role":
			out = "User/Role specific logging is not yet implemented"
			s.ChannelMessageSend(m.ChannelID, out)
			return

		// TODO: Implement core/config level changing
		case "core":
			out = "Changing core logging via command is not yet implemented"
			s.ChannelMessageSend(m.ChannelID, out)
			return

		// Process commands/aliases here
		// TODO: Refactor this if CommandRegistrar.GetAliasedCommand and
		// BotConfig.GetAliasedCommand are refactored. See the TODO
		// on those methods for details
		default:
			out = sprintf("Command `%s` isn't registered, nor is it an alias", obj)
			logger.LogDebug(cmd, sprintf("Checking for command %s", obj))

			if cmdobj, err = reg.Command(obj); err != nil {
				logger.LogDebug(cmd,
					sprintf("Command %s isn't registered, checking for aliases",
						obj))

				if ok, ac = c.GetAliasedCommand(obj); ok == false {
					s.ChannelMessageSend(m.ChannelID, out)
					return "", nil
				}

				if cmdobj, err = reg.Command(ac); err != nil {
					s.ChannelMessageSend(m.ChannelID, out)
					errOut := sprintf("Configured command alias is invalid: `%s`",
						ac)
					return "", errors.New(errOut)
				}
			}

			// Set that commands loglevel
			out = sprintf("Level for `%s` is now: _**%d**_", cmdobj.Name(), l)
			s.ChannelMessageSend(m.ChannelID, out)
			cmdobj.SetLoglevel(l)
			continue
		}
	}

	return "", nil
}
