package gamedb

import (
	"avorioncontrol/ifaces"
	"database/sql"
	"errors"

	// Load sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// TrackingDB describes a database of tracked playerdata
type TrackingDB struct {
	db *sql.DB
}

// OpenTrackingDB returns a reference to a TrackingDB object given a
//	valid path to a sqlite database (or a filepath to a file that doesn't
//	exist)
func OpenTrackingDB(file string) (*TrackingDB, error) {
	var err error
	tracking := &TrackingDB{}
	if tracking.db, err = sql.Open("sqlite3", file); err != nil {
		return nil, err
	}

	tracking.Init()

	return tracking, nil
}

// Init initializes a TrackingDB object provided it has been assigned
//	a database file
func (t *TrackingDB) Init() error {
	var (
		s   *sql.Stmt
		err error
	)

	if t.db == nil {
		return errors.New("DB not assigned")
	}

	create := `BEGIN TRANSACTION;
		CREATE TABLE IF NOT EXISTS "factions" (
			"ID"     INTEGER PRIMARY KEY AUTOINCREMENT,
			"NAME"   TEXT,
			"GAMEID" INTEGER
		);
		CREATE TABLE IF NOT EXISTS "jumps" (
			"ID"        INTEGER PRIMARY KEY AUTOINCREMENT,
			"SECTOR"    INTEGER,
			"FACTION"   INTEGER,
			"SHIP NAME"	TEXT,
			"TIME"	    REAL
		);
		CREATE TABLE IF NOT EXISTS "sectors" (
			"ID" INTEGER PRIMARY KEY AUTOINCREMENT,
			"X"  INTEGER,
			"Y"  INTEGER
		);
		COMMIT;`

	if s, err = t.db.Prepare(create); err != nil {
		return err
	}

	if _, err = s.Exec(); err != nil {
		return err
	}

	return nil
}

// AddJump adds a jump to the tracking DB
func (t *TrackingDB) AddJump(si, fi, k int64, j ifaces.JumpInfo) error {
	var (
		s   *sql.Stmt
		err error
	)

	query := `INSERT INTO jumps ("SECTOR","FACTION","SHIP NAME","TIME","KIND")
		VALUES(?,?,?,?,?);`
	if s, err = t.db.Prepare(query); err != nil {
		return nil
	}

	if _, err = s.Exec(si, fi, j.Jump.Time.Unix(), j.Jump.Name, k); err != nil {
		return err
	}

	return nil
}

// TrackSector add a sector to the DB of tracked sector instances
func (t *TrackingDB) TrackSector(s *ifaces.Sector) error {
	var (
		id int64
	)

	t.db.QueryRow(`SELECT ID
		FROM	sectors
		WHERE	"X" = "?"
		AND		"Y" = "?"
		LIMIT	1`, s.X, s.Y).Scan(&id)

	if id != 0 {
		s.Index = id
		return nil
	}

	p, err := t.db.Prepare(`INSERT INTO sectors ("X", "Y") VALUES(?,?)`)
	if err != nil {
		return err
	}

	_, err = p.Exec(s.X, s.Y)
	if err != nil {
		return err
	}

	row := t.db.QueryRow(`SELECT MAX(ID) FROM sectors;`)
	row.Scan(&id)
	s.Index = id

	return nil
}
