package commands

import (
	"avorioncontrol/ifaces"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getJumpsCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		reg = cmd.Registrar()
		out = newCommandOutput(cmd, "Get Jump History")
		obj ifaces.IHaveShips
		err error
		cnt int
	)

	// Make sure we have the args we need
	if !HasNumArgs(a, 2, -1) {
		return nil, &ErrInvalidArgument{
			message: sprintf("`%s` was passed the wrong number of arguments", a[0]),
			cmd:     cmd}
	}

	ref := strings.Join(a[2:], " ")
	out.Quoted = true

	// FIXME: Prevent overflows here
	if cnt, err = strconv.Atoi(a[1]); err != nil {
		return nil, &ErrInvalidArgument{
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
		return nil, &ErrInvalidArgument{
			message: sprintf("`%s` is not a valid player or alliance reference", ref),
			cmd:     cmd}
	}

	loc, err := time.LoadLocation(c.TimeZone())
	if err != nil {
		return nil, &ErrInvalidTimezone{
			tz:  c.TimeZone(),
			cmd: cmd}
	}

	out.Header = "Jumps for " + obj.Name()
	if cnt > 250 {
		cnt = 250
	}

	if jumps := obj.GetLastJumps(cnt); len(jumps) > 0 {
		for _, j := range jumps {
			t := j.Time.In(loc)
			out.AddLine(sprintf("%d/%02d/%02d %02d:%02d:%02d  (%d:%d) %s",
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second(), j.X, j.Y, j.Name))
		}
	} else {
		out.AddLine(sprintf("Player **%s** has no recorded jump history", ref))
	}

	out.Construct()
	return out, nil
}
