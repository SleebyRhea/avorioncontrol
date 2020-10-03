package ifaces

import "time"

// ChatData describes datapassed between Discord and the Server
type ChatData struct {
	Name string
	UID  string
	Msg  string
}

// JumpInfo describes a ship jump
type JumpInfo struct {
	// Jump *ShipCoordData
	Time time.Time
	Name string
	Kind string
	FID  int
	X    int
	Y    int
}

// ShipCoordData describes a set of coords for a ship
type ShipCoordData struct {
	X    int
	Y    int
	Name string
	Time time.Time
}

// Sector defines a sector in an Avorion galaxy
type Sector struct {
	Index int64
	X     int
	Y     int

	Jumphistory []*JumpInfo
}
