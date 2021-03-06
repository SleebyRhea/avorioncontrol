package ifaces

import "avorioncontrol/logger"

// IMessageBus describes an interface to a message bus object
type IMessageBus interface {
	NewSubscription(string) (chan interface{}, func())
	Listen(string) (chan interface{}, func())
	logger.ILogger
}
