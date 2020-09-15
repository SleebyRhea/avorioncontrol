package commands

import (
	"AvorionControl/ifaces"
	"AvorionControl/logger"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getJumpsCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		obj ifaces.IHaveShips
		reg *CommandRegistrar
		cmd *CommandRegistrant
		err error
		cnt int
	)

	// Make sure we have the args we need
	if !HasNumArgs(a, 3, -1) {
		return wrongArgsCmd(s, m, a, c)
	}

	kind := a[1]
	ref := strings.Join(a[3:], " ")

	// FIXME: Prevent overflows here
	if cnt, err = strconv.Atoi(a[2]); err != nil {
		msg := sprintf("`%s` is not a valid integer.", a[2])
		s.ChannelMessageSend(m.ChannelID, msg)
		return "", err
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	if kind == "player" {
		if p := reg.server.PlayerFromName(ref); p != nil {
			obj = p
		} else if p := reg.server.Player(ref); p != nil {
			obj = p
		} else if p := reg.server.PlayerFromDiscord(ref); p != nil {
			obj = p
		}
		// } else if kind == "alliance" {
		// 	if a := reg.server.AllianceFromName(ref); a != nil {
		// 		obj = a
		// 	} else if a := reg.server.Player(ref); a != nil {
		// 		obj = a
		// 	}
	} else {
		msg := sprintf("`%s` is not a valid option (player or alliance expected)", kind)
		s.ChannelMessageSend(m.ChannelID, msg)
	}

	if obj == nil {
		msg := sprintf("`%s` is not a valid player or alliance reference", ref)
		s.ChannelMessageSend(m.ChannelID, msg)
		return "", err
	}

	loc, err := time.LoadLocation(c.TimeZone())
	if err != nil {
		msg := sprintf("Failed to load timezone: `%s`", c.TimeZone())
		s.ChannelMessageSend(m.ChannelID, msg)
		logger.LogError(cmd, err.Error())
		return "", nil
	}

	if jumps := obj.GetLastJumps(cnt); len(jumps) > 0 {
		msg := sprintf("**Jumps for %s**:```", obj.Name())
		for _, j := range jumps {
			t := j.Time.In(loc)
			msg = sprintf("%s\n%d/%02d/%02d %02d:%02d:%02d - (%d:%d) %s", msg,
				t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second(), j.X, j.Y, j.Name)
		}
		msg = msg + "```"
		s.ChannelMessageSend(m.ChannelID, msg)
	} else {
		msg := sprintf("Player **%s** has no recorded jump history", ref)
		s.ChannelMessageSend(m.ChannelID, msg)
	}

	return "", err
}
