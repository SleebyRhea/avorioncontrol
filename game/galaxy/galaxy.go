package galaxy

import "avorioncontrol/game/player"

// Galaxy represents an Avorion galaxy with tracked state
type Galaxy struct {
	pcache *player.Cache
	// acache
}

// New returns a new instance of Galaxy
func New() *Galaxy {
	g := &Galaxy{
		pcache: &player.Cache{},
	}

	return g
}

// Players returns the player cache object
func (g *Galaxy) Players() *player.Cache {
	return g.pcache
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
