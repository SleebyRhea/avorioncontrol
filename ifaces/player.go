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

	logger.ILogger
	FactionID() string
}

// ITrackedPlayer defines an interface to an a player that has tracking
type ITrackedPlayer interface {
	IHaveShips
}

// INetPlayer describes a an interface to a player that can connect
//	over the internet
type INetPlayer interface {
	SetOnline(bool)
	SetIP(string)
	IP() net.IP
}

// IModeratablePlayer describes an interface to a player that can be
//	be moderated
type IModeratablePlayer interface {
	// DirectMessage(string)
	// Mail(string)
	Kick(string, func(interface{}) error, func(interface{}) error) error
	Ban(string, func(interface{}) error, func(interface{}) error) error
}

// IDiscordIntegratedPlayer describes an interface to a player that has
//	Discord integration
type IDiscordIntegratedPlayer interface {
	DiscordID() string
}

// ISteamPlayer describes an interface to a player that has a SteamID
type ISteamPlayer interface {
	Steam64ID() string
}
