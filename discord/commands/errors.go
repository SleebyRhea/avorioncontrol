package commands

// ErrInvalidArgument describes an invalid attempt to use a command
// due to incorrect arguments
type ErrInvalidArgument struct {
	cmd     *CommandRegistrant
	sub     *CommandRegistrant
	message string
}

// Command returns the command object that encountered an error
func (e *ErrInvalidArgument) Command() *CommandRegistrant {
	return e.cmd
}

// Subcommand returns the subcommand object that encountered an error
func (e *ErrInvalidArgument) Subcommand() *CommandRegistrant {
	return e.sub
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
	cmd *CommandRegistrant
	sub *CommandRegistrant
	tz  string
}

// Command returns the command object that encountered an error
func (e *ErrInvalidTimezone) Command() *CommandRegistrant {
	return e.cmd
}

// Subcommand returns the subcommand object that encountered an error
func (e *ErrInvalidTimezone) Subcommand() *CommandRegistrant {
	return e.sub
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
	name string
	cmd  *CommandRegistrant
	sub  *CommandRegistrant
}

// Command returns the command object that encountered an error
func (e *ErrInvalidCommand) Command() *CommandRegistrant {
	return e.cmd
}

// Subcommand returns the subcommand object that encountered an error
func (e *ErrInvalidCommand) Subcommand() *CommandRegistrant {
	return e.sub
}

func (e *ErrInvalidCommand) Error() string {
	if e.name == "" {
		return "Invalid command supplied"
	}

	return sprintf("Command `%s` is invalid", e.name)
}

// ErrUnauthorizedUsage describes an attempt to run a command by someone
// unauthorized to do so
type ErrUnauthorizedUsage struct {
	cmd *CommandRegistrant
	sub *CommandRegistrant
}

// Command returns the command object that encountered an error
func (e *ErrUnauthorizedUsage) Command() *CommandRegistrant {
	return e.cmd
}

// Subcommand returns the subcommand object that encountered an error
func (e *ErrUnauthorizedUsage) Subcommand() *CommandRegistrant {
	return e.sub
}

func (e *ErrUnauthorizedUsage) Error() string {
	return "You do not have permission to use that command"
}

// ErrInvalidAlias describes an attempt to use an alias that doesn't
// exist
type ErrInvalidAlias struct {
	cmd   *CommandRegistrant
	sub   *CommandRegistrant
	alias string
}

// Command returns the command object that encountered an error
func (e *ErrInvalidAlias) Command() *CommandRegistrant {
	return e.cmd
}

// Subcommand returns the subcommand object that encountered an error
func (e *ErrInvalidAlias) Subcommand() *CommandRegistrant {
	return e.sub
}

func (e *ErrInvalidAlias) Error() string {
	if e.alias == "" {
		return "Invalid command alias"
	}

	return sprintf("Alias `%s` is invalid", e.alias)
}

// ErrCommandDisabled describes an attempt to use a command that has
// been disabled
type ErrCommandDisabled struct {
	cmd *CommandRegistrant
	sub *CommandRegistrant
}

// Command returns the command object that encountered an error
func (e *ErrCommandDisabled) Command() *CommandRegistrant {
	return e.cmd
}

// Subcommand returns the subcommand object that encountered an error
func (e *ErrCommandDisabled) Subcommand() *CommandRegistrant {
	return e.sub
}

func (e *ErrCommandDisabled) Error() string {
	if e.cmd == nil {
		return "That command has been disabled"
	}

	return sprintf("The command `%s` has been disabled", e.cmd.Name())
}

// ErrCommandError describes a generic non-fatal error that occurred
// during command processing.
type ErrCommandError struct {
	cmd     *CommandRegistrant
	sub     *CommandRegistrant
	message string
}

// Command returns the command object that encountered an error
func (e *ErrCommandError) Command() *CommandRegistrant {
	return e.cmd
}

// Subcommand returns the subcommand object that encountered an error
func (e *ErrCommandError) Subcommand() *CommandRegistrant {
	return e.sub
}

func (e *ErrCommandError) Error() string {
	if e.message == "" {
		return "Command encountered an unspecified failure"
	}

	return e.message
}

// ErrInvalidSubcommand describes an error in which a provided subcommand
// does not exist
type ErrInvalidSubcommand struct {
	cmd     *CommandRegistrant
	sub     *CommandRegistrant
	subname string
}

// Command returns the command object that encountered an error
func (e *ErrInvalidSubcommand) Command() *CommandRegistrant {
	return e.cmd
}

// Subcommand returns the subcommand object that encountered an error
func (e *ErrInvalidSubcommand) Subcommand() *CommandRegistrant {
	return e.sub
}

func (e *ErrInvalidSubcommand) Error() string {
	if e.subname == "" || e.cmd == nil {
		return "Subcommand is not registered"
	}

	return sprintf("%s does not have the subcommand `%s` registered",
		e.cmd.Name(), e.subname)
}
