package events

import (
	"AvorionControl/ifaces"
	"AvorionControl/logger"
)

func initB() {
	New("EventPlayerChat",
		"^\\s*<(.+?)> (.*)",
		handlePlayerChat)

	New("EventPlayerJoin",
		"^\\s*Player logged in: (.+?), index: ([0-9]+)\\s*$",
		handleEventPlayerJoin)

	New("EventPlayerLeft",
		"^\\s*Player logged off: (.+?), index: ([0-9]+):?\\s*$",
		handleEventPlayerLeft)

	New("EventServerLag",
		"^\\s*Server frame took over [0-9]+ seconds?\\.?\\s*$",
		handleEventServerLag)

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
}

func handleEventPlayerLeft(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	logger.LogOutput(srv, in)

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
	if m[1] != "Server" {
		out := m[2]
		if len(out) >= 2000 {
			logger.LogInfo(srv, "Truncated player message for sending")
			out = out[0:1900]
			out += "...(truncated)"
		}

		if p := srv.PlayerFromName(m[1]); p != nil {
			srv.DCOutput() <- ifaces.ChatData{
				Name: m[1],
				UID:  p.DiscordUID(),
				Msg:  out}
		} else {
			logger.LogWarning(srv, "Unable to locate player: "+m[1])
			srv.DCOutput() <- ifaces.ChatData{
				Name: m[1],
				Msg:  out}
		}
	}
}

func handleEventServerLag(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogWarning(srv, in)
}

func defaultEventHandler(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogOutput(srv, in)
}
