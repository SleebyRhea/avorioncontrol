package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func setaliasCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	reg := cmd.Registrar()

	if !HasNumArgs(a, 2, 2) {
		return "", &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`),
			cmd:     cmd}
	}

	if !regexp.MustCompile("^[a-zA-Z]{1,10}$").MatchString(a[2]) {
		return "", &ErrInvalidAlias{
			alias: a[2],
			cmd:   cmd}
	}

	if reg.IsRegistered(a[1]) == false {
		return "", &ErrInvalidCommand{
			name: a[1],
			cmd:  cmd}
	}

	if err := c.SetAliasCommand(a[1], a[2]); err != nil {
		logger.LogError(cmd, err.Error())
		return "", &ErrCommandError{
			message: "Failed to configure Alias: " + err.Error(),
			cmd:     cmd}
	}

	c.SaveConfiguration()
	return "", nil
}
