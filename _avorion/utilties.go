package avorion

import (
	"avorioncontrol/ifaces"
	"os"
	"regexp"
)

// REVIEW: These regexp objects need to be replaced with a function that produces
// a parser object for the command that is provided. Ideally, this would be replaced
// with a *Server method that handles this.

/**
 * Substring Match Indexes:
 * 0  Entire string
 * 1  Player index
 * 2  Current coordinates (X)
 * 3  Current coordinates (Y)
 * 4  Ship Count
 * 5  Station Count
 * 6  Money
 * 7  Iron
 * 8  Titanium
 * 9  Naonite
 * 10 Trinium
 * 11 Xanian
 * 12 Ogonite
 * 13 Avorion
 * 14 Player Name
**/
var rePlayerData = regexp.MustCompile(
	`^\s*player: ([0-9]+) (-?[0-9]{1,3}):(-?[0-9]{1,3}) ([0-9]+) ([0-9]+) ` +
		`credits:(-?[0-9]+) iron:(-?[0-9]+) titanium:(-?[0-9]+) naonite:(-?[0-9]+) ` +
		`trinium:(-?[0-9]+) xanian:(-?[0-9]+) ogonite:(-?[0-9]+) avorion:(-?[0-9]+) (.*)$`)

/**
 * Substring Match Indexes:
 * 0  Entire string
 * 1  Alliance index
 * 2  Ship Count
 * 3  Station Count
 * 4  Money
 * 5  Iron
 * 6  Titanium
 * 7  Naonite
 * 8  Trinium
 * 9  Xanian
 * 10 Ogonite
 * 11 Avorion
 * 12 Alliance Name
**/
var reAllianceData = regexp.MustCompile(`^\s*alliance: ([0-9]+) ([0-9]+) ([0-9]+) ` +
	`credits:(-?[0-9]+) iron:(-?[0-9]+) titanium:(-?[0-9]+) naonite:(-?[0-9]+) ` +
	`trinium:(-?[0-9]+) xanian:(-?[0-9]+) ogonite:(-?[0-9]+) avorion:(-?[0-9]+) (.*)$`)

type jumpsByTime []ifaces.ShipCoordData

func (t jumpsByTime) Len() int {
	return len(t)
}

func (t jumpsByTime) Less(i, j int) bool {
	return t[i].Time.Unix() < t[j].Time.Unix()
}

func (t jumpsByTime) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Check if a file exists or is a directory.
func exists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
