package discord

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// DataCache represents a cache of member nicknames on registered Guilds
type DataCache struct {
	loglevel int

	guildcache map[string]bool
	namecache  map[string]map[string]string
	colorcache map[string]map[string]CachedColor

	mutex struct {
		guilds *sync.Mutex
		color  *sync.Mutex
		name   *sync.Mutex
	}
}

// CachedColor describes a role color that has been cached with its Number,
// hex code, and shortcode
type CachedColor struct {
	Integer int
	String  string
	Short   string
}

/************************/
/* IFace logger.ILogger */
/************************/

// SetLoglevel sets the current loglevel for the object
func (d *DataCache) SetLoglevel(l int) {
	d.loglevel = l
}

// Loglevel returns the current loglevel for the object
func (d *DataCache) Loglevel() int {
	return d.loglevel
}

// UUID returns the UUID for the Logger
func (d *DataCache) UUID() string {
	return "DataCache"
}

/********************/
/* Struct DataCache */
/********************/

// UpdateCache updates the internal name cache for nicknames and linked player
// names. Linked player names takes precedence.
func (d *DataCache) UpdateCache(s *discordgo.Session, gs ifaces.IGameServer) {
	logger.LogInfo(d, "Updating Discord data cache")

	newnamecache := make(map[string]map[string]string, 0)
	newcolorcache := make(map[string]map[string]CachedColor, 0)

GUILDGET:
	for gid, doupdate := range d.guildcache {
		if !doupdate {
			logger.LogInfo(d, "Skipping: "+gid)
			continue
		}

		logger.LogInfo(d, "Updating: "+gid)

		_, err := s.State.Guild(gid)
		if err != nil {
			logger.LogError(d, "cannot update invalid guild name cache: "+gid)
			return
		}

		newnamecache[gid] = make(map[string]string, 0)
		newcolorcache[gid] = make(map[string]CachedColor, 0)

		if gs != nil {
			logger.LogDebug(d, "Checking IGameServer players for linked discord users")
			for _, p := range gs.Players() {
				if _pdid := p.DiscordUID(); _pdid != "" {
					newnamecache[gid][_pdid] = p.Name()
				}
			}
		}

		var (
			last       string
			guildroles []*discordgo.Role
			members    []*discordgo.Member
		)

		last = ""
		logger.LogInfo(d, "Getting members list for: "+gid)

		guildroles, err = s.GuildRoles(gid)
		if err != nil {
			logger.LogError(d, "discordgo (*Session).GuildRoles: ")
			continue GUILDGET
		}

		guildroles = rolesort(guildroles)

		for {
			logger.LogDebug(d, "Fetching members from: "+last)
			_members, err := s.GuildMembers(gid, last, 1000)
			if err != nil {
				logger.LogError(d, "discordgo (*Session).GuildMembers: "+err.Error())
				continue GUILDGET
			}

			if len(_members) > 0 {
				members = append(members, _members...)
				last = _members[len(_members)-1].User.ID
				continue
			}

			break
		}

	MEMBERGET:
		for _, m := range members {
			logger.LogDebug(d, "Running update for: "+m.User.String())
			name := m.User.Username
			if m.Nick != "" {
				name = m.Nick
			}

			if pname, ok := newnamecache[gid][m.User.ID]; ok {
				name = pname
			}

			userroles := make([]*discordgo.Role, 0)

			for _, rid := range m.Roles {
				_r, err := s.State.Role(gid, rid)
				if err != nil {
					logger.LogError(d, "discordgo (*Session).State.Role(): "+err.Error())
					continue MEMBERGET
				}
				userroles = append(userroles, _r)
			}

			userroles = rolesort(userroles)

			color := 0
			if len(m.Roles) > 0 {
			COLOR:
				for _, mr := range userroles {
					for _, gr := range guildroles {
						logger.LogDebug(d, fmt.Sprintf("Checking %s: %s vs %s w/ color %d",
							m.User.String(), gr.Name, mr.Name, mr.Color))
						if mr.ID == gr.ID && mr.Color != 0 {
							color = mr.Color
							logger.LogDebug(d, "Found: "+mr.Name)
							break COLOR
						}
					}
				}
			}

			var colorstring string
			var shortstring string

			if color != 0 {
				colorstring = fmt.Sprintf("%x", 0xFFFFFF&color)
				shortstring = strings.Join(regexp.MustCompile("^(.).(.).(.).$").
					FindStringSubmatch(colorstring)[1:], ``)
			}

			newnamecache[gid][m.User.ID] = name
			newcolorcache[gid][m.User.ID] = CachedColor{
				Integer: color,
				String:  colorstring,
				Short:   shortstring}

			if d.Loglevel() >= 3 {
				logger.LogDebug(d, fmt.Sprintf("Setting %s to use %s and %s",
					m.User.String(), newnamecache[gid][m.User.ID],
					newcolorcache[gid][m.User.ID].String))
			}
		}
	}

	// Now that the more expensive operations are done, lock the mutex, nil out
	// the old data and replace
	d.mutex.name.Lock()
	d.mutex.color.Lock()
	logger.LogDebug(d, `UpdateCache() locking`)

	d.namecache = nil
	d.colorcache = nil

	d.namecache = newnamecache
	d.colorcache = newcolorcache

	logger.LogDebug(d, `UpdateCache() unlocked`)
	d.mutex.name.Unlock()
	d.mutex.color.Unlock()
}

// GetName returns the cached nickname for that guildmember
func (d *DataCache) GetName(s *discordgo.Session, gid,
	uid string) (string, bool) {
	if gid == "" || uid == "" {
		return "", false
	}

	d.mutex.name.Lock()
	logger.LogDebug(d, `GetName() locking`)
	defer func() {
		logger.LogDebug(d, `GetName() unlocked`)
		d.mutex.name.Unlock()
	}()

	_, ok := d.namecache[gid]
	if !ok {
		return "", false
	}

	name, ok := d.namecache[gid][uid]
	if !ok {
		return "", false
	}

	return name, true
}

// GetColor returns the cached discord role color for that guildmember
func (d *DataCache) GetColor(s *discordgo.Session, gid,
	uid string) (int, string, string) {
	if gid == "" || uid == "" {
		return 0, "", ""
	}

	d.mutex.color.Lock()
	logger.LogDebug(d, `GetColor() locking`)
	defer func() {
		logger.LogDebug(d, `GetColor() unlocked`)
		d.mutex.color.Unlock()
	}()

	_, ok := d.colorcache[gid]
	if !ok {
		return 0, "", ""
	}

	data, ok := d.colorcache[gid][uid]
	if !ok {
		return 0, "", ""
	}

	return data.Integer, data.String, data.Short
}

// AddGuild adds a guild to the known guilds cache
func (d *DataCache) AddGuild(gid string) {
	if gid == "" {
		return
	}

	if _, ok := d.guildcache[gid]; ok {
		return
	}

	d.guildcache[gid] = true
	logger.LogInfo(d, fmt.Sprintf("Added Guild %s to cache", gid))
}

func rolesort(unsorted []*discordgo.Role) []*discordgo.Role {
	var sorted []*discordgo.Role

ROLESORT:
	for _, u := range unsorted {
		if len(sorted) == 0 {
			sorted = append(sorted, u)
			continue
		}

		for i, s := range sorted {
			if u.Position < s.Position && i == len(sorted)-1 {
				sorted = append(sorted, u)
				continue ROLESORT
			}

			if u.Position > s.Position {
				sorted = append(sorted, &discordgo.Role{})
				copy(sorted[i+1:], sorted[i:])
				sorted[i] = u
				continue ROLESORT
			}
		}
	}

	return sorted
}
