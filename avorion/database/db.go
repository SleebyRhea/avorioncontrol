package database

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
	if t.db == nil {
		return errors.New("DB not assigned")
	}

	create := `BEGIN TRANSACTION;
		CREATE TABLE IF NOT EXISTS "jumps" (
			"ID"	INTEGER PRIMARY KEY AUTOINCREMENT,
			"SECTOR"	INTEGER,
			"PLAYER"	INTEGER,
			"TIME"	REAL,
			"SHIP NAME"	TEXT
		);
		CREATE TABLE IF NOT EXISTS "players" (
			"ID"	INTEGER,
			"NAME"	TEXT,
			PRIMARY KEY("ID")
		);
		CREATE TABLE IF NOT EXISTS "sectors" (
			"ID"	INTEGER PRIMARY KEY AUTOINCREMENT,
			"X"	INTEGER,
			"Y"	INTEGER
		);
		COMMIT;`

	if s, err := t.db.Prepare(create); err != nil {
		return err
	} else {
		if _, err := s.Exec(); err != nil {
			return err
		}
	}

	return nil
}

// AddJump adds a jump to the tracking DB
func (t *TrackingDB) AddJump(i, p int, j ifaces.JumpInfo) error {
	query := `BEGIN TRANSACTION;
		INSERT INTO "jumps" (?,?,?,?);
		COMMIT;`
	if s, err := t.db.Prepare(query); err != nil {
		return nil
	} else {
		if _, err := s.Exec(i, p, j.Jump.Time.Unix(), j.Jump.Name); err != nil {
			return err
		}
	}

	return nil
}
