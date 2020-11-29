package events

import (
	"errors"
	"regexp"
	"time"
)

var events []*Event             // Iteration
var eventsMap map[string]*Event // Reference
var benchTimer *time.Timer

// Initialize handles event initialization and should be run before the
// server starts
func Initialize() {
	// For re-init
	events = nil
	eventsMap = nil
	benchTimer = nil

	events = make([]*Event, 0)
	eventsMap = make(map[string]*Event)
	benchTimer = time.NewTimer(time.Minute * 5)
	initB()
}

// New makes,registers, and returns a game event using the Regex that is used
// to detect it. Panic if the Regexp that was provided matches the Regexp of a
// previously registered Event. Returns the index of the registered event
func New(n, re string, h EventHandler) (*Event, error) {
	for _, e := range events {
		if e.Capture.String() == re || e.name == n {
			return nil, errors.New("Cannot register the same event multiple times")
		}
	}

	ge := &Event{
		name:     n,
		loglevel: 3,
		Capture:  regexp.MustCompile(re),
		Handler:  h}

	events = append(events, ge)
	eventsMap[n] = ge

	return ge, nil
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
