package avorion

import "AvorionControl/ifaces"

func init() {
	var _ ifaces.IPlayer = (*Player)(nil)
	var _ ifaces.IGameServer = (*Server)(nil)
}
