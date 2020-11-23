package avorion

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"net"
	"regexp"
	"strconv"

	"time"
)

var resourceMap map[string]int

const (
	steamUIDRegex   = `^\s*([0-9]+) .*$`
	steamUIDCommand = `playerinfo %s -s -o`
)

var steamUIDRegexp = regexp.MustCompile(steamUIDRegex)

// Player is a player that has connected to the server at some point, and has
// data present in the game db
type Player struct {
	ifaces.IPlayer

	// avorion.Player
	index     string
	steam64   int64
	discordid string

	// ifaces.Player
	ip       net.IP
	name     string
	online   bool
	server   *Server
	loglevel int

	// playerdata
	resources   map[string]int64
	jumphistory []ifaces.ShipCoordData
}

// Update gathers playerdata from the server and updates our cache
func (p *Player) Update() error {
	return nil
}

// UpdateFromData updates the players information using the data from
//	a successful rePlayerData match
func (p *Player) UpdateFromData(d [15]string) error {
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

// AddJump registers a jump that a player took into a system
func (p *Player) AddJump(sc ifaces.ShipCoordData) {
	sc.Time = time.Now()
	p.jumphistory = append(p.jumphistory, sc)
	if len(p.jumphistory) > 1000 {
		p.jumphistory = p.jumphistory[1:]
	}

	sector := p.server.Sector(sc.X, sc.Y)
	fid64, _ := strconv.ParseInt(p.index, 10, 32)
	fid := int(fid64)

	// Add a pointer to the players jump to the sector history for our
	//	usage later on
	jump := &ifaces.JumpInfo{
		Name: sc.Name,
		Kind: "player",
		Time: sc.Time,
		FID:  fid,
		X:    sc.X,
		Y:    sc.Y}

	sector.Jumphistory = append(sector.Jumphistory, jump)

	id, _ := strconv.Atoi(p.Index())
	p.server.tracking.AddJump(sector.Index, int64(id), 0, *jump)
	logger.LogDebug(p, "Updated jumphistory")
}

/************************/
/* IFace ifaces.IPlayer */
/************************/

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
	p.server.RunCommand("kick " + p.Name())
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
	p.server.RunCommand(sprintf(rconPlayerDiscord, p.index, uid))
}

// DiscordUID returns a players Discord ID
func (p *Player) DiscordUID() string {
	return p.discordid
}

// Message messages a player
func (p *Player) Message(string) {
}

/*****************************/
/* IFace ifaces.ISteamPlayer */
/*****************************/

// SteamUID returns the steamUID or 0 of the player
// TODO: Modify this function, as well as the GameDB to store this data
// in the sqlite database
func (p *Player) SteamUID() int64 {
	if p.steam64 != 0 {
		return p.steam64
	}

	cmd := sprintf(steamUIDCommand, p.Index())
	out, err := p.server.RunCommand(cmd)
	if err != nil {
		return 0
	}

	m := steamUIDRegexp.FindStringSubmatch(out)
	if len(m) < 2 {
		return 0
	}

	sid, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0
	}

	logger.LogDebug(p, "Setting player steamcmd to: "+m[1])
	p.steam64 = sid
	return sid
}

/************************/
/* IFace logger.ILogger */
/************************/

// UUID returns the UUID of a player
func (p *Player) UUID() string {
	return sprintf("Player:%s:%s", p.index, p.name)
}

// Loglevel returns the loglevel of a player
func (p *Player) Loglevel() int {
	return p.loglevel
}

// SetLoglevel sets the loglevel of a player
func (p *Player) SetLoglevel(l int) {
	p.loglevel = l
}

/***************************/
/* IFace ifaces.IHaveShips */
/***************************/

// GetLastJumps returns up to (max) jumps that this player has
// performed recently
//
// TODO: This should return both the jumps and how many were found
// so that we can avoid an extra len call later
func (p *Player) GetLastJumps(limit int) []ifaces.ShipCoordData {
	var jumps []ifaces.ShipCoordData

	var l = len(p.jumphistory)
	var i = l - 1
	var n = 0

	if l == 0 {
		return jumps
	}

	// If -1 is used just return the entire history, but in reverse
	if limit < 0 {
		limit = l
	}

	for n < limit {
		if i < 0 {
			break
		}
		jumps = append(jumps, p.jumphistory[i])
		n++
		i--
	}

	return jumps
}

// SetJumpHistory sets the jump history for a player
func (p *Player) SetJumpHistory([]ifaces.ShipCoordData) {

}
