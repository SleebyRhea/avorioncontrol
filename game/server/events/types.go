package events

import (
	"avorioncontrol/ifaces"
	"regexp"
)

// EventHandler is a function that parses a logfile and performs an action
type EventHandler func(*Event, string,
	ifaces.IGalaxyCache,
	ifaces.IConfigurator,
	func(interface{}) error, // RCON publisher
	func(interface{}) error, // Discord chat publisher
	func(interface{}) error) // Discord log publisher

// Event describes logged output that the ifaces.Server has output that can
// be acted upon in some fashion
type Event struct {
	FString  string
	name     string
	loglevel int
	Capture  *regexp.Regexp
	Handler  EventHandler
}
