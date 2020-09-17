package commands

import (
	"avorioncontrol/ifaces"
	"log"
	"time"
)

func init() {
	time.LoadLocation("America/New_York")
}

// HasNumArgs - Determine if a set of command arguments is between min and max
//  @a BotArgs    Argument set to process
//  @min int      Minimum number of positional arguments
//  @max int      Maximum number of positional arguments
//
//  You can use -1 in place of either min or max (or both) to disable the check
//  for that range.
func HasNumArgs(a BotArgs, min, max int) bool {
	if len(a) == 0 || len(a[0]) == 0 {
		log.Fatal("Empty argument list passed to commands.HasNumArgs")
		return false
	}

	if min == -1 {
		min = 0
	}

	if max == -1 {
		max = len(a[1:]) + 1
	}

	if len(a[1:]) > max || len(a[1:]) < min {
		return false
	}

	return true
}

// reverseSlice reverse an arbtrary slice
func reverseJumps(j []*ifaces.JumpInfo) []*ifaces.JumpInfo {
	var jumps []*ifaces.JumpInfo

	var l = len(j)
	var i = l - 1

	if l == 0 {
		return jumps
	}

	for {
		if 0 > i {
			break
		}
		jumps = append(jumps, j[i])
		i--
	}

	return jumps
}

func newArgument(a, b string) CommandArgument {
	return CommandArgument{a, b}
}

// InitializeCommandRegistry Register the commands for commands/everyonecmds
//	@r *CommandRegistrar    The command resistrar that we are initializing
func InitializeCommandRegistry(r *CommandRegistrar) {
	arg := newArgument

	// Global Commands
	r.Register("help",
		"Output help text for a command",
		"help <command> (subcommand) ...",
		[]CommandArgument{
			arg("command",
				"The base command that you would like help with"),
			arg("subcommand",
				"The subcommand that you would like help with")},
		helpCmd)

	r.Register("list",
		"Output a list of all available commands",
		"list",
		make([]CommandArgument, 0),
		listCmd)

	// Debug Commands
	r.Register("ping",
		"Get a \"Pong!\" respons",
		"ping",
		make([]CommandArgument, 0),
		pingCmd)

	r.Register("pong",
		"Get a \"Ping~\" respons",
		"pong",
		make([]CommandArgument, 0),
		pongCmd)

	// Admin Commands
	r.Register("loglevel",
		"Set the log level for a given command(s) or object(s)",
		"loglevel <number> <object> <object> <object> ...",
		[]CommandArgument{
			arg("number",
				"The level of logging being set. Must be a number between 0 and 3"),
			arg("object",
				"Object or command that is being configured. Ex: ping, or guild")},
		loglevelCmnd)

	r.Register("setprefix",
		"Set the command prefix",
		"setprefix <prefix>",
		[]CommandArgument{
			arg("prefix",
				"The prefix that is to be applied.")},
		setprefixCmd)

	r.Register("setalias",
		"Set a command alias",
		"setalias <command> <alias>",
		[]CommandArgument{
			arg("command", "Name of the command the new alias will apply to"),
			arg("alias", "Name of the alias that is being created")},
		setaliasCmd)

	// ifaces.Server (Avorion)
	r.Register("rcon",
		"Run a command in Avorion and return its result",
		"rcon <command> ...",
		[]CommandArgument{
			arg("command", "Name of the command to run"),
			arg("...", "The commands arguments")},
		rconCmnd)

	r.Register("getjumps",
		"Get the last n jumps for a player or alliance",
		"getjumps <type> <max> <name>",
		[]CommandArgument{
			arg("type", "player or alliance type"),
			arg("max", "Max jumps to list"),
			arg("name", "Name of the object to check")},
		getJumpsCmnd)

	r.Register("getcoordhistory",
		"Get all of the logged jumps made to a sector (IN DEV)",
		"getcoordhistory <x:y> <x:y> ...",
		[]CommandArgument{
			arg("x", "x coordinate for a Sector"),
			arg("y", "y coordinate for a sector")},
		getCoordHistoryCmnd)

	r.Register("getplayers",
		"List the tracked players",
		"getplayers",
		make([]CommandArgument, 0),
		getPlayersCmnd)

	r.Register("setchatchannel",
		"Sets the channel to output server chat into",
		"setchatchannel channelid",
		[]CommandArgument{
			arg("channelid", "UID of the channel to send server chat messages to")},
		setChatChannelCmnd)

	r.Register("settimezone",
		"Sets the channel to output server chat into",
		"settimezone timezone",
		[]CommandArgument{
			arg("timezone", "Timezone to set (reference: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)")},
		setTimezoneCmnd)

	r.Register("server",
		"Control the state of the Avorion server",
		"server <start|stop|restart>",
		make([]CommandArgument, 0),
		proxySubCmnd)
	r.Register("stop",
		"Stop the Avorion server (if its up)",
		"stop",
		make([]CommandArgument, 0),
		stopServerCmnd, "server")
	r.Register("start",
		"Start the Avorion server (if its down)",
		"start",
		make([]CommandArgument, 0),
		startServerCmnd, "server")
	r.Register("restart",
		"Restart the Avorion server",
		"restart",
		make([]CommandArgument, 0),
		restartServerCmnd, "server")
}
