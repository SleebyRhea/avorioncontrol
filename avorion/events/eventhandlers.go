package events

import (
	"AvorionControl/ifaces"
	"AvorionControl/logger"
	"fmt"
	"strconv"
)

func initB() {
	New("EventShipTrackInit",
		"^\\s*shipTrackInitEvent: (-?[0-9]+) (-?[0-9]+):(-?[0-9]+) (.*)$",
		handleEventShipTrackInit)

	New("EventPlayerChat",
		"^\\s*<(.+?)> (.*)",
		handlePlayerChat)

	New("EventShipJump",
		"^\\s*shipJumpEvent: (-?[0-9]+) (-?[0-9]+):(-?[0-9]+) (.*)$",
		handleEventShipJump)

	New("EventPlayerJoin",
		"^\\s*Player logged in: (.+?), index: ([0-9]+)\\s*$",
		handleEventPlayerJoin)

	New("EventPlayerLeft",
		"^\\s*Player logged off: (.+?), index: ([0-9]+):?\\s*$",
		handleEventPlayerLeft)

	New("EventServerLag",
		"^\\s*Server frame took over [0-9]+ seconds?\\.?\\s*$",
		handleEventServerLag)

	New("NilCommandEvent",
		"^\\s*nilCommandEvent: (.*)$",
		handleNilCommand)

	New("EventNone",
		".*",
		defaultEventHandler)
}

func handleEventConnection(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	go func() { oc <- m[1] }()
}

func handleEventPlayerJoin(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	logger.LogOutput(srv, in)
	if p := srv.Player(m[2]); p == nil {
		srv.NewPlayer(m[2], in)
	} else {
		p.SetOnline(true)
		p.Update()
	}

	srv.SendChat(ifaces.ChatData{
		Msg:  fmt.Sprintf("Player %s has logged in", m[1]),
		Name: "Server"})
}

func handleEventShipTrackInit(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
}

func handleEventShipJump(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)

	// We already use a regex to make sure we capture the correct values
	x, _ := strconv.Atoi(m[2])
	y, _ := strconv.Atoi(m[3])
	n := m[4]

	data := ifaces.ShipCoordData{X: x, Y: y, Name: n}

	if p := srv.Player(m[1]); p != nil {
		p.UpdateCoords(data)
	} else if a := srv.Alliance(m[1]); a != nil {
		a.UpdateCoords(data)
	}
}

func handleEventPlayerLeft(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	logger.LogOutput(srv, in)

	srv.SendChat(ifaces.ChatData{
		Msg:  fmt.Sprintf("Player %s has logged off", m[1]),
		Name: "Server"})

	if p := srv.Player(m[2]); p != nil {
		p.SetOnline(false)
		return
	}

	logger.LogError(srv, "Player logged off, but has no tracking: "+m[2])
}

func handlePlayerChat(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogOutput(srv, in)
	m := e.Capture.FindStringSubmatch(in)
	if m[1] != "Server" && m[1] != "D" {
		out := m[2]
		if len(out) >= 2000 {
			logger.LogInfo(srv, "Truncated player message for sending")
			out = out[0:1900]
			out += "...(truncated)"
		}

		output := ifaces.ChatData{
			Name: m[1],
			Msg:  out}

		srv.SendChat(output)
	}
}

func handleNilCommand(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
}

func handleEventServerLag(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogWarning(srv, in)
}

func defaultEventHandler(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogOutput(srv, in)
}
