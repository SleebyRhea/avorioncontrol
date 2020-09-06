package avorion

import (
	"AvorionControl/gameserver"
	"AvorionControl/logger"
)

func init() {
	var reg = gameserver.RegisterEventHandler

	gameserver.InitEvents()

	reg("EventPlayerJoin", "^\\s*Player logged in: (.+?), index: ([0-9]+)\\s*$",
		handleEventPlayerJoin)

	reg("EventPlayerLeft", "^\\s*Player logged off: (.+?), index: ([0-9]+):?\\s*$",
		handleEventPlayerLeft)

	reg("EventServerLag", "^\\s*Server frame took over [0-9]+ seconds?\\.?\\s*$",
		handleEventServerLag)

	reg("EventNone", ".*", defaultEventHandler)
}

func handleEventConnection(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	go func() { oc <- m[1] }()
}

func handleEventPlayerJoin(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	var (
		p     gameserver.Player
		m     = e.Capture.FindStringSubmatch(in)
		index = m[2]
	)

	logger.LogOutput(gs, in, gs.WSOutput())
	if p = gs.Player(index); p == nil {

		gs.NewPlayer(index, in)
		return
	}

	p.SetOnline(true)
	p.GetData()
}

func handleEventPlayerLeft(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	logger.LogOutput(gs, in, gs.WSOutput())

	if p := gs.Player(m[2]); p != nil {
		p.SetOnline(false)
		return
	}

	logger.LogError(gs, "Player logged off, but has no tracking: "+m[2])
}

func handleEventServerLag(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	logger.LogWarning(gs, in, gs.WSOutput())
}

func defaultEventHandler(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	logger.LogOutput(gs, in)
}
