package avorion

import (
	"AvorionControl/gameserver"
	"AvorionControl/logger"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var resourceMap map[string]int

// Player is a player that has connected to the server at some point, and has
// data present in the game db
type Player struct {
	gameserver.IPlayer

	// avorion.Player
	index     string
	steam64   int64
	coords    [2]int
	oldcoords [][2]int
	resources map[string]int64
	discordid string

	// gameserver.Player
	ip       net.IP
	name     string
	online   bool
	server   *Server
	loglevel int
}

// Update updates our tracked data for the player
func (p *Player) Update() error {
	logger.LogDebug(p.server, "Attempting to update player: "+p.index)
	cmd := "playerinfo -i -s -t -c -p " + fmt.Sprint(p.index)

	if p.online {
		cmd += " -a"
	}

	out, err := p.server.RunCommand(cmd)

	if err != nil {
		logger.LogError(p.server, fmt.Sprintf(
			"Failed to acquire player data for: %s (%s)", p.index, p.name))
		return err
	}

	out = strings.TrimSuffix(out, "\n")
	logger.LogDebug(p.server, fmt.Sprintf("Processing: (%s)", out))
	for _, o := range strings.Split(out, "\n") {
		if m := rePlayerDataFull.FindStringSubmatch(o); m != nil {
			logger.LogDebug(p.server, "Found online player string")
			p.online = true

			var (
				x   int
				y   int
				err error
			)

			// Perform data sanity checks
			id, _ := strconv.ParseInt(m[1], 10, 64)

			if p.steam64 == 0 {
				p.steam64 = id
			}

			if id != p.steam64 {
				return fmt.Errorf("Failed to update playerdata, Steam64 ID mismatch (%s != %d)",
					m[1], p.steam64)
			}

			if m[2] != p.index {
				return fmt.Errorf("Failed to update playerdata, index mismatch (%s != %s)",
					m[2], p.index)
			}

			coords := strings.Split(m[3], ":")

			if x, err = strconv.Atoi(coords[0]); err != nil {
				return fmt.Errorf("Failed to update playerdata, x couldn't be converted: (%s)",
					m[3])
			}

			if y, err = strconv.Atoi(coords[1]); err != nil {
				return fmt.Errorf("Failed to update playerdata, y couldn't be converted: (%s)",
					m[3])
			}

			p.UpdateCoords(x, y)
			p.name = m[4]
			if p.ip = net.ParseIP(m[6]); p.ip == nil {
				return fmt.Errorf("Failed to parse player IP address: (%s)", m[4])
			}

			logger.LogInfo(p.server, fmt.Sprintf(
				"Updated player information for %s (%s:%s)", p.index, p.name, p.ip.String()))

			p.discordid, _ = p.server.RunCommand("getlinkeddiscord " + m[2])
			continue
		}

		if m := rePlayerDataOffline.FindStringSubmatch(o); m != nil {
			logger.LogDebug(p.server, "Found offline player string")
			p.online = false

			// Perform data sanity checks
			id, _ := strconv.ParseInt(m[1], 10, 64)

			if p.steam64 == 0 {
				p.steam64 = id
			}

			if id != p.steam64 {
				return fmt.Errorf("Failed to update playerdata, Steam64 ID mismatch (%s != %d)",
					m[1], p.steam64)
			}

			if m[2] != p.index {
				return fmt.Errorf("Failed to update playerdata, index mismatch (%s != %s)",
					m[2], p.index)
			}

			p.name = m[3]
			p.discordid, _ = p.server.RunCommand("getlinkeddiscord " + m[2])
			logger.LogInfo(p.server, fmt.Sprintf(
				"Updated player information for %s (%s)", p.index, p.name))
			continue
		}

		if m := rePlayerAlliance.FindStringSubmatch(o); m != nil {
			// TODO: Add alliance tracking. avorion.Allliance needs to be implented
			logger.LogDebug(p.server, "Found alliance string (unimplemented)")
			continue
		}

		return fmt.Errorf("Unable to parse the following line: (%s)", o)
	}

	return nil
}

/******************/
/* avorion.Player */
/******************/

// Index returns the players in-game index
func (p *Player) Index() string {
	return p.index
}

// Steam64 returns the players steam64 ID
func (p *Player) Steam64() int64 {
	return p.steam64
}

// UpdateCoords saves the previous coordinates that the player was in, and sets
// their current position. Saves up to 100 previous coordinate locations
func (p *Player) UpdateCoords(x, y int) {
	p.oldcoords = append(p.oldcoords, p.coords)
	p.oldcoords = p.oldcoords[1:]
	p.coords = [2]int{x, y}
}

/*********************/
/* gameserver.Player */
/*********************/

// IP returns the IP address that the player used to connect this session
func (p *Player) IP() net.IP {
	return p.ip
}

// SetIP sets or updates a players IP address
func (p *Player) SetIP(ips string) {
	p.ip = net.ParseIP(ips)
}

// Name returns the name of the player
func (p *Player) Name() string {
	return p.name
}

// Kick kicks the player
func (p *Player) Kick(r string) {
	p.server.RunCommand("kick" + p.Name())
}

// Ban bans the player
func (p *Player) Ban(r string) {
	p.server.RunCommand("ban " + p.Name())
}

// Online returns the current online status of the player
func (p *Player) Online() bool {
	return p.online
}

// SetOnline updates the player status to the boolean passed
func (p *Player) SetOnline(o bool) {
	p.online = o
}

// SetDiscordUID sets a players Discord ID
func (p *Player) SetDiscordUID(uid string) {
	p.discordid = uid
}

// DiscordUID returns a players Discord ID
func (p *Player) DiscordUID() string {
	return p.discordid
}

// Message messages a player
func (p *Player) Message(string) {
}

/************************/
/* IFace logger.ILogger */
/************************/

// UUID returns the UUID of a player
func (p *Player) UUID() string {
	return fmt.Sprintf("%s:%s", p.index, p.name)
}

// Loglevel returns the loglevel of a player
func (p *Player) Loglevel() int {
	return p.loglevel
}

// SetLoglevel sets the loglevel of a player
func (p *Player) SetLoglevel(l int) {
	p.loglevel = l
}
