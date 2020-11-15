package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func loglevelCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var (
		reg    = cmd.Registrar()
		cmdobj *CommandRegistrant

		l   int
		err error
		ac  string
		ok  bool
		out string
	)

	if !HasNumArgs(a, 2, -1) {
		return "", &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`, a[0]),
			cmd:     cmd}
	}

	// Validate our loglevel
	if l, err = strconv.Atoi(a[1]); err != nil || l > 3 || l < 0 {
		out = sprintf("`%s` is not a valid loglevel. Valid levels include:\n",
			a[1]) + "```\n0 - Only output service level info\n" +
			"1 - Show warnings (default)\n" +
			"2 - Show informational output\n" +
			"3 - Debug mode\n```"
		return "", &ErrCommandError{
			message: out,
			cmd:     cmd}
	}

	logger.LogDebug(cmd, sprintf("Using loglevel %d", l))

	for _, obj := range a[2:] {
		switch obj {
		case "guild":
			out = sprintf("Default command loglevel is now: _**%d**_", l)
			reg.SetLoglevel(l)
			s.ChannelMessageSend(m.ChannelID, out)
			continue

		// REVIEW: Is user/role specific logging something I want to implement?
		case "user", "role":
			return "", &ErrCommandError{
				message: "User/Role specific logging is not yet implemented",
				cmd:     cmd}

		case "default":
			c.SetLoglevel(l)
			c.SaveConfiguration()
			s.ChannelMessageSend(m.ChannelID,
				sprintf("Default command loglevel is now: _**%d**_", l))
			return "", nil

		// Process commands/aliases here
		// TODO: Refactor this if CommandRegistrar.GetAliasedCommand and
		// BotConfig.GetAliasedCommand are refactored. See the TODO
		// on those methods for details
		default:
			out = sprintf("Command `%s` isn't registered, nor is it an alias", obj)
			logger.LogDebug(cmd, sprintf("Checking for command %s", obj))

			if cmdobj, err = reg.Command(obj); err != nil {
				logger.LogDebug(cmd,
					sprintf("Command [%s] isn't registered, checking for aliases", obj))

				if ok, ac = c.GetAliasedCommand(obj); !ok {
					s.ChannelMessageSend(m.ChannelID, out)
					return "", nil
				}

				if cmdobj, err = reg.Command(ac); err != nil {
					logger.LogError(cmd, "Found configured alias that is invalid: "+ac)
					s.ChannelMessageSend(m.ChannelID, out)
					return "", nil
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
