package gamedb

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	// Load sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	factionKind = [3]string{
		"player", "alliance", "npc"}
)

// TrackingDB describes a database of tracked playerdata
type TrackingDB struct {
	dbpath   string
	loglevel int
}

// New returns a reference to a TrackingDB object given a
//	valid path to a sqlite database (or a filepath to a file that doesn't
//	exist)
func New(file string) (*TrackingDB, error) {
	var (
		db  *sql.DB
		err error
	)

	t := &TrackingDB{dbpath: file}
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	return t, nil
}

// Init initializes a TrackingDB object provided it has been assigned
//	a database file
func (t *TrackingDB) Init() ([]*ifaces.Sector, error) {
	var (
		db  *sql.DB
		err error
	)

	db, err = sql.Open("sqlite3", t.dbpath)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "factions" (
		"ID"     INTEGER PRIMARY KEY AUTOINCREMENT,
		"NAME"   TEXT,
		"KIND"	 INTEGER,
		"GAMEID" INTEGER);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "jumps" (
		"ID"        INTEGER PRIMARY KEY AUTOINCREMENT,
		"SECTOR"    INTEGER,
		"FACTION"   INTEGER,
		"SHIP NAME"	TEXT,
		"TIME"	    REAL,
		"KIND"		  INTEGER);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "sectors" (
		"ID" INTEGER PRIMARY KEY AUTOINCREMENT,
		"X"  INTEGER,
		"Y"  INTEGER);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "integrations" (
		"ID" INTEGER PRIMARY KEY AUTOINCREMENT,
		"FACTION"		INTEGER,
		"DISCORD"		TEXT);`)
	if err != nil {
		return nil, err
	}

	// Get all of the sectors that have been tracked
	sectors := make([]*ifaces.Sector, 0)

	// TODO: This is only here until the major server refactor is underway.
	factions := make([]struct {
		Index int64
		Name  string
		ID    int64
		Kind  int
	}, 0)

	srows, err := db.Query(`select * from sectors;`)
	if err != nil {
		return nil, err
	}

	frows, err := db.Query(`select * from factions;`)
	if err != nil {
		return nil, err
	}

	for srows.Next() {
		sec := &ifaces.Sector{
			Jumphistory: make([]*ifaces.JumpInfo, 0)}
		srows.Scan(&sec.Index, &sec.X, &sec.Y)
		sectors = append(sectors, sec)
	}

	for frows.Next() {
		f := struct {
			Index int64
			Name  string
			ID    int64
			Kind  int
		}{}

		frows.Scan(&f.Index, &f.Name, &f.Kind, &f.ID)
		factions = append(factions, f)
	}

	srows.Close()
	frows.Close()

	var (
		jumpid    int64
		sectorid  int
		factionid int
		jumptime  float64
		kind      int
		name      string
		count     int64
	)

	for _, sec := range sectors {
		// Pull the last 100 jumps from the db for that sector
		jrows, err := db.Query(`SELECT * FROM jumps WHERE SECTOR=?
			ORDER BY rowid ASC LIMIT 100;`, sec.Index)

		if err != nil {
			return nil, err
		}

		count = 0

		for jrows.Next() {
			jrows.Scan(&jumpid, &sectorid, &factionid, &name, &jumptime, &kind)

			if err = jrows.Err(); err != nil {
				jrows.Close()
				logger.LogError(t, err.Error())
				return nil, err
			}

			if kind > len(factionKind) || kind < 0 {
				kind = 3
			}

			j := &ifaces.JumpInfo{
				Time: time.Unix(int64(jumptime), 0),
				Name: name,
				FID:  factionid,
				X:    sec.X,
				Y:    sec.Y,
				Kind: factionKind[kind]}

			sec.Jumphistory = append(sec.Jumphistory, j)
			count++
		}

		if err := jrows.Close(); err != nil {
			logger.LogError(t, err.Error())
		}

		if err := jrows.Err(); err != nil {
			return nil, err
		}

		logger.LogDebug(t, fmt.Sprintf("Loaded sector %d (%d_%d), which had %d jumps",
			sec.Index, sec.X, sec.Y, count))
	}

	return sectors, nil
}

// AddJump adds a jump to the tracking DB
func (t *TrackingDB) AddJump(si, fi, k int64, j ifaces.JumpInfo) error {
	var (
		s   *sql.Stmt
		db  *sql.DB
		err error
	)

	db, err = sql.Open("sqlite3", t.dbpath)
	if err != nil {
		return err
	}

	q := `INSERT INTO jumps ("SECTOR","FACTION","SHIP NAME","TIME","KIND")
		VALUES(?,?,?,?,?);`

	defer db.Close()

	if s, err = db.Prepare(q); err != nil {
		logger.LogError(t, fmt.Sprintf("AddJump: %s",
			err.Error()))
		return err
	}

	defer s.Close()

	if _, err = s.Exec(si, fi, j.Name, j.Time.Unix(), k); err != nil {
		logger.LogError(t, fmt.Sprintf("AddJump: %s",
			err.Error()))
		return err
	}

	logger.LogDebug(t, "AddJump: Success")
	return nil
}

// TrackSector add a sector to the DB of tracked sector instances
func (t *TrackingDB) TrackSector(sec *ifaces.Sector) error {
	var (
		db  *sql.DB
		err error
		id  int64
	)

	db, err = sql.Open("sqlite3", t.dbpath)
	if err != nil {
		return err
	}

	defer db.Close()

	db.QueryRow(`SELECT ID FROM sectors WHERE	"X" = "?" AND "Y" = "?"
		LIMIT	1`, sec.X, sec.Y).Scan(&id)

	if id != 0 {
		sec.Index = id
		return nil
	}

	p, err := db.Prepare(`INSERT INTO sectors ("X", "Y") VALUES(?,?)`)
	if err != nil {
		logger.LogError(t, fmt.Sprintf("TrackSector: %s",
			err.Error()))
		return err
	}

	_, err = p.Exec(sec.X, sec.Y)
	if err != nil {
		logger.LogError(t, fmt.Sprintf("TrackSector: %s",
			err.Error()))
		return err
	}

	row := db.QueryRow(`SELECT MAX(ID) FROM sectors;`)
	row.Scan(&id)
	sec.Index = id
	logger.LogDebug(t, "TrackSector: Added sector to DB")

	return nil
}

// TrackPlayer adds a player to the tracking DB
func (t *TrackingDB) TrackPlayer(p ifaces.IPlayer) error {
	db, err := sql.Open("sqlite3", t.dbpath)
	if err != nil {
		return err
	}

	defer db.Close()

	var (
		fid int64
		rid int64
		res sql.Result

		selQ = `SELECT ID FROM factions WHERE GAMEID=? LIMIT 1;`
		addQ = `INSERT INTO factions ("NAME","KIND","GAMEID") VALUES (?,?,?);`
	)

	fid, err = strconv.ParseInt(p.Index(), 10, 64)
	if err != nil {
		return err
	}

	logger.LogDebug(t, "Checking if player exists: "+p.Name())
	row := db.QueryRow(selQ, fid)
	row.Scan(&rid)
	if row.Err() != nil {
		return row.Err()
	}

	if rid > 0 {
		logger.LogDebug(t, fmt.Sprintf("Found player in DB: %d|%s", fid, p.Name()))
		return nil
	}

	logger.LogDebug(t, "Adding player to DB: "+p.Name())
	res, err = db.Exec(addQ, p.Name(), 0, fid)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	logger.LogDebug(t, fmt.Sprintf("Added player to DB: %d|%s", fid, p.Name()))

	return nil
}

// TrackAlliance adds an alliance to the tracking DB
func (t *TrackingDB) TrackAlliance(a ifaces.IAlliance) error {
	db, err := sql.Open("sqlite3", t.dbpath)
	if err != nil {
		return err
	}

	defer db.Close()

	var (
		fid int64
		rid int64
		res sql.Result

		selQ = `SELECT ID FROM factions WHERE GAMEID=? LIMIT 1;`
		addQ = `INSERT INTO factions ("NAME","KIND","GAMEID") VALUES (?,?,?);`
	)

	fid, err = strconv.ParseInt(a.Index(), 10, 64)
	if err != nil {
		return err
	}

	logger.LogDebug(t, "Checking if alliance exists: "+a.Name())
	row := db.QueryRow(selQ, fid)
	row.Scan(&rid)
	if row.Err() != nil {
		return row.Err()
	}

	if rid > 0 {
		logger.LogDebug(t, fmt.Sprintf("Found alliance in DB: %d|%s", fid, a.Name()))
		return nil
	}

	logger.LogDebug(t, "Adding alliance to DB: "+a.Name())
	res, err = db.Exec(addQ, a.Name(), 1, fid)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	logger.LogDebug(t, fmt.Sprintf("Added alliance to DB: %d|%s", fid, a.Name()))

	return nil
}

// AddIntegration adds a tracked integration request to our database
func (t *TrackingDB) AddIntegration(discordid string, p ifaces.IPlayer) error {
	db, err := sql.Open("sqlite3", t.dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	var (
		fid  int64
		addQ = `INSERT INTO integrations ("FACTION", "DISCORD") VALUES (?,?);`
	)

	fid, err = strconv.ParseInt(p.Index(), 10, 64)
	if err != nil {
		return err
	}

	_, err = db.Exec(addQ, fid, discordid)
	if err != nil {
		return err
	}

	logger.LogInfo(t, fmt.Sprintf("Added Discord integration for the user [%s]",
		p.Name()))
	return nil
}

// RemoveIntegration removes an existing discord integration from the database
func (t *TrackingDB) RemoveIntegration(p ifaces.IPlayer) error {
	db, err := sql.Open("sqlite3", t.dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	var (
		fid  int64
		delQ = `DELETE FROM integrations WHERE DISCORD=? AND FACTION=?;`
	)

	fid, err = strconv.ParseInt(p.Index(), 10, 64)
	_, err = db.Exec(delQ, fid, p.DiscordUID())
	if err != nil {
		return err
	}

	logger.LogInfo(t, fmt.Sprintf("Cleared Discord integration for the user [%s]",
		p.Name()))
	return nil
}

// SetDiscordToPlayer gets the Discord UID from the faction ID and sets the
// DiscordUID for the player
func (t *TrackingDB) SetDiscordToPlayer(p ifaces.IPlayer) error {
	db, err := sql.Open("sqlite3", t.dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	var (
		fid  int64
		did  string
		selQ = `SELECT DISCORD FROM integrations WHERE FACTION=? LIMIT 1;`
	)

	fid, err = strconv.ParseInt(p.Index(), 10, 64)
	row := db.QueryRow(selQ, fid)
	row.Scan(&did)
	if err := row.Err(); err != nil {
		return err
	}

	p.SetDiscordUID(did)
	logger.LogDebug(t, fmt.Sprintf("Processed integration for [%s] (%s)", p.Name(),
		did))
	return nil
}

/************************/
/* IFace logger.ILogger */
/************************/

// UUID returns the UUID of an avorion.Server
func (t *TrackingDB) UUID() string {
	return "GameDB"
}

// Loglevel returns the loglevel of an avorion.Server
func (t *TrackingDB) Loglevel() int {
	return t.loglevel
}

// SetLoglevel sets the loglevel of an avorion.Server
func (t *TrackingDB) SetLoglevel(l int) {
	t.loglevel = l
}
