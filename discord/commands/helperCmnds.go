package commands

import (
	"avorioncontrol/ifaces"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var sprintf = fmt.Sprintf

// Command to be used when the command being created is intended to be used with
// subcommands
func proxySubCmnd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var out string

	if !HasNumArgs(a, 1, -1) {
		return "", &ErrInvalidArgument{
			message: sprintf("`%s` was not passed a subcommand to run", a[0]),
			cmd:     cmd}
	}

	_, cmdlets := cmd.Subcommands()
	for _, cmdlet := range cmdlets {
		if a[1] == cmdlet.Name() {
			return cmdlet.exec(s, m, a, c, cmdlet)
		}
	}

	return out, &ErrInvalidSubcommand{
		subname: a[1],
		cmd:     cmd}
}
