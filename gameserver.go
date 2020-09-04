package main

import (
	"net"
	"regexp"
)

var gameServers []GameServer
var illegalNamesRe []*regexp.Regexp

// GameServer - A GameServer describes an interface to a full GameServer
type GameServer interface {
	Commandable
	Versioned
	Loggable
	Playable
	Server
	LoginMessager
	PasswordLockable
	Seeded
	Websocketer
}

// GameData is a datastructure that represents the current state of a GameServer
type GameData struct {
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

// OutputSender sends output from a GameServer to a channel
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
	SetIP(string)
	Name() string
	Kick(string)
	Ban(string)
	IP() net.IP
}

// PlayerData represents a player, and includes the name and their IP address.
// the IP address is in string form, as this is intended to be marshalled into a
// json file.
type PlayerData struct {
	Name string
	IP   string
}

// Server -
type Server interface {
	IsUp() bool
	Stop() error
	Start() error
	Restart() error
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

// GamePlayerData returns a playerdata object for a GameServer
func GamePlayerData(gs GameServer) []*PlayerData {
	d := make([]*PlayerData, 0)
	for _, p := range gs.Players() {
		d = append(d, &PlayerData{
			Name: p.Name(),
			IP:   p.IP().String(),
		})
	}
	return d
}

// GameStatus constructs a new GameData struct from the given GameServer
func GameStatus(gs GameServer) *GameData {
	return &GameData{
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

func init() {
	// Prepare our application data
	gameServers = make([]GameServer, 0)
}
