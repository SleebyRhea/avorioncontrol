package events

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"fmt"
	"regexp"
	"strconv"
)

var discChatRe = regexp.MustCompile(`^\s*<D> <.*?#[0-9]{4}> (.*)$`)
var modURLBase = `https://steamcommunity.com/sharedfiles/filedetails/?id=`

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
		`^\s*playerJoinEvent: ([0-9]+) (.+?)\s*$`,
		handleEventPlayerJoin)

	New("EventPlayerLeft",
		`^\s*playerLeftEvent: ([0-9]+) (.+?)\s*$`,
		handleEventPlayerLeft)

	New("EventPlayerKick",
		`^\s*doPlayerKickEvent: ([0-9]+) (.*?)\s*$`,
		handleEventPlayerKick)

	New("EventPlayerBan",
		`^\s*doPlayerBanEvent: ([0-9]+) (.*?)\s*$`,
		handleEventPlayerBan)

	// New("EventDiscordIntegrationRequest",
	// 	`^\s*discordIntegrationRequestEvent: ([0-9]+) ([0-9]+)`,
	// 	handleDiscordIntegrationRequest)

	New("EventModUpdate",
		`^\s*Downloading ([0-9]+) \[[^\s]+ of [^\s]+ \| 100%\]\s*$`,
		handleModUpdate)
}

func handleEventPlayerJoin(e *Event, in string,
	gc ifaces.IGalaxyCache, cfg ifaces.IConfigurator,
	sendServ, sendChat, sendLog func(interface{}) error) {

	m := e.Capture.FindStringSubmatch(in)
	if p := gc.Players().FromFactionID(m[1]); p == nil {
		p, _ := gc.Players().NewPlayer(``, ``, ``, ``)
		p.SetOnline(true)
	} else {
		p.SetOnline(true)
	}
}

func handleEventPlayerLeft(e *Event, in string,
	gc ifaces.IGalaxyCache, cfg ifaces.IConfigurator,
	sendServ, sendChat, sendLog func(interface{}) error) {

	m := e.Capture.FindStringSubmatch(in)
	if p := gc.Players().FromFactionID(m[1]); p != nil {
		p.SetOnline(false)
		return
	}

	logger.LogError(e, "Player logged off, but has no tracking: "+m[2])
}

func handleEventShipTrackInit(e *Event, in string,
	gc ifaces.IGalaxyCache, cfg ifaces.IConfigurator,
	sendServ, sendChat, sendLog func(interface{}) error) {
}

func handleEventShipJump(e *Event, in string,
	gc ifaces.IGalaxyCache, cfg ifaces.IConfigurator,
	sendServ, sendChat, sendLog func(interface{}) error) {

	m := e.Capture.FindStringSubmatch(in)
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

func handlePlayerChat(e *Event, in string,
	gc ifaces.IGalaxyCache, cfg ifaces.IConfigurator,
	sendServ, sendChat, sendLog func(interface{}) error) {
	logger.LogChat(srv, in)
	// Catch our own discord messages
	if discChatRe.MatchString(in) {
		return
	}

	m := e.Capture.FindStringSubmatch(in)
	if m[1] != "Server" && m[1] != "Discord" {
		out := m[2]
		if len(out) >= 2000 {
			logger.LogInfo(e, "Truncated player message for sending")
			out = out[0:1900]
			out += "...(truncated)"
		}

		output := ifaces.ChatData{
			Name: m[1],
			Msg:  out}

		sendChat(output)
	}
}

func handleNilCommand(e *Event, in string,
	gc ifaces.IGalaxyCache, cfg ifaces.IConfigurator,
	sendServ, sendChat, sendLog func(interface{}) error) {
}

// func handleDiscordIntegrationRequest(cfg ifaces.IConfigurator, sendServ, sendChat, sendLog chan interface{}
//  e *Event, in string, oc chan string) {
// 	m := e.Capture.FindStringSubmatch(in)
// 	srv.AddIntegrationRequest(m[1], m[2])
// 	logger.LogInfo(srv, "Received Discord integration request")
// }

func handleEventPlayerKick(e *Event, in string,
	gc ifaces.IGalaxyCache, cfg ifaces.IConfigurator,
	sendServ, sendChat, sendLog func(interface{}) error) {

	m := e.Capture.FindStringSubmatch(in)
	p := srv.Player(m[1])

	// If the player cannot be found, we *do* still want to kick them, so just
	// run the ban and output an error
	if p == nil {
		logger.LogError(srv, fmt.Sprintf("Failed to locate player index: %s", m[1]))
		srv.RunCommand(fmt.Sprintf(`kick %s %s`, m[1], m[2]))
		return
	}

	p.Kick(m[2])
}

func handleEventPlayerBan(e *Event, in string,
	gc ifaces.IGalaxyCache, cfg ifaces.IConfigurator,
	sendServ, sendChat, sendLog func(interface{}) error) {

	m := e.Capture.FindStringSubmatch(in)
	p := gc.Players().FromFactionID(m[1])

	// If the player cannot be found, we *do* still want to ban them, so just
	// run the ban and output an error
	if p == nil {
		logger.LogError(e, fmt.Sprintf("Failed to locate player index: %s", m[1]))
		sendServ <- ifaces.RconCommand{
			Command:   `ban`,
			Arguments: []string{m[1], m[2]}}
	}

	// p.Ban(m[2])

	sendLog(ifaces.ChatData{
		Msg: fmt.Sprintf("**Kicked Player:** `%s`\n**Reason:** _%s_",
			p.Name(), m[2])})
}

func handleModUpdate(e *Event, in string,
	gc ifaces.IGalaxyCache, cfg ifaces.IConfigurator,
	sendServ, sendChat, sendLog func(interface{}) error) {

	m := e.Capture.FindStringSubmatch(in)
	out := fmt.Sprintf("Updated %s%s", modURLBase, m[1])

	output := ifaces.ChatData{
		Name: `Startup`,
		Msg:  out}

	sendLog <- output
}

func defaultEventHandler(cfg ifaces.IConfigurator, sendServ, sendChat, sendLog chan interface{},
	gc ifaces.IGalaxyCache, e *Event, in string, oc chan string) {
	logger.LogOutput(srv, in)
}
