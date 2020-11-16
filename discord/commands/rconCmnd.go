package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func rconCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		reg = cmd.Registrar()
		srv = reg.server
		out = newCommandOutput(cmd, "RCON")

		rconout string
		rcmd    string
		err     error
	)

	if !HasNumArgs(a, 1, -1) {
		return nil, &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`, cmd.Name()),
			cmd:     cmd}
	}

	rcmd = strings.Join(a[1:], " ")
	out.Description = rcmd

	if rconout, err = srv.RunCommand(rcmd); err != nil {
		return nil, &ErrCommandError{
			message: sprintf("Failed to run `%s`. Error:\n```%s```", rcmd, err.Error()),
			cmd:     cmd}
	}

	if strings.ReplaceAll(rconout, " ", "") == "" {
		return nil, nil
	}

	for _, line := range strings.Split(rconout, "\n") {
		logger.LogDebug(cmd, "RCON: "+line)
		out.AddLine(line)
	}

	logger.LogInfo(cmd, sprintf("%s ran the rcon command: [%s]", m.Author.String(),
		strings.Join(a[1:], " ")))

	out.Quoted = true
	out.Construct()
	return out, nil
}
