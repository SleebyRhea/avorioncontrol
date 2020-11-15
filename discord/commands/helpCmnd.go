package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func helpCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var (
		maincmd *CommandRegistrant
		reg     = cmd.Registrar()
		err     error
		out     string
	)

	if len(a[1:]) < 1 {
		return listCmd(s, m, a, c, cmd)
	}

	if maincmd, err = reg.Command(a[1]); err != nil {
		return "", &ErrInvalidCommand{
			name: a[1],
			cmd:  cmd}
	}

	if c.CommandDisabled(maincmd.Name()) {
		return "", &ErrCommandDisabled{cmd: cmd}
	}

	if len(a[1:]) > 1 {
		if c, cmdlets := maincmd.Subcommands(); c > 0 {
			for _, sub := range cmdlets {
				if a[2] == sub.Name() {
					maincmd = sub
					break
				}
			}
		}

		if maincmd.Name() == a[1] {
			return "", &ErrInvalidSubcommand{
				subname: a[2],
				cmd:     maincmd}
		}
	}

	if out, err = maincmd.Help(); err != nil {
		logger.LogError(maincmd, err.Error())
		return "", &ErrCommandError{
			message: "Error getting help",
			cmd:     cmd}
	}

	if _, err = s.ChannelMessageSend(m.ChannelID, out); err != nil {
		logger.LogError(cmd, err.Error())
	}

	return "", nil
}
