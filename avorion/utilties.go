package avorion

import "regexp"

// TODO: These regexp objects should be replaced with functions that parse player
// data and output a map of the players data

/**
 * Substring Match Indexes:
 * 1	Steam64 ID
 * 2  Player Index
 * 3 	Coordinates
 * 4 	Player name
 * 5 	Playtime
 * 6 	IP address
 * 7 	Credits
 * 8	Iron
 * 9	Titanium
 * 10	Naonite
 * 11 Trinium
 * 12	Xanium
 * 13	Ogonite
 * 14	Avorion
 * 15	Ships
 * 16 Stations
 **/
var rePlayerDataFull = regexp.MustCompile(
	"^([0-9]+) ([0-9]+) \\((-?[0-9]{1,3}:-?[0-9]{1,3})\\) (.+?) currently logged in, " +
		"playtime: (.+?) (" + ipv4re + "):[0-9]+ ([0-9\\.]+) Credits, " +
		"([0-9\\.]+) Iron, ([0-9\\.]+) Titanium, ([0-9\\.]+) Naonite, " +
		"([0-9\\.]+) Trinium, ([0-9\\.]+) Xanion, ([0-9\\.]+) Ogonite, " +
		"([0-9\\.]+) Avorion, ([0-9\\.]+) Ships, ([0-9\\.]+) Stations\\.?$")

/**
 * Substring Match Indexes:
 * 1	Steam64 ID
 * 2  Player Index
 * 3 	Player name
 * 4 	Credits
 * 5	Iron
 * 6	Titanium
 * 7	Naonite
 * 8  Trinium
 * 9	Xanium
 * 10	Ogonite
 * 11	Avorion
 * 12	Ships
 * 13 Stations
 **/
var rePlayerDataOffline = regexp.MustCompile(
	"^([0-9]+) ([0-9]+) (.+?) " +
		"([0-9\\.]+) Credits, ([0-9\\.]+) Iron, ([0-9\\.]+) Titanium, " +
		"([0-9\\.]+) Naonite, ([0-9\\.]+) Trinium, ([0-9\\.]+) Xanion, " +
		"([0-9\\.]+) Ogonite, ([0-9\\.]+) Avorion, " +
		"([0-9\\.]+) Ships, ([0-9\\.]+) Stations\\.?$")

var rePlayerAlliance = regexp.MustCompile(
	"^(.+) Alliance: ([0-9]+) (.+?) ([0-9\\.]+) Credits, ([0-9\\.]+) Iron, " +
		"([0-9\\.]+) Titanium, ([0-9\\.]+) Naonite, ([0-9\\.]+) Trinium, " +
		"([0-9\\.]+) Xanion, ([0-9\\.]+) Ogonite, ([0-9\\.]+) Avorion, " +
		"([0-9\\.]+) Ships, ([0-9\\.]+) Stations\\.? ([0-9\\.]+) " +
		"Members: \\(([0-9,]+)\\)$")
