package avorion

import "avorioncontrol/ifaces"

func init() {
	var _ ifaces.IPlayer = (*Player)(nil)
	var _ ifaces.IGameServer = (*Server)(nil)
}
