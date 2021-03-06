package pubsub

import (
	"avorioncontrol/logger"
	"sync"
)

// RCONBUSID is the string that describes the MessageBus subscription ID for
// messaging the Avorion server
const (
	RCONBUSID     = `ServerRCON`
	DISCCHATBUSID = `DiscordChat`
	DISCLOGBUSID  = `DiscordLog`
)

// New returns a new MessageBus
func New(exit chan struct{}) *MessageBus {
	sub := &MessageBus{
		mutex: &sync.Mutex{},
		in:    make(map[string]chan interface{}, 0),
		out:   make(map[string][]chan interface{}, 0),
	}
	return sub
}

// MessageBus is a messaging bus interface to a generic type struct.
type MessageBus struct {
	loglevel int
	mutex    *sync.Mutex

	in  map[string]chan interface{}
	out map[string][]chan interface{}
}

// UUID returns the UUID of a pubbus.Subscription
func (bus *MessageBus) UUID() string {
	return "MessagingBus"
}

// Loglevel returns the loglevel of an avorion.Server
func (bus *MessageBus) Loglevel() int {
	return bus.loglevel
}

// SetLoglevel sets the loglevel of an avorion.Server
func (bus *MessageBus) SetLoglevel(l int) {
	bus.loglevel = l
}

// NewSubscription creates a new subscription and returns a channel with which
// to publish on, and a function to cancel the subscription
func (bus *MessageBus) NewSubscription(name string) (chan interface{}, func()) {
	if name == "" {
		logger.LogError(bus, "invalid subscription name provided (cannot be empty)")
		return nil, nil
	}

	w := make(chan interface{}, 10)
	if _, ok := bus.in[name]; ok {
		w = bus.in[name]
	} else {
		bus.in[name] = w
	}

	cancel := func() {
	}

	// Publish our data to our output channels
	go func() {
		for {
			data := <-w
			bus.mutex.Lock()
			for _, out := range bus.out[name] {
				out <- data
			}
			bus.mutex.Unlock()
		}
	}()

	return w, cancel
}

// Listen provides an interface to listen for a message, and a function
// to cancel the subscription to the messaging bus.
func (bus *MessageBus) Listen(name string) (chan interface{}, func()) {
	r := make(chan interface{}, 10)
	if _, ok := bus.out[name]; !ok {
		bus.out[name] = make([]chan interface{}, 0)
	}

	bus.out[name] = append(bus.out[name], r)

	cancel := func() {
		bus.mutex.Lock()
		defer bus.mutex.Unlock()

		if _, ok := bus.out[name]; ok {
			for i, ch := range bus.out[name] {
				if ch == r {
					bus.out[name][i] = bus.out[name][len(bus.out[name])-1]
					bus.out[name][len(bus.out[name])] = nil
					bus.out[name] = bus.out[name][:len(bus.out[name])-1]
					break
				}
			}
		}
	}

	return r, cancel
}
