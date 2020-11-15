package commands

import (
	"avorioncontrol/ifaces"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getJumpsCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var (
		reg = cmd.Registrar()
		obj ifaces.IHaveShips
		err error
		cnt int
	)

	// Make sure we have the args we need
	if !HasNumArgs(a, 2, -1) {
		return "", &ErrInvalidArgument{
			message: sprintf("`%s` was passed the wrong number of arguments", a[0]),
			cmd:     cmd}
	}

	ref := strings.Join(a[2:], " ")

	// FIXME: Prevent overflows here
	if cnt, err = strconv.Atoi(a[1]); err != nil {
		return "", &ErrInvalidArgument{
			message: sprintf("`%s` is not a valid number.", a[1]),
			cmd:     cmd}
	}

	if p := reg.server.PlayerFromName(ref); p != nil {
		obj = p
	} else if p := reg.server.PlayerFromDiscord(ref); p != nil {
		obj = p
	} else if p := reg.server.Player(ref); p != nil {
		obj = p
	} else if a := reg.server.Alliance(ref); a != nil {
		obj = a
	} else if a := reg.server.AllianceFromName(ref); a != nil {
		obj = a
	}

	if obj == nil {
		return "", &ErrInvalidArgument{
			message: sprintf("`%s` is not a valid player or alliance reference", ref),
			cmd:     cmd}
	}

	loc, err := time.LoadLocation(c.TimeZone())
	if err != nil {
		return "", &ErrInvalidTimezone{
			tz:  c.TimeZone(),
			cmd: cmd}
	}

	if cnt > 25 {
		cnt = 25
	}

	if jumps := obj.GetLastJumps(cnt); len(jumps) > 0 {
		msg := sprintf("**Jumps for %s**:```", obj.Name())
		for _, j := range jumps {
			t := j.Time.In(loc)
			msg = sprintf("%s\n%d/%02d/%02d %02d:%02d:%02d - (%d:%d) %s", msg,
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second(), j.X, j.Y, j.Name)
			if len(msg) > 1900 {
				msg += "\n...(truncated due to length)"
				break
			}
		}

		msg += "```"
		s.ChannelMessageSend(m.ChannelID, msg)
	} else {
		msg := sprintf("Player **%s** has no recorded jump history", ref)
		s.ChannelMessageSend(m.ChannelID, msg)
	}

	return "", nil
}
