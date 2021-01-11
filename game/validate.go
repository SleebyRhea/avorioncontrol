package game

import (
	"avorioncontrol/game/player"
	"avorioncontrol/game/server"
	"avorioncontrol/ifaces"
)

func init() {
	var _ ifaces.IPlayer = (*player.Player)(nil)
	var _ ifaces.IAlliance = (*alliance.Alliance)(nil)
	var _ ifaces.IGameServer = (*server.Server)(nil)
}
