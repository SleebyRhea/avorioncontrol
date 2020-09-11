package events

import (
	"AvorionControl/ifaces"
	"AvorionControl/logger"
	"fmt"
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

	srv.SendChat(ifaces.ChatData{
		Msg:  fmt.Sprintf("Player %s has logged in", m[1]),
		Name: "Server"})
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
	if m[1] != "Server" {
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

func handleEventServerLag(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogWarning(srv, in)
}

func defaultEventHandler(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogOutput(srv, in)
}
