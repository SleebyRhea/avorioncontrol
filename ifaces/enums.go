package ifaces

// enums describing the status of a server
const (
	ServerOffline    = 0
	ServerOnline     = 1
	ServerStarting   = 2
	ServerStopping   = 3
	ServerRestarting = 4
	ServerCrashed    = 255
)
