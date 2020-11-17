package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

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

// Emit creates and outputs a CommandOutput with error formatting for the given
// error object.
func (e *ErrInvalidArgument) Emit(s *discordgo.Session, channel string) {
	var cmd = e.cmd

	// Do not run if there is no command object to generate an error for. This is
	// an error.
	if e.cmd == nil {
		if e.sub != nil {
			cmd = e.sub
		} else {
			return
		}
	}

	out := newCommandOutput(cmd, "Command Error")
	out.Status = ifaces.CommandFailure
	out.AddLine(e.Error())
	out.Construct()

	embed, _, _ := GenerateOutputEmbed(out, out.ThisPage())
	s.ChannelMessageSendEmbed(channel, embed)
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

// Emit creates and outputs a CommandOutput with error formatting for the given
// error object.
func (e *ErrInvalidTimezone) Emit(s *discordgo.Session, channel string) {
	var cmd = e.cmd

	// Do not run if there is no command object to generate an error for. This is
	// an error.
	if e.cmd == nil {
		if e.sub != nil {
			cmd = e.sub
		} else {
			return
		}
	}

	out := newCommandOutput(cmd, "Command Error")
	out.Status = ifaces.CommandFailure
	out.AddLine(e.Error())
	out.Construct()

	embed, _, _ := GenerateOutputEmbed(out, out.ThisPage())
	s.ChannelMessageSendEmbed(channel, embed)
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

// Emit creates and outputs a CommandOutput with error formatting for the given
// error object.
func (e *ErrInvalidCommand) Emit(s *discordgo.Session, channel string) {
	var cmd = e.cmd

	// Do not run if there is no command object to generate an error for. This is
	// an error.
	if e.cmd == nil {
		if e.sub != nil {
			cmd = e.sub
		} else {
			return
		}
	}

	out := newCommandOutput(cmd, "Command Error")
	out.Status = ifaces.CommandFailure
	out.AddLine(e.Error())
	out.Construct()

	embed, _, _ := GenerateOutputEmbed(out, out.ThisPage())
	s.ChannelMessageSendEmbed(channel, embed)
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

// Emit creates and outputs a CommandOutput with error formatting for the given
// error object.
func (e *ErrUnauthorizedUsage) Emit(s *discordgo.Session, channel string) {
	var cmd = e.cmd

	// Do not run if there is no command object to generate an error for. This is
	// an error.
	if e.cmd == nil {
		if e.sub != nil {
			cmd = e.sub
		} else {
			return
		}
	}

	out := newCommandOutput(cmd, "Command Error")
	out.Status = ifaces.CommandFailure
	out.AddLine(e.Error())
	out.Construct()

	embed, _, _ := GenerateOutputEmbed(out, out.ThisPage())
	s.ChannelMessageSendEmbed(channel, embed)
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

// Emit creates and outputs a CommandOutput with error formatting for the given
// error object.
func (e *ErrInvalidAlias) Emit(s *discordgo.Session, channel string) {
	var cmd = e.cmd

	// Do not run if there is no command object to generate an error for. This is
	// an error.
	if e.cmd == nil {
		if e.sub != nil {
			cmd = e.sub
		} else {
			return
		}
	}

	out := newCommandOutput(cmd, "Command Error")
	out.Status = ifaces.CommandFailure
	out.AddLine(e.Error())
	out.Construct()

	embed, _, _ := GenerateOutputEmbed(out, out.ThisPage())
	s.ChannelMessageSendEmbed(channel, embed)
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

// Emit creates and outputs a CommandOutput with error formatting for the given
// error object.
func (e *ErrCommandDisabled) Emit(s *discordgo.Session, channel string) {
	var cmd = e.cmd

	// Do not run if there is no command object to generate an error for. This is
	// an error.
	if e.cmd == nil {
		if e.sub != nil {
			cmd = e.sub
		} else {
			return
		}
	}

	out := newCommandOutput(cmd, "Command Error")
	out.Status = ifaces.CommandFailure
	out.AddLine(e.Error())
	out.Construct()

	embed, _, _ := GenerateOutputEmbed(out, out.ThisPage())
	s.ChannelMessageSendEmbed(channel, embed)
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

// Emit creates and outputs a CommandOutput with error formatting for the given
// error object.
func (e *ErrCommandError) Emit(s *discordgo.Session, channel string) {
	var cmd = e.cmd

	// Do not run if there is no command object to generate an error for. This is
	// an error.
	if e.cmd == nil {
		if e.sub != nil {
			cmd = e.sub
		} else {
			return
		}
	}

	out := newCommandOutput(cmd, "Command Error")
	out.Status = ifaces.CommandFailure
	out.AddLine(e.Error())
	out.Construct()

	embed, _, _ := GenerateOutputEmbed(out, out.ThisPage())
	s.ChannelMessageSendEmbed(channel, embed)
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

// Emit creates and outputs a CommandOutput with error formatting for the given
// error object.
func (e *ErrInvalidSubcommand) Emit(s *discordgo.Session, channel string) {
	var cmd = e.cmd

	// Do not run if there is no command object to generate an error for. This is
	// an error.
	if e.cmd == nil {
		if e.sub != nil {
			cmd = e.sub
		} else {
			return
		}
	}

	out := newCommandOutput(cmd, "Command Error")
	out.Status = ifaces.CommandFailure
	out.AddLine(e.Error())
	out.Construct()

	embed, _, _ := GenerateOutputEmbed(out, out.ThisPage())
	s.ChannelMessageSendEmbed(channel, embed)
}
