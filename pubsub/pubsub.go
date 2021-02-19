package pubsub

import (
	"avorioncontrol/logger"
	"sync"
)

// New returns a new Subscription
func New(exit chan struct{}) *Subscription {
	sub := &Subscription{
		mutex: &sync.Mutex{},
		in:    make(map[string]chan interface{}, 0),
		out:   make(map[string][]chan interface{}, 0),
	}
	return sub
}

// Subscription is a messaging bus interface to a generic type struct.
type Subscription struct {
	loglevel int
	mutex    *sync.Mutex

	in  map[string]chan interface{}
	out map[string][]chan interface{}
}

// UUID returns the UUID of a pubsub.Subscription
func (sub *Subscription) UUID() string {
	return "MessagingBus"
}

// Loglevel returns the loglevel of an avorion.Server
func (sub *Subscription) Loglevel() int {
	return sub.loglevel
}

// SetLoglevel sets the loglevel of an avorion.Server
func (sub *Subscription) SetLoglevel(l int) {
	sub.loglevel = l
}

// NewSubscription creates a new subscription and returns a channel with which
// to publish on, and a function to cancel the subscription
func (sub *Subscription) NewSubscription(name string) (chan interface{}, func()) {
	if name == "" {
		logger.LogError(sub, "invalid subscription name provided (cannot be empty)")
		return nil, nil
	}

	w := make(chan interface{}, 10)
	if _, ok := sub.in[name]; ok {
		w = sub.in[name]
	} else {
		sub.in[name] = w
	}

	cancel := func() {
	}

	// Publish our data to our output channels
	go func() {
		for {
			data := <-w
			sub.mutex.Lock()
			for _, out := range sub.out[name] {
				out <- data
			}
			sub.mutex.Unlock()
		}
	}()

	return w, cancel
}

// Subscribe provides an interface to listen for a message, and a function
// to cancel the subscription to the messaging bus.
func (sub *Subscription) Subscribe(name string) (chan interface{}, func()) {
	r := make(chan interface{}, 10)
	if _, ok := sub.out[name]; !ok {
		sub.out[name] = make([]chan interface{}, 0)
	}

	sub.out[name] = append(sub.out[name], r)

	cancel := func() {
		sub.mutex.Lock()
		defer sub.mutex.Unlock()

		if _, ok := sub.out[name]; ok {
			for i, ch := range sub.out[name] {
				if ch == r {
					sub.out[name][i] = sub.out[name][len(sub.out[name])-1]
					sub.out[name][len(sub.out[name])] = nil
					sub.out[name] = sub.out[name][:len(sub.out[name])-1]
					break
				}
			}
		}
	}

	return r, cancel
}
