package commands

import (
	"AvorionControl/logger"
	"fmt"
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
	return c.name + "Cmd:" + c.registrar.GuildID
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
//  TODO: Convert this into an embed
//  Current Format:
//    **Command:**
//    `cmdname` - Command description
//    ```
//    Command usage
//    ```
//    **Arguments:**
//    ```
//    1: Arg1 - Description
//    2: Arg2 - Description
//    ```
//    **Subcommands:**
//    ```
//    1: sub1 - Description
//    2: sub1 - Description
//    ```
func (c *CommandRegistrant) Help() (out string, _ error) {
	out = fmt.Sprintf("**`%s`** - %s\n```\n%s\n```\n", c.name,
		c.description, c.usage)

	if len(c.args) > 0 {
		out = fmt.Sprintf("%s**Arguments:**\n```\n", out)
		for _, a := range c.args {
			out = fmt.Sprintf("%s%s - %s\n", out, a[0], a[1])
		}
		out = fmt.Sprintf("%s```\n", out)
	}

	if len(c.cmdlets) > 0 {
		out = fmt.Sprintf("%s**Subcommands:**\n```\n", out)
		for _, sc := range c.cmdlets {
			out = fmt.Sprintf("%s%s - %s\n", out, sc.Name(),
				sc.description)
		}
		out = fmt.Sprintf("%s```\n", out)
	}

	return out, nil
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
