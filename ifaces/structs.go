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
	Jump *ShipCoordData
	Name string
	Kind string
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
	X int
	Y int

	// slice of pointers to player jumpdata structs
	Jumphistory []*JumpInfo
}
