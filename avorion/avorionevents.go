package avorion

import (
	"AvorionControl/gameserver"
	"AvorionControl/logger"
	"fmt"
	"strings"
)

func initEventsDefs() {
	gameserver.Init()

	ipReString := "[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}"

	gameserver.RegisterEventHandler("EventConnection",
		"^("+ipReString+"):[0-9]{1,5} is connecting...$", handleEventConnection)

	gameserver.RegisterEventHandler("EventPlayerJoin",
		"^Player logged in: (.+), index: ([0-9]+)$", handleEventPlayerJoin)

	gameserver.RegisterEventHandler("EventPlayerLeft",
		"^<Server> (.+) left the galaxy$", handleEventPlayerLeft)

	gameserver.RegisterEventHandler("EventServerLag",
		"^Server frame took over [0-9]+ seconds?.$", handleEventServerLag)

	gameserver.RegisterEventHandler("EventNone",
		".*", defaultEventHandler)
}

func handleEventConnection(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	go func() { oc <- m[1] }()
}

func handleEventPlayerJoin(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	logger.LogInfo(gs, in, gs.WSOutput())
	m := e.Capture.FindStringSubmatch(in)
	gs.RunCommand("playerinfo -p " + m[1] + " -a -c -t -s")
}

func handleEventServerLag(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	logger.LogWarning(gs, in, gs.WSOutput())
}

func handleEventPlayerLeft(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	name := strings.TrimPrefix(in, "<Server> ")
	name = strings.TrimSuffix(name, " left the galaxy")
	gs.RemovePlayer(name)
}

func handleEventPlayerBoot(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	logger.LogInfo(gs, fmt.Sprintf("Failed connection: %s [%s]", m[1], m[2]),
		gs.WSOutput())
}

func handleEventPlayerBan(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	logger.LogInfo(gs, in, gs.WSOutput())
}

func handleEventServerPass(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	m := e.Capture.FindStringSubmatch(in)
	gs.SetPassword(m[1])
}

func defaultEventHandler(gs gameserver.Server, e *gameserver.Event, in string,
	oc chan string) {
	logger.LogOutput(gs, in)
}
