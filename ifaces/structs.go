package ifaces

import (
	"regexp"
	"time"
)

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

// ServerStatus is a struct that describes the current status of an
//	AvorionServer
type ServerStatus struct {
	Name    string
	Output  string
	Players string

	Status        int
	TotalPlayers  int
	PlayersOnline int
	Alliances     int
	Sectors       int

	INI *ServerGameConfig
}

// ServerGameConfig describes an object that contains the current
//	game configuration
type ServerGameConfig struct {
	PVP                 bool
	Name                string
	Collision           string
	Seed                string
	Steam               bool
	Version             string
	Difficulty          int
	BlockLimit          int64
	VolumeLimit         int64
	MaxPlayerShips      int64
	MaxPlayerSlots      int64
	MaxPlayerStations   int64
	MaxAllianceSlots    int64
	MaxAllianceShips    int64
	MaxAllianceStations int64
}

// LoggedServerEvent describes an event that can be tracked and logged
type LoggedServerEvent struct {
	Name    string
	FString string
	Regex   *regexp.Regexp
}
