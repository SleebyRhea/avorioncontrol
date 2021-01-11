package player

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"fmt"
	"strconv"
	"sync"
)

// Cache describes a cache of player data
type Cache struct {
	loglevel int

	mutex struct {
		dmap *sync.Mutex
		smap *sync.Mutex
		fmap *sync.Mutex
		nmap *sync.Mutex
	}

	plrlist []*Player
	discord map[string]*Player
	steam64 map[string]*Player
	faction map[string]*Player
	namestr map[string]*Player
}

func (c *Cache) lockAll() {
	c.mutex.dmap.Lock()
	c.mutex.smap.Lock()
	c.mutex.fmap.Lock()
	c.mutex.nmap.Lock()
}

func (c *Cache) unlockAll() {
	c.mutex.dmap.Unlock()
	c.mutex.smap.Unlock()
	c.mutex.fmap.Unlock()
	c.mutex.nmap.Unlock()
}

// NewPlayer initializes and returns a new player
func (c *Cache) NewPlayer(sid, fid, name, did string, gs ifaces.IGameServer) (
	ifaces.IPlayer, error) {
	if fid == "" {
		return nil, &ErrEmptyFactionID{}
	}

	if sid == "" {
		return nil, &ErrEmptySteam64ID{}
	}

	if name == "" {
		return nil, &ErrEmptyName{}
	}

	_, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return nil, &ErrBadSteam64ID{}
	}

	player := &Player{
		name:      name,
		factionid: fid,
		steam64id: sid,
		discordid: did,
		isonline:  false}

	c.lockAll()
	defer c.unlockAll()

	c.plrlist = append(c.plrlist, player)
	c.namestr[name] = player
	c.steam64[sid] = player
	c.faction[fid] = player
	if did != "" {
		c.discord[did] = player
	}

	logger.LogDebug(player, "Registered player")
	return player, nil
}

// UUID returns the name of the onject for logging
func (c *Cache) UUID() string {
	return "PlayerCache"
}

// SetLoglevel sets a player Caches level of logging
func (c *Cache) SetLoglevel(level int) {
	c.lockAll()
	defer c.unlockAll()

	logger.LogInfo(c, fmt.Sprintf("Setting logging level to %d", level))
	c.loglevel = level
	for _, p := range c.plrlist {
		p.SetLoglevel(level)
	}
}

// Loglevel returns the integer loglevel for Cache
func (c *Cache) Loglevel() int {
	return c.loglevel
}

// FromFactionID returns a player given a faction ID
func (c *Cache) FromFactionID(ref string) ifaces.IPlayer {
	c.mutex.fmap.Lock()
	defer c.mutex.fmap.Unlock()

	if p, ok := c.faction[ref]; ok {
		return p
	}

	return nil
}

// FromDiscordID returns a player given a discord ID
func (c *Cache) FromDiscordID(ref string) ifaces.IPlayer {
	c.mutex.fmap.Lock()
	defer c.mutex.dmap.Unlock()

	if p, ok := c.discord[ref]; ok {
		return p
	}

	return nil
}

// FromSteamID returns a player given a steam ID
func (c *Cache) FromSteamID(ref string) ifaces.IPlayer {
	c.mutex.fmap.Lock()
	defer c.mutex.smap.Unlock()

	if p, ok := c.steam64[ref]; ok {
		return p
	}

	return nil
}

// FromName returns a player given a name string
func (c *Cache) FromName(ref string) ifaces.IPlayer {
	c.mutex.fmap.Lock()
	defer c.mutex.nmap.Unlock()

	if p, ok := c.namestr[ref]; ok {
		return p
	}

	return nil
}

// SetPlayerName sets a given player to have the given name
func (c *Cache) SetPlayerName(name string, player *Player) {
	c.mutex.nmap.Lock()
	defer c.mutex.nmap.Unlock()
	c.namestr[name] = nil
	c.namestr[name] = player
	player.name = name
}

// SetPlayerDiscord sets a given player to use a given Discord ID
func (c *Cache) SetPlayerDiscord(did string, player *Player) error {
	c.mutex.dmap.Lock()
	defer c.mutex.dmap.Unlock()

	if _, ok := c.discord[did]; ok {
		return &ErrDiscordMapped{}
	}

	c.discord[did] = player
	logger.LogInfo(player, fmt.Sprintf(`Assigned Discord ID %s to %s`, did,
		player.steam64id))
	return nil
}
