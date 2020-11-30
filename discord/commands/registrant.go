package commands

import (
	"avorioncontrol/logger"
)

// CommandRegistrant - Command definition
type CommandRegistrant struct {
	name        string
	description string

	exec         BotCommand
	cmdlets      []*CommandRegistrant
	args         []CommandArgument
	usage        string
	hasauthlevel bool
	registrar    *CommandRegistrar

	// Implements logger.Loggable
	loglevel int
}

// SetLoglevel - Set the current loglevel
func (c *CommandRegistrant) SetLoglevel(l int) {
	c.loglevel = l
	logger.LogInfo(c, sprintf("Setting loglevel to %d", l))
}

// UUID - Return a commands name and guild information
func (c *CommandRegistrant) UUID() string {
	return c.name + "Cmnd:" + c.registrar.GuildID
}

// Loglevel - Return the current loglevel for the CommandRegistrant
func (c *CommandRegistrant) Loglevel() int {
	var l int
	if c.loglevel < 0 {
		l = c.registrar.Loglevel()
	} else {
		l = c.loglevel
	}

	return l
}

// Name - Return a commands name.
func (c *CommandRegistrant) Name() string {
	return c.name
}

// Help - Return a commands help text with following markdown formatting
// preapplied
func (c *CommandRegistrant) Help() *CommandOutput {
	out := newCommandOutput(c, "Command Help")
	out.Description = c.description
	out.Header = "Usage"
	out.AddLine(sprintf("> %s", c.usage))

	if len(c.args) > 0 {
		out.AddLine("**Arguments**")
		for _, a := range c.args {
			out.AddLine(sprintf("> %s - %s", a[0], a[1]))
		}
	}

	if len(c.cmdlets) > 0 {
		out.AddLine("**Subcommands**")
		for _, sc := range c.cmdlets {
			out.AddLine(sprintf("> %s - %s", sc.Name(), sc.description))
		}
	}

	out.Construct()
	return out
}

// Subcommands - Return all the count of the subcommands added to a command, and
// their slice.
func (c *CommandRegistrant) Subcommands() (int, []*CommandRegistrant) {
	if c.cmdlets != nil {
		return len(c.cmdlets), c.cmdlets
	}
	return 0, nil
}

// Registrar - Return a commands master CommandRegistrar
func (c *CommandRegistrant) Registrar() *CommandRegistrar {
	return c.registrar
}
