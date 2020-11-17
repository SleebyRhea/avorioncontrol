package ifaces

import (
	"avorioncontrol/logger"
	"net"
)

// IPlayer describes a an IGameServer player
type IPlayer interface {
	IDiscordIntegratedPlayer
	IModeratablePlayer
	ITrackedPlayer
	ISteamPlayer

	INetPlayer
}

// ITrackedPlayer defines an interface to an a player that has tracking
type ITrackedPlayer interface {
	logger.ILogger
	IHaveShips
	Index() string
	Message(string)
	AddJump(ShipCoordData)

	Update() error
	UpdateFromData([15]string) error
}

// INetPlayer describes a an interface to a player that can connect
//	over the internet
type INetPlayer interface {
	IP() net.IP
	SetIP(string)

	Online() bool
	SetOnline(bool)
}

// IModeratablePlayer describes an interface to a player that can be
//	be moderated
type IModeratablePlayer interface {
	Kick(string)
	Ban(string)
}

// IDiscordIntegratedPlayer describes an interface to a player that has
//	Discord integration
type IDiscordIntegratedPlayer interface {
	DiscordUID() string
	SetDiscordUID(string)
}

// ISteamPlayer describes an interface to a player that has a SteamID
type ISteamPlayer interface {
	SteamUID() int64
}
