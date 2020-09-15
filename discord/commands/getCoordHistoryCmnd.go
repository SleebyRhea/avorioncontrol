package commands

import (
	"AvorionControl/ifaces"
	"AvorionControl/logger"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type jumpInfo struct {
	Jump ifaces.ShipCoordData
	Name string
	Kind string
}

func getCoordHistoryCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		cmd *CommandRegistrant

		match []string
		err   error
	)

	// Require at least one set of coords
	if !HasNumArgs(a, 1, -1) {
		return wrongArgsCmd(s, m, a, c)
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	jumps := make([]jumpInfo, 0)
	coords := make([][2]int, 0)
	coordRe := regexp.MustCompile(`^(-?[0-9]{1,3}):(-?[0-9]{1,3})$`)

	// Validate the coords that we were given and store them as ints for easy
	//	comparison
	for _, c := range a[1:] {
		logger.LogDebug(cmd, "Operating on: "+c)
		if match = coordRe.FindStringSubmatch(c); match == nil {
			s.ChannelMessageSend(m.ChannelID, sprintf("Invalid coordinate given: `%s`",
				c))
			return "", nil
		}

		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		coords = append(coords, [2]int{x, y})
	}

	// TODO: Track individual sector history to make this kind of thing unneeded
	for _, p := range reg.server.Players() {
		for _, j := range p.GetLastJumps(-1) {
			logger.LogDebug(cmd, "Processing player: "+p.Name())
			for _, c := range coords {
				logger.LogDebug(cmd, sprintf("Checking jump to coord: (%d:%d)",
					j.X, j.Y))
				if j.X == c[0] && j.Y == c[1] {
					logger.LogDebug(cmd, "Found match")
					jumps = append(jumps, jumpInfo{
						Name: p.Name(),
						Jump: j,
						Kind: "player"})
				}
			}
		}
	}

	// TODO: Track individual sector history to make this kind of thing unneeded
	for _, a := range reg.server.Alliances() {
		logger.LogDebug(cmd, "Processing alliance: "+a.Name())
		for _, j := range a.GetLastJumps(-1) {
			logger.LogDebug(cmd, sprintf("Checking jump to coord: (%d:%d)",
				j.X, j.Y))
			for _, c := range coords {
				if j.X == c[0] && j.Y == c[1] {
					logger.LogDebug(cmd, "Found match")
					jumps = append(jumps, jumpInfo{
						Name: a.Name(),
						Jump: j,
						Kind: "alliance"})
				}
			}
		}
	}

	if len(jumps) == 0 {
		msg := "No results found"
		s.ChannelMessageSend(m.ChannelID, msg)
		return "", nil
	}

	loc, _ := time.LoadLocation(c.TimeZone())

	// TODO: This really should use an embed for paginated output
	msg := "**Jump history search results:**\n```"
	for _, j := range jumps {
		tl := j.Jump.Time.In(loc)
		t := sprintf("%d/%02d/%02d %02d:%02d:%02d", tl.Year(), tl.Month(), tl.Day(),
			tl.Hour(), tl.Minute(), tl.Second())
		msg = sprintf("%s\n%s (%d:%d) %s/%s \"%s\"",
			msg, t, j.Jump.X, j.Jump.Y, j.Name, j.Kind, j.Jump.Name)
	}

	msg = msg + "\n```"
	s.ChannelMessageSend(m.ChannelID, msg)
	return "", err
}
