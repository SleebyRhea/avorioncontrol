package main

import "regexp"

var gameEvents []*GameEvent             // Iteration
var gameEventsMap map[string]*GameEvent // Reference

func initEvents() {
	gameEvents = make([]*GameEvent, 0)
	gameEventsMap = make(map[string]*GameEvent)
}

// gameEventHandler - A function that takes a GameServer, some output and output
// channel (used as needed) and processes a given gameevent
type gameEventHandler func(GameServer, *GameEvent, string, chan string)

// GameEvent -
// TODO: Have GameEvent implement Loggable
type GameEvent struct {
	name    string
	Capture *regexp.Regexp
	Handler gameEventHandler
}

// GetEventFromString -
func GetEventFromString(in string) *GameEvent {
	for _, e := range gameEvents {
		if e.Capture.MatchString(in) {
			return e
		}
	}

	return gameEvents[len(gameEvents)-1]
}

// RegisterGameEventHandler - Register a game event using the Regex that is used to detect
// it. Panic if the Regexp that was provided matches the Regexp of a previously
// registered GameEvent. Returns the index of the registered event
func RegisterGameEventHandler(n, re string, f gameEventHandler) int {
	for _, e := range gameEvents {
		if e.Capture.String() == re || e.name == n {
			panic("Cannot register the same event multiple times")
		}
	}

	ge := &GameEvent{
		name:    n,
		Capture: regexp.MustCompile(re),
		Handler: f}

	gameEvents = append(gameEvents, ge)
	gameEventsMap[n] = ge

	return len(gameEvents)
}

// EventType - Given string and its source GameServer{}, determine the event type
// that was provided. Returns a -1 if none was found.
func EventType(s string, gs GameServer) int {
	for i, ge := range gameEvents {
		if ge.Capture.MatchString(s) {
			return i
		}
	}
	return -1
}

func defaultEventHandler(gs GameServer, e *GameEvent, in string, oc chan string) {
	LogOutput(gs, in)
}
