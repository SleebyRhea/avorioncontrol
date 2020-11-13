package commands

import (
	"avorioncontrol/ifaces"
	"time"

	"github.com/bwmarrin/discordgo"
)

func setTimezoneCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	if !HasNumArgs(a, 1, 1) {
		return "", &ErrInvalidArgument{sprintf(
			`%s was passed the wrong number of arguments`, a[0])}
	}

	if _, err := time.LoadLocation(a[1]); err != nil {
		return "", &ErrCommandError{sprintf(
			"Incorrect timezone: `%s` (%s)", a[1], err.Error())}
	}

	c.SetTimeZone(a[1])
	c.SaveConfiguration()
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return "", err
}
