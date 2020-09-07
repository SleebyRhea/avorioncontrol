package events

import "AvorionControl/gameserver"

// GetFromString returns a reference to a game event given a matching string
func GetFromString(in string) *Event {
	for _, e := range events {
		if e.Capture.MatchString(in) {
			return e
		}
	}
	return nil
}

// EventType - Given string and its source Server{}, determine the event type
// that was provided. Returns a -1 if none was found.
func EventType(s string, srv gameserver.IServer) int {
	for i, ge := range events {
		if ge.Capture.MatchString(s) {
			return i
		}
	}
	return -1
}

func init() {
	events = make([]*Event, 0)
	eventsMap = make(map[string]*Event)
	initB()
}
