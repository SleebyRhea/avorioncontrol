package server

// ErrInvalidFactionID describes an error in which an invalid Faction ID was given
type ErrInvalidFactionID struct {
}

func (e *ErrInvalidFactionID) Error() string {
	return "invalid Faction ID provided"
}

// ErrFailedToStart describes an error in which Avorion failed to start
type ErrFailedToStart struct {
}

func (e *ErrFailedToStart) Error() string {
	return "the Avorion binary failed to start"
}

// ErrFailedInit describes an error in which Avorion failed finish initializing
type ErrFailedInit struct {
}

func (e *ErrFailedInit) Error() string {
	return "the Avorion process failed to successfully initialize"
}

// ErrRconFailedToStart describes an error in which RCON failed to run
type ErrRconFailedToStart struct {
}

func (e *ErrRconFailedToStart) Error() string {
	return "failed to run rcon binary"
}

// ErrCommandFailedToRun describes an error in which a command failed to run
type ErrCommandFailedToRun struct {
}

func (e *ErrCommandFailedToRun) Error() string {
	return "failed to run command"
}

// ErrCommandTimedOut describes an error in which a command timed out
type ErrCommandTimedOut struct {
}

func (e *ErrCommandTimedOut) Error() string {
	return "command timed out before it could complete"
}

// ErrCommandInvalid describes an error in which a command doesnt exist
type ErrCommandInvalid struct {
	Cmd string
}

func (e *ErrCommandInvalid) Error() string {
	return "invalid rcon command: " + e.Cmd
}

// ErrServerOffline describes an error in which a method cannot be run
// due to the server being offline
type ErrServerOffline struct {
}

func (e *ErrServerOffline) Error() string {
	return "server is currently offline"
}

// ErrServerOnline describes an error in which a method cannot be run
// due to the server being online
type ErrServerOnline struct {
}

func (e *ErrServerOnline) Error() string {
	return "server is already online"
}
