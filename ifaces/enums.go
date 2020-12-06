package ifaces

// enums describing the status of a server
const (
	ServerOffline           = 0
	ServerOnline            = 1
	ServerStarting          = 2
	ServerStopping          = 3
	ServerRestarting        = 4
	ServerCrashedOffline    = 255
	ServerCrashedRecovered  = 256
	ServerCrashedStarting   = 257
	ServerCrashedStopping   = 258
	ServerCrashedRestarting = 259

	ServerOfflineColor           = 16711680
	ServerOnlineColor            = 6749952
	ServerStartingColor          = 3381759
	ServerStoppingColor          = 16776960
	ServerRestartingColor        = 3381759
	ServerCrashedOfflineColor    = 15158332
	ServerCrashedRecoveredColor  = 15158332
	ServerCrashedStartingColor   = 15158332
	ServerCrashedStoppingColor   = 15158332
	ServerCrashedRestartingColor = 15158332

	ServerOfflineString           = "Offline"
	ServerOnlineString            = "Online"
	ServerStartingString          = "Initializing"
	ServerStoppingString          = "Stopping"
	ServerRestartingString        = "Re-Initializing"
	ServerCrashedOfflineString    = "Crashed (Dead)"
	ServerCrashedRecoveredString  = "Crashed (Recovered)"
	ServerCrashedStartingString   = "Crashed (Recovering)"
	ServerCrashedStoppingString   = "Crashed (Attempting Graceful Exit)"
	ServerCrashedRestartingString = "Crashed (Restarting)"

	CommandSuccess = 0
	CommandFailure = 1
	CommandWarning = 2

	difficultyBeginner = -3
	difficultyEasy     = -2
	difficultyNormal   = -1
	difficultyVeteran  = 0
	difficultyExpert   = 1
	difficultyHardcore = 2
	difficultyInsane   = 3

	difficultyBeginnerString = "Beginner"
	difficultyEasyString     = "Easy"
	difficultyNormalString   = "Normal"
	difficultyVeteranString  = "Veteran"
	difficultyExpertString   = "Expert"
	difficultyHardcoreString = "Hardcore"
	difficultyInsaneString   = "Insane"
)

var difficultyMap map[int]string
var stateMapString map[int]string
var stateMapColor map[int]int

func init() {
	difficultyMap = make(map[int]string, 0)
	difficultyMap[difficultyBeginner] = difficultyBeginnerString
	difficultyMap[difficultyEasy] = difficultyEasyString
	difficultyMap[difficultyNormal] = difficultyNormalString
	difficultyMap[difficultyVeteran] = difficultyVeteranString
	difficultyMap[difficultyExpert] = difficultyExpertString
	difficultyMap[difficultyHardcore] = difficultyHardcoreString
	difficultyMap[difficultyInsane] = difficultyInsaneString

	stateMapString = make(map[int]string, 0)
	stateMapString[ServerOffline] = ServerOfflineString
	stateMapString[ServerOnline] = ServerOnlineString
	stateMapString[ServerStarting] = ServerStartingString
	stateMapString[ServerStopping] = ServerStoppingString
	stateMapString[ServerRestarting] = ServerRestartingString
	stateMapString[ServerCrashedOffline] = ServerCrashedOfflineString
	stateMapString[ServerCrashedRecovered] = ServerCrashedRecoveredString
	stateMapString[ServerCrashedStarting] = ServerCrashedStartingString
	stateMapString[ServerCrashedStopping] = ServerCrashedStoppingString
	stateMapString[ServerCrashedRestarting] = ServerCrashedRestartingString

	stateMapColor = make(map[int]int, 0)
	stateMapColor[ServerOffline] = ServerOfflineColor
	stateMapColor[ServerOnline] = ServerOnlineColor
	stateMapColor[ServerStarting] = ServerStartingColor
	stateMapColor[ServerStopping] = ServerStoppingColor
	stateMapColor[ServerRestarting] = ServerRestartingColor
	stateMapColor[ServerCrashedOffline] = ServerCrashedOfflineColor
	stateMapColor[ServerCrashedRecovered] = ServerCrashedRecoveredColor
	stateMapColor[ServerCrashedStarting] = ServerCrashedStartingColor
	stateMapColor[ServerCrashedStopping] = ServerCrashedStoppingColor
	stateMapColor[ServerCrashedRestarting] = ServerCrashedRestartingColor
}

// Difficulty returns the difficulty string for a given int
func Difficulty(d int) string {
	if d > difficultyInsane {
		d = difficultyInsane
	}

	if d < difficultyBeginner {
		d = difficultyBeginner
	}

	return difficultyMap[d]
}

// State returns the matching state string and color for a given server state
func State(s int) (string, int) {
	if s > ServerRestarting && s < ServerCrashedOffline {
		s = ServerRestarting
		return stateMapString[s], stateMapColor[s]
	}

	if s > ServerCrashedRestarting {
		s = ServerCrashedRestarting
		return stateMapString[s], stateMapColor[s]
	}

	if s < ServerOffline {
		s = ServerOffline
		return stateMapString[s], stateMapColor[s]
	}

	return stateMapString[s], stateMapColor[s]
}
