package events

import (
	"regexp"
)

var events []*Event             // Iteration
var eventsMap map[string]*Event // Reference

// New makes,registers, and returns a game event using the Regex that is used
// to detect it. Panic if the Regexp that was provided matches the Regexp of a
// previously registered Event. Returns the index of the registered event
func New(n, re string, h EventHandler) *Event {
	for _, e := range events {
		if e.Capture.String() == re || e.name == n {
			panic("Cannot register the same event multiple times")
		}
	}

	ge := &Event{
		name:     n,
		loglevel: 3,
		Capture:  regexp.MustCompile(re),
		Handler:  h}

	events = append(events, ge)
	eventsMap[n] = ge

	return ge
}

// Name returns the name of the Event
func (e *Event) Name() string {
	return e.name
}

/************************/
/* IFace logger.ILogger */
/************************/

// UUID returns an events UUID
func (e *Event) UUID() string {
	return e.name
}

// Loglevel returns an events loglevel
func (e *Event) Loglevel() int {
	return e.loglevel
}

// SetLoglevel sets an events loglevel
func (e *Event) SetLoglevel(l int) {
	e.loglevel = l
}
