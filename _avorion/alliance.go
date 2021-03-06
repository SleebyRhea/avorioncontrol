package avorion

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"fmt"
	"strconv"
	"time"
)

// Alliance defines a player alliance in Avorion
type Alliance struct {
	index    string
	name     string
	leader   *Player
	members  []*Player
	loglevel int
	server   *Server

	// alliance data
	resources   map[string]int64
	jumphistory []ifaces.ShipCoordData
}

// Message sends an in-game message to all members of an alliance
func (a *Alliance) Message(string) {
}

// Index returns the faction index of an alliance
func (a *Alliance) Index() string {
	return a.index
}

// Name returns the name of the Alliance
func (a *Alliance) Name() string {
	return a.name
}

// Update updates the Alliance internal data
func (a *Alliance) Update() error {
	return nil
}

// UpdateFromData updates the alliances information using the data from
//	a successful reAllianceData match
func (a *Alliance) UpdateFromData(d [13]string) error {
	return nil
}

// AddJump registers a jump that a player took into a system
func (a *Alliance) AddJump(sc ifaces.ShipCoordData) {
	sc.Time = time.Now()
	a.jumphistory = append(a.jumphistory, sc)
	if len(a.jumphistory) > 1000 {
		a.jumphistory = a.jumphistory[1:]
	}

	fid64, _ := strconv.ParseInt(a.index, 10, 32)
	fid := int(fid64)
	s := a.server.Sector(sc.X, sc.Y)

	// Add a pointer to the players jump to the sector history for our usage
	//	later on
	jump := &ifaces.JumpInfo{
		Name: sc.Name,
		FID:  fid,
		Kind: "alliance",
		Time: sc.Time,
		X:    sc.X,
		Y:    sc.Y}

	s.Jumphistory = append(s.Jumphistory, jump)

	id, _ := strconv.Atoi(a.Index())
	a.server.tracking.AddJump(s.Index, int64(id), 1, *jump)

	logger.LogDebug(a, "Updated jumphistory")
}

/************************/
/* IFace logger.ILogger */
/************************/

// UUID returns the UUID of an alliance
func (a *Alliance) UUID() string {
	return fmt.Sprintf("Alliance:%s:%s", a.index, a.name)
}

// Loglevel returns the loglevel of an alliance
func (a *Alliance) Loglevel() int {
	return a.loglevel
}

// SetLoglevel sets the loglevel of an alliance
func (a *Alliance) SetLoglevel(l int) {
	a.loglevel = l
}

/***************************/
/* IFace ifaces.IHaveShips */
/***************************/

// GetLastJumps returns up to (max) jumps that this player has performed recently
// TODO: This should return both the jumps and how many were found
func (a *Alliance) GetLastJumps(limit int) []ifaces.ShipCoordData {
	var jumps []ifaces.ShipCoordData

	var l = len(a.jumphistory)
	var i = l - 1
	var n = 0

	if l == 0 {
		return jumps
	}

	// If -1 is used just return the entire history, but in reverse (for easy search)
	if limit < 0 {
		limit = l
	}

	for n < limit {
		if i < 0 {
			break
		}
		jumps = append(jumps, a.jumphistory[i])
		n++
		i--
	}

	return jumps
}

// SetJumpHistory sets the jump history for an alliance
func (a *Alliance) SetJumpHistory([]ifaces.ShipCoordData) {

}
