package gameserver

import (
	"regexp"
)

// EventHandler - A function that takes a Server, some output and output
// channel (used as needed) and processes a given Event
type EventHandler func(Server, *Event, string, chan string)

var events []*Event             // Iteration
var eventsMap map[string]*Event // Reference

// InitEvents initalizes the local Event map and slice that are used for tracking
func InitEvents() {
	if events == nil {
		events = make([]*Event, 0)
	}

	if eventsMap == nil {
		eventsMap = make(map[string]*Event)
	}
}

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

// Name returns the name of the Event
func (ge *Event) Name() string {
	return ge.name
}

/************/
/* Logger */
/************/

// UUID -
func (ge *Event) UUID() string {
	return ge.name
}

// Loglevel -
func (ge *Event) Loglevel() int {
	return ge.loglevel
}

// SetLoglevel -
func (ge *Event) SetLoglevel(l int) {
	ge.loglevel = l
}

// GetEventFromString returns a reference to a game event given a matching string
func GetEventFromString(in string) *Event {
	for _, e := range events {
		if e.Capture.MatchString(in) {
			return e
		}
	}

	return nil
}

// RegisterEventHandler - Register a game event using the Regex that is used to detect
// it. Panic if the Regexp that was provided matches the Regexp of a previously
// registered Event. Returns the index of the registered event
func RegisterEventHandler(n, re string, f EventHandler) int {
	for _, e := range events {
		if e.Capture.String() == re || e.name == n {
			panic("Cannot register the same event multiple times")
		}
	}

	ge := &Event{
		name:     n,
		Capture:  regexp.MustCompile(re),
		Handler:  f,
		loglevel: 3}

	events = append(events, ge)
	eventsMap[n] = ge

	return len(events)
}

// EventType - Given string and its source Server{}, determine the event type
// that was provided. Returns a -1 if none was found.
func EventType(s string, gs Server) int {
	for i, ge := range events {
		if ge.Capture.MatchString(s) {
			return i
		}
	}
	return -1
}
