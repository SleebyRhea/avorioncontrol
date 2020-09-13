package ifaces

import (
	"AvorionControl/logger"
)

// IAlliance describes a an IServer player alliance
type IAlliance interface {
	ITrackedAlliance
}

// ITrackedAlliance defines an interface to an an alliance that has tracking
type ITrackedAlliance interface {
	logger.ILogger
	Name() string
	Index() string
	Update() error
	Message(string)
	UpdateCoords(ShipCoordData)
}
