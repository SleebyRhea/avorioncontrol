package commands

// ErrInvalidArgument describes an invalid attempt to use a command
// due to incorrect arguments
type ErrInvalidArgument struct {
	message string
}

func (e *ErrInvalidArgument) Error() string {
	if e.message == "" {
		return "Invalid argument supplied"
	}

	return e.message
}

// ErrInvalidTimezone describes an attempt to use an invalid timezone
// that was configured
type ErrInvalidTimezone struct {
	tz string
}

func (e *ErrInvalidTimezone) Error() string {
	if e.tz == "" {
		return "Invalid timezone configured"
	}

	return sprintf("Configured timezone `%s` is invalid", e.tz)
}

// ErrInvalidCommand describes an attempt to run a command that doesn't
// exist
type ErrInvalidCommand struct {
	command string
}

func (e *ErrInvalidCommand) Error() string {
	if e.command == "" {
		return "Invalid command supplied"
	}

	return sprintf("Command `%s` is invalid", e.command)
}

// ErrUnauthorizedUsage describes an attempt to run a command by someone
// unauthorized to do so
type ErrUnauthorizedUsage struct {
	command string
}

func (e *ErrUnauthorizedUsage) Error() string {
	return "You do not have permission to use that command"
}

// ErrInvalidAlias describes an attempt to use an alias that doesn't
// exist
type ErrInvalidAlias struct {
	alias string
}

func (e *ErrInvalidAlias) Error() string {
	if e.alias == "" {
		return "Invalid command alias"
	}

	return sprintf("Alias `%s` is invalid")
}

// ErrCommandDisabled describes an attempt to use a command that has
// been disabled
type ErrCommandDisabled struct {
	command string
}

func (e *ErrCommandDisabled) Error() string {
	if e.command == "" {
		return "That command has been disabled"
	}

	return sprintf("The command `%s` has been disabled")
}

// ErrCommandError describes a generic non-fatal error that occurred
// during command processing.
type ErrCommandError struct {
	message string
}

func (e *ErrCommandError) Error() string {
	if e.message == "" {
		return "Command encountered an unspecified failure"
	}

	return e.message
}
