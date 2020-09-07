package avorion

import "AvorionControl/gameserver"

func init() {
	var _ gameserver.IPlayer = (*Player)(nil)
	var _ gameserver.IServer = (*Server)(nil)
}
