package avorion

import (
	"AvorionControl/ifaces"
	"fmt"
)

// Alliance defines a player alliance in Avorion
type Alliance struct {
	index    string
	name     string
	leader   *Player
	members  []*Player
	loglevel int

	// alliance data
	resources map[string]int64
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

// UpdateCoords updates the coordinate DB of the Alliance
func (a *Alliance) UpdateCoords(ifaces.ShipCoordData) {
}

// Update updates the Alliance internal data
func (a *Alliance) Update() error {
	return nil
}

/************************/
/* IFace logger.ILogger */
/************************/

// UUID returns the UUID of an alliance
func (a *Alliance) UUID() string {
	return fmt.Sprintf("%s:%s", a.index, a.name)
}

// Loglevel returns the loglevel of an alliance
func (a *Alliance) Loglevel() int {
	return a.loglevel
}

// SetLoglevel sets the loglevel of an alliance
func (a *Alliance) SetLoglevel(l int) {
	a.loglevel = l
}
