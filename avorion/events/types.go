package events

import (
	"avorioncontrol/ifaces"
	"regexp"
)

// EventHandler is a function that parses a logfile and performs an action
type EventHandler func(ifaces.IGameServer, *Event, string, chan string)

// Event describes logged output that the ifaces.Server has output that can
// be acted upon in some fashion
//
// TODO: Have Event implement Logger
type Event struct {
	name     string
	loglevel int
	Capture  *regexp.Regexp
	Handler  EventHandler
}
