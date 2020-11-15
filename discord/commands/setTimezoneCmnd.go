package commands

import (
	"avorioncontrol/ifaces"
	"time"

	"github.com/bwmarrin/discordgo"
)

func setTimezoneCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	if !HasNumArgs(a, 1, 1) {
		return "", &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`, a[0]),
			cmd:     cmd}
	}

	if _, err := time.LoadLocation(a[1]); err != nil {
		return "", &ErrCommandError{
			message: sprintf("Incorrect timezone: `%s` (%s)", a[1], err.Error()),
			cmd:     cmd}
	}

	c.SetTimeZone(a[1])
	c.SaveConfiguration()
	return sprintf("User %s changed the timezone to %s", m.Author.String(), a[1]),
		nil
}
