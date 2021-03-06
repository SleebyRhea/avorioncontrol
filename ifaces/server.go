package ifaces

import (
	"avorioncontrol/logger"
	"context"
)

// IGameServer describes an interface to a server with full capability
type IGameServer interface {
	IServer
	ILockableServer
	IVersionedServer
	logger.ILogger
}

// IServer defines an interface to a Gameserver that we can change the status
//	of and log
type IServer interface {
	IsUp() bool
	Stop(context.Context, IConfigurator, IGalaxyCache, IPlayerCache) error
	InitializeEvents(IConfigurator)
}

// IVersionedServer describes an interface to an IGameserver's version information
type IVersionedServer interface {
	Version() string
}

// ILockableServer describes an interface to lock a server with a password
type ILockableServer interface {
	Password() string
}
