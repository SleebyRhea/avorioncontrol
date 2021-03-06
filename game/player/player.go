package player

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"fmt"
	"net"
	"strconv"
)

const (
	kickString = "**Kicked Player:** `%s`\n**Steam64:** `%s`\n**Reason:** _%s_"
	banString  = "**Banned Player:** `%s`\n**Steam64:** `%s`\n**Reason:** _%s_"
)

func init() {
	var _ ifaces.IPlayer = (*Player)(nil)
}

// Player is a player in Avorion that is being tracked
type Player struct {
	name      string
	steam64id string
	discordid string
	factionid string

	// State tracking
	isonline bool
	ip       net.IP

	// Logging
	loglevel int
}

// Name returns the players current name or an empty string
func (p *Player) Name() string {
	return p.name
}

// FactionID returns a string containing the players faction ID
func (p *Player) FactionID() string {
	return p.factionid
}

// Steam64ID returns a string containing the players steam64 ID
func (p *Player) Steam64ID() string {
	return p.steam64id
}

// DiscordID returns a string containing the players Discord ID
func (p *Player) DiscordID() string {
	return p.discordid
}

// SetDiscordID sets the players Discord ID
// 	Deprecated: This is a stub, and will be removed
func (p *Player) SetDiscordID(did string) error {
	return nil
}

// Online retrieves the players current online status
func (p *Player) Online() bool {
	return p.isonline
}

// SetOnline sets the players current online status
func (p *Player) SetOnline(on bool) {
	p.isonline = on
}

// IP returns the players current IP address (if they are logged in)
func (p *Player) IP() net.IP {
	if p.Online() {
		return p.ip
	}
	return nil
}

// SetIP sets the players current IP address.
func (p *Player) SetIP(ip string) {
	p.ip = net.ParseIP(ip)
}

// UUID returns the UUID of a player
func (p *Player) UUID() string {
	return `Player:` + p.steam64id
}

// Loglevel returns the loglevel of a player
func (p *Player) Loglevel() int {
	return p.loglevel
}

// SetLoglevel is a stub that implements ILogger
func (p *Player) SetLoglevel(l int) {
	logger.LogInfo(p, fmt.Sprintf("Setting logging level to %d", l))
	p.loglevel = l
}

// SetName sets the players name
// 	Deprecated: This is a stub, and will be removed
func (p *Player) SetName(_ string) {
	return
}

// Kick kicks a player from the server and sends a notification
// 	Deprecated: This is a stub, and will be removed
func (p *Player) Kick(reason string) {
	logger.LogWarning(p, "Kick() is a stub, please use IGameServer.KickPlayer")
	return
}

// Ban bans a player from the server and sends a notification
// 	Deprecated: This is a stub, and will be removed
func (p *Player) Ban(reason string) {
	logger.LogWarning(p, "Kick() is a stub, please use IGameServer.BanPlayer")
}

// Message sends a private message to a player
// 	Deprecated: This is a stub, and will be removed
func (p *Player) Message(string) {
	logger.LogWarning(p, "call to decprecated method Method()")
}

// Index returns the faction index for the player
// 	Deprecated: User FactionID() instead
func (p *Player) Index() string {
	logger.LogWarning(p, "call to decprecated method Index()")
	return p.FactionID()
}

// SteamUID returns the players steam64
// 	Deprecated: Use Steam64ID() instead
func (p *Player) SteamUID() int64 {
	logger.LogWarning(p, "call to decprecated method SteamUID()")
	sid, _ := strconv.ParseInt(p.steam64id, 10, 64)
	return sid
}

// DiscordUID returns the players Discord ID
// 	Deprecated: Use DiscordID() instead
func (p *Player) DiscordUID() string {
	logger.LogWarning(p, "call to decprecated method DiscordUID()")
	return p.DiscordID()
}

// SetDiscordUID sets the Discord ID for a player
// 	Deprecated: Use SetDiscordID() instead
func (p *Player) SetDiscordUID(did string) {
	logger.LogWarning(p, "call to decprecated method SetDiscordUID()")
	p.SetDiscordID(did)
}

// Update gets player data from the server and updates the player
// 	Deprecated: This is a stub and cannot not be used
func (p *Player) Update() error {
	logger.LogWarning(p, "call to decprecated method Update()")
	return nil
}

// UpdateFromData updates the players information using the data
// from a successful login
// 	Deprecated: This is a stub and cannot not be used
func (p *Player) UpdateFromData(d [15]string) error {
	logger.LogWarning(p, "call to decprecated method UpdateFromData()")
	return nil
}

// SetJumpHistory sets the jump history for a player
//	Deprecated: This is a stub, and cannot be used
func (p *Player) SetJumpHistory(jh []ifaces.ShipCoordData) {
	return
}

// GetLastJumps returns the last N jumps for a player
// 	Deprecated: This is a stub, and cannot be used
func (p *Player) GetLastJumps(n int) []ifaces.ShipCoordData {
	return nil
}

// AddJump adds a jump to a player
// 	Deprecated: This is a stub, and cannot be used
func (p *Player) AddJump(j ifaces.ShipCoordData) {
	return
}
