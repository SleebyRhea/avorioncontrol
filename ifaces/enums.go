package ifaces

// enums describing the status of a server
const (
	ServerOffline    = 0
	ServerOnline     = 1
	ServerStarting   = 2
	ServerStopping   = 3
	ServerRestarting = 4
	ServerCrashed    = 255

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

func init() {
	difficultyMap = make(map[int]string, 0)
	difficultyMap[difficultyBeginner] = difficultyBeginnerString
	difficultyMap[difficultyEasy] = difficultyEasyString
	difficultyMap[difficultyNormal] = difficultyNormalString
	difficultyMap[difficultyVeteran] = difficultyVeteranString
	difficultyMap[difficultyExpert] = difficultyExpertString
	difficultyMap[difficultyHardcore] = difficultyHardcoreString
	difficultyMap[difficultyInsane] = difficultyInsaneString
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
