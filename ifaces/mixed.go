package ifaces

// IHaveShips describes an object that is capable of returning jumphistory of ships
type IHaveShips interface {
	Name() string
	GetLastJumps(int) []ShipCoordData
	SetJumpHistory([]ShipCoordData)
}
