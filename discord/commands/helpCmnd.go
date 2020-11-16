package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func helpCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		maincmd *CommandRegistrant
		out     = newCommandOutput(cmd, "Command Help")
		reg     = cmd.Registrar()
		err     error
		help    string
	)

	if len(a[1:]) < 1 && a[0] == "help" {
		return listCmd(s, m, a, c, cmd)
	}

	if a[0] == "help" {
		if maincmd, err = reg.Command(a[1]); err != nil {
			return nil, &ErrInvalidCommand{
				name: a[1],
				cmd:  cmd}
		}
	} else {
		a = a[0:1]
		maincmd = cmd
	}

	if c.CommandDisabled(maincmd.Name()) {
		return nil, &ErrCommandDisabled{cmd: cmd}
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
			return nil, &ErrInvalidSubcommand{
				subname: a[2],
				cmd:     maincmd}
		}
	}

	if help, err = maincmd.Help(); err != nil {
		logger.LogError(maincmd, err.Error())
		return nil, &ErrCommandError{
			message: "Error getting help",
			cmd:     cmd}
	}

	for _, line := range strings.Split(help, "\n") {
		out.AddLine(line)
	}

	out.Description = sprintf("_%s_", maincmd.description)
	out.Header = "Usage"
	out.Construct()
	return out, nil
}
