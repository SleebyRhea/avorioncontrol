package commands

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

	// // Debug Commands
	r.Register("ping",
		"Get a \"Pong!\" response",
		"ping",
		make([]CommandArgument, 0),
		pingCmd)

	r.Register("pong",
		"Get a \"Ping~\" respons",
		"pong",
		make([]CommandArgument, 0),
		pongCmd)

	// // Admin Commands
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

	r.Register("status",
		"Get the current server status",
		"status",
		make([]CommandArgument, 0),
		statusCmnd)

	r.Register("getjumps",
		"Get the last n jumps for a player or alliance",
		"getjumps <number> <name>",
		[]CommandArgument{
			arg("number", "Number of jumps to list (25 max)"),
			arg("name", "Player or Alliance name")},
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

	r.Register("reload",
		"Reloads the active configuration from our config file",
		"reload",
		make([]CommandArgument, 0),
		reloadConfigCmnd)

	r.Register("setchatchannel",
		"Sets the channel to output server chat into",
		"setchatchannel channelid",
		[]CommandArgument{
			arg("channelid", "UID of the channel to send server chat messages to")},
		setChatChannelCmnd)

	r.Register("setstatuschannel",
		"Sets the channel in which the server will update it's status embed",
		"setstatuschannel channelid",
		[]CommandArgument{
			arg("channelid", "UID of the channel to send server chat messages to")},
		setStatusChannelCmnd)

	r.Register("settimezone",
		"Sets the channel to output server chat into",
		"settimezone timezone",
		[]CommandArgument{
			arg("timezone", "Timezone to set (reference: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)")},
		setTimezoneCmnd)

	r.Register("server",
		"Control the state of the Avorion server",
		"server <subcommand>",
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

	r.Register("admin",
		"Configure admin level privileges",
		"admin <subcommand>",
		make([]CommandArgument, 0),
		proxySubCmnd)
	r.Register("roles",
		"List admin level roles",
		"roles",
		make([]CommandArgument, 0),
		showAdminRolesSubCmnd, "admin")
	r.Register("commands",
		"List commands that need admin",
		"commands",
		make([]CommandArgument, 0),
		showAdminCmndsSubCmnd, "admin")
	r.Register("addrole",
		"Add a set level of authorization to a role",
		"addrole <role> <level>",
		[]CommandArgument{
			arg("role", "role to add auth privileges to"),
			arg("level", "number representing authorization level (max 10)")},
		addAdminRoleSubCmnd, "admin")
	r.Register("delrole",
		"Remove authorization from a role",
		"delrole <role>",
		[]CommandArgument{
			arg("role", "role to add auth privileges to")},
		removeAdminRoleSubCmnd, "admin")
	r.Register("addcommand",
		"Require authorization for a command",
		"addcommand <command> <level>",
		[]CommandArgument{
			arg("level", "number representing authorization level (max 10)"),
			arg("command", "role to add auth privileges to")},
		addAdminCmndSubCmnd, "admin")
	r.Register("delcommand",
		"Remove the auth requirements from a command",
		"delcommand <command>",
		[]CommandArgument{
			arg("command", "command to have remove requirements from")},
		removeAdminCmndSubCmnd, "admin")
	r.Register("self",
		"Output your authorization level",
		"delcommand <command>",
		make([]CommandArgument, 0),
		showSelfCmndSubCmnd, "admin")

	r.Register("mod",
		"Configure mods installed on the Avorion server",
		"mod <add|remove|list>",
		make([]CommandArgument, 0),
		proxySubCmnd)
	r.Register("add",
		"Add a mod or mods to the server configuration",
		"add <workshopid> <workshopid> ...",
		[]CommandArgument{
			arg("workshopid", "Steam workshop ID of a mod to add")},
		modAddSubCmnd, "mod")
	r.Register("allow",
		"Allow a mod to be installed on the client",
		"add <workshopid> <workshopid> ...",
		[]CommandArgument{
			arg("workshopid", "Steam workshop ID of a mod to add")},
		modAllowSubCmnd, "mod")
	r.Register("disallow",
		"Remove a mod from the allowed client mods list",
		"add <workshopid> <workshopid> ...",
		[]CommandArgument{
			arg("workshopid", "Steam workshop ID of a mod to add")},
		modDisallowSubCmnd, "mod")
	r.Register("remove",
		"Remove a mod or mods from the server configuration",
		"remove <workshopid> <workshopid> ...",
		[]CommandArgument{
			arg("workshopid", "Steam workshop ID of a mod to add")},
		modRemoveSubCmnd, "mod")
	r.Register("list",
		"List the workshop mods that are currently configured to be installed",
		"list",
		make([]CommandArgument, 0),
		listModsSubCmnd, "mod")

	r.Register("modlist",
		"List the workshop mods that are currently configured to be installed",
		"list",
		make([]CommandArgument, 0),
		listModsSubCmnd)
	r.Register("showuserintegrations",
		"Show the users that have integrated Discord and their in-game player",
		"showuserintegrations",
		make([]CommandArgument, 0),
		getIntegratedCmnd)

	r.Register("broadcast",
		"Send all players an email, with an attachment used as the message body",
		"broadcast",
		make([]CommandArgument, 0),
		sendBroadcastCmnd)
}
