package gameserver

import (
	"AvorionControl/logger"
	"net"
)

// Server - A Server describes an interface to a full Server
type Server interface {
	Commandable
	Versioned
	logger.Logger
	Playable
	LoginMessager
	PasswordLockable
	Seeded
	Websocketer

	IsUp() bool
	Stop() error
	Start() error
	Restart() error
}

// Data is a datastructure that represents the current state of a Server
type Data struct {
	WorldName   string
	Online      bool
	Seed        string
	MOTD        string
	Password    string
	Players     []*PlayerData
	PlayerCount int
	Loglevel    int
	Version     string
}

// OutputSender sends output from a Server to a channel
type OutputSender interface {
	SetSendChannel(chan []byte)
}

// Playable - Define an object that can track the players that have joined
type Playable interface {
	Player(string) Player
	Players() []Player
	NewPlayer(string, string) Player
	RemovePlayer(string) bool
}

// Player - Define a player than can join a server and has various details
// regarding its connection tracked
type Player interface {
	Name() string
	GetData() error

	Kick(string)
	Ban(string)

	SetIP(string)
	IP() net.IP

	Online() bool
	SetOnline()
	SetOffline()
}

// PlayerData represents a player, and includes the name and their IP address.
// the IP address is in string form, as this is intended to be marshalled into a
// json file.
type PlayerData struct {
	Name string
	IP   string
}

// Websocketer is an object that is able to output to the Guis websocket ub
type Websocketer interface {
	WSOutput() chan []byte
}

// Versioned is an interface to objects with verisons
type Versioned interface {
	SetVersion(string)
	Version() string
}

// PasswordLockable is an interface to an object that can have its password set
type PasswordLockable interface {
	Password() string
	SetPassword(string)
}

// LoginMessager is an interface to an object that can have an MOTD set
type LoginMessager interface {
	MOTD() string
	SetMOTD(string)
}

// Seeded is an interface to an object that has a seed
type Seeded interface {
	Seed() string
	SetSeed(string)
}

// Commandable - A Commandable object must implement the function RunCommand
type Commandable interface {
	RunCommand(string) (string, error)
}

// SendCommand - Send a command to a Commandable() object
func SendCommand(s string, cs Commandable) {
	cs.RunCommand(s)
}

// GamePlayerData returns a playerdata object for a Server
func GamePlayerData(gs Server) []*PlayerData {
	d := make([]*PlayerData, 0)
	for _, p := range gs.Players() {
		d = append(d, &PlayerData{
			Name: p.Name(),
			IP:   p.IP().String(),
		})
	}
	return d
}

// GameStatus constructs a new Data struct from the given Server
func GameStatus(gs Server) *Data {
	return &Data{
		WorldName:   "Avorion",
		Online:      gs.IsUp(),
		Seed:        gs.Seed(),
		Password:    gs.Password(),
		Players:     GamePlayerData(gs),
		PlayerCount: len(gs.Players()),
		Loglevel:    gs.Loglevel(),
		Version:     gs.Version(),
	}
}
