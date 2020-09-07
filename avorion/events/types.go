package events

import (
	"AvorionControl/gameserver"
	"regexp"
)

// EventHandler is a function that parses a logfile and performs an action
type EventHandler func(gameserver.IServer, *Event, string, chan string)

// Event describes logged output that the gameserver.Server has output that can
// be acted upon in some fashion
//
// TODO: Have Event implement Logger
type Event struct {
	name     string
	loglevel int
	Capture  *regexp.Regexp
	Handler  EventHandler
}
