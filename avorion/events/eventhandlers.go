package events

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"fmt"
	"regexp"
	"strconv"
)

var discChatRe = regexp.MustCompile(`^\s*<D> <.*?#[0-9]{4}> (.*)$`)

func initB() {
	New("EventShipTrackInit",
		`^\s*shipTrackInitEvent: (-?[0-9]+) (-?[0-9]+):(-?[0-9]+) (.*)$`,
		handleEventShipTrackInit)

	New("EventPlayerChat",
		`^\s*<(.+?)> (.*)`,
		handlePlayerChat)

	New("EventShipJump",
		`^\s*shipJumpEvent: (-?[0-9]+) (-?[0-9]+):(-?[0-9]+) (.*)$`,
		handleEventShipJump)

	New("EventPlayerJoin",
		`^\s*Player logged in: (.+?), index: ([0-9]+)\s*$`,
		handleEventPlayerJoin)

	New("EventPlayerLeft",
		`^\s*Player logged off: (.+?), index: ([0-9]+):?\s*$`,
		handleEventPlayerLeft)

	New("EventServerLag",
		`^\s*Server frame took over [0-9]+ seconds?\.?\s*$`,
		handleEventServerLag)

	New("EventDiscordIntegrationRequest",
		`^\s*discordIntegrationRequestEvent: ([0-9]+) ([0-9]+)`,
		handleDiscordIntegrationRequest)

	New("NilCommandEvent",
		`^\s*nilCommandEvent: (.*)$`,
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
		srv.NewPlayer(m[2], m)
	} else {
		p.SetOnline(true)
		p.Update()
	}

	srv.SendChat(ifaces.ChatData{
		Msg:  fmt.Sprintf("Player %s has logged in", m[1]),
		Name: "Avorion"})
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
		p.AddJump(data)
	} else if a := srv.Alliance(m[1]); a != nil {
		a.AddJump(data)
	}
}

func handleEventPlayerLeft(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	logger.LogOutput(srv, in)

	srv.SendChat(ifaces.ChatData{
		Msg:  fmt.Sprintf("Player %s has logged off", m[1]),
		Name: "Avorion"})

	if p := srv.Player(m[2]); p != nil {
		p.SetOnline(false)
		return
	}

	logger.LogError(srv, "Player logged off, but has no tracking: "+m[2])
}

func handlePlayerChat(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogOutput(srv, in)
	// Catch our own discord messages
	if discChatRe.MatchString(in) {
		return
	}

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

func handleNilCommand(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
}

func handleEventServerLag(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogWarning(srv, in)
}

func handleDiscordIntegrationRequest(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	srv.AddIntegrationRequest(m[1], m[2])
}

func defaultEventHandler(srv ifaces.IGameServer, e *Event, in string,
	oc chan string) {
	logger.LogOutput(srv, in)
}
