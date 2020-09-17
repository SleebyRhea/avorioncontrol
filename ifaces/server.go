package ifaces

import (
	"avorioncontrol/logger"
)

// IGameServer describes an interface to a server with full capability
type IGameServer interface {
	IServer
	IMOTDServer
	ISeededServer
	IGalaxyServer
	ILockableServer
	IPlayableServer
	IVersionedServer
	ICommandableServer
	IDiscordIntegratedServer
}

// IServer defines an interface to a Gameserver that we can change the status
//	of and log
type IServer interface {
	IsUp() bool
	Stop(bool) error
	Start(bool) error
	Restart() error
	Config() IConfigurator

	logger.ILogger
}

// IGalaxyServer describes an interface to a server with a sectored galaxy
type IGalaxyServer interface {
	Sector(int, int) *Sector
}

// IPlayableServer defines an object that can track the players that have joined
type IPlayableServer interface {
	Players() []IPlayer
	RemovePlayer(string)
	NewPlayer(string, []string) IPlayer

	Player(string) IPlayer
	PlayerFromName(string) IPlayer
	PlayerFromDiscord(string) IPlayer

	Alliance(string) IAlliance
	Alliances() []IAlliance
	NewAlliance(string, []string) IAlliance
}

// IVersionedServer describes an interface to an IGameserver's version information
type IVersionedServer interface {
	Version() string
	SetVersion(string)
}

// ILockableServer describes an interface to lock a server with a password
type ILockableServer interface {
	Password() string
	SetPassword(string)
}

// ICommandableServer describes an interface to an IGameServer that can have run
//	game commands
type ICommandableServer interface {
	RunCommand(string) (string, error)
}

// IMOTDServer describes an interface to a server that can set an MOTD
type IMOTDServer interface {
	MOTD() string
	SetMOTD(string)
}

// ISeededServer is an interface to an object that has a seed
type ISeededServer interface {
	Seed() string
	SetSeed(string)
}

// IDiscordIntegratedServer describes an IGameServer that is capable of integrating
//	with Discord
type IDiscordIntegratedServer interface {
	AddIntegrationRequest(string, string)
	ValidateIntegrationPin(string, string) bool
	SendChat(ChatData)
}
