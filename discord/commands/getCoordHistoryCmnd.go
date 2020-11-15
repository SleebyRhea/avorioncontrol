package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getCoordHistoryCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {
	var (
		reg   = cmd.Registrar()
		match []string
	)

	// Require at least one set of coords
	if !HasNumArgs(a, 1, -1) {
		return "", &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`, cmd.Name()),
			cmd:     cmd}
	}

	jumps := make([]ifaces.JumpInfo, 0)
	coords := make([][2]int, 0)
	coordRe := regexp.MustCompile(`^(-?[0-9]{1,3}):(-?[0-9]{1,3})$`)

	// Validate the coords that we were given and store them as ints for easy
	//	comparison
	for _, c := range a[1:] {
		logger.LogDebug(cmd, "Operating on: "+c)
		if match = coordRe.FindStringSubmatch(c); match == nil {
			return "", &ErrInvalidArgument{
				message: sprintf("Invalid coordinate given: `%s`", c),
				cmd:     cmd}
		}

		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])

		// Make sure our coordinates are actually between -500 and 500
		if x > 500 || x < -500 {
			return "", &ErrInvalidArgument{
				message: sprintf("Coordinate **x** is out of range: `%d`", x),
				cmd:     cmd}
		}

		// Make sure our coordinates are actually between -500 and 500
		if y > 500 || y < -500 {
			return "", &ErrInvalidArgument{
				message: sprintf("Coordinate **y** is out of range: `%d`", y),
				cmd:     cmd}
		}

		coords = append(coords, [2]int{x, y})
	}

	// Migrate this to a method call on sectors or a utility function
	for _, c := range coords {
		logger.LogDebug(cmd, sprintf("Checking for jumps to sector: (%d:%d)", c[0], c[1]))
		sector := reg.server.Sector(c[0], c[1])
		if len(sector.Jumphistory) > 0 {
			orderedjumps := reverseJumps(sector.Jumphistory)
			for _, j := range orderedjumps {
				// Sectors have jumps saved as a pointer to the player jump for both easy
				// clearing, and for efficiency
				jumps = append(jumps, *j)
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
		var (
			obj ifaces.IHaveShips
		)

		fid := strconv.FormatInt(int64(j.FID), 10)

		tl := j.Time.In(loc)
		t := sprintf("%d/%02d/%02d %02d:%02d:%02d", tl.Year(), tl.Month(), tl.Day(),
			tl.Hour(), tl.Minute(), tl.Second())

		switch j.Kind {
		case "player":
			if obj = reg.server.Player(fid); obj == nil {
				logger.LogError(cmd, "Got an invalid ifaces.IHaveShips object")
				return "", &ErrCommandError{
					message: "Error, bad data type encountered. Please review the logs.",
					cmd:     cmd}
			}

		case "alliance":
			if obj = reg.server.Alliance(fid); obj == nil {
				logger.LogError(cmd, "Got an invalid ifaces.IHaveShips object")
				return "", &ErrCommandError{
					message: "Error, bad data type encountered. Please review the logs.",
					cmd:     cmd}
			}

			logger.LogError(cmd, "Got an invalid ifaces.IHaveShips object")
			return "", &ErrCommandError{
				message: "Error, bad data type encountered. Please review the logs.",
				cmd:     cmd}

		default:
			logger.LogError(cmd, "Got an invalid ifaces.IHaveShips object")
			return "", &ErrCommandError{
				message: "Error, bad data type encountered. Please review the logs.",
				cmd:     cmd}
		}

		suffix := sprintf("\n%s (%d:%d) %s/%s \"%s\"",
			t, j.X, j.Y, obj.Name(), j.Kind, j.Name)
		if len(suffix+msg) > 1900 {
			msg += "\n...(truncated due to length)"
		} else {
			msg = msg + suffix
		}
	}

	msg = msg + "\n```"
	s.ChannelMessageSend(m.ChannelID, msg)
	return "", nil
}
