package commands

import (
	"avorioncontrol/ifaces"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func statusCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		out = newCommandOutput(cmd, "Server Status")
		reg = cmd.Registrar()
		srv = reg.server

		ret string
		err error
	)

		out.Monospace = true
	rcmd := "status"

	if ret, err = srv.RunCommand(rcmd); err != nil {
		return nil, &ErrCommandError{
			message: sprintf("Failed to run `%s`. Error:\n```%s```", rcmd, err.Error()),
			cmd:     cmd}
	}

	if strings.ReplaceAll(ret, " ", "") != "" {
		for _, line := range strings.Split(ret, "\n")[0:10] {
			out.AddLine(line)
		}
	} else {
		return nil, &ErrCommandError{
			message: "Invalid output recieved from status",
			cmd:     cmd}
	}

	out.Construct()
	return out, nil
}
