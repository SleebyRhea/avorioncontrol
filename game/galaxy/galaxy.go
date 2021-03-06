package galaxy

import (
	"avorioncontrol/game/player"
	"avorioncontrol/ifaces"
)

// Galaxy represents an Avorion galaxy with tracked state
type Galaxy struct {
	name   string
	path   string
	pcache *player.Cache
}

// New returns a new instance of Galaxy
func New() *Galaxy {
	g := &Galaxy{
		pcache: &player.Cache{},
	}

	return g
}

// Players returns the player cache object
func (g *Galaxy) Players() ifaces.IPlayerCache {
	return g.pcache
}

// Alliances returns the alliance cache object
func (g *Galaxy) Alliances() {
}

// Name returns the name of the Galaxy
func (g *Galaxy) Name() string {
	return g.name
}

// Path return the filepath to the Galaxy
func (g *Galaxy) Path() string {
	return g.path
}

// Players returns the player cache object
func (g *Galaxy) Players() ifaces.IPlayerCache {
	return g.pcache
}

// Alliances returns the alliance cache object
func (g *Galaxy) Alliances() {
}

// Do this on startup
// if _, err := os.Stat(galaxydir); os.IsNotExist(err) {
// 	err := os.Mkdir(galaxydir, 0700)
// 	if err != nil {
// 		logger.LogError(s, "os.Mkdir: "+err.Error())
// 	}
// }

// if err := s.config.BuildModConfig(); err != nil {
// 	return err
// }
