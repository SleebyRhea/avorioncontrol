package commands

import (
	"avorioncontrol/ifaces"
	"log"
	"time"
)

func init() {
	time.LoadLocation("America/New_York")
}

// HasNumArgs - Determine if a set of command arguments is between min and max
//  @a BotArgs    Argument set to process
//  @min int      Minimum number of positional arguments
//  @max int      Maximum number of positional arguments
//
//  You can use -1 in place of either min or max (or both) to disable the check
//  for that range.
func HasNumArgs(a BotArgs, min, max int) bool {
	if len(a) == 0 || len(a[0]) == 0 {
		log.Fatal("Empty argument list passed to commands.HasNumArgs")
		return false
	}

	if min == -1 {
		min = 0
	}

	if max == -1 {
		max = len(a[1:]) + 1
	}

	if len(a[1:]) > max || len(a[1:]) < min {
		return false
	}

	return true
}

// reverseSlice reverse an arbtrary slice
func reverseJumps(j []*ifaces.JumpInfo) []*ifaces.JumpInfo {
	var jumps []*ifaces.JumpInfo

	var l = len(j)
	var i = l - 1

	if l == 0 {
		return jumps
	}

	for {
		if i < 0 {
			break
		}
		jumps = append(jumps, j[i])
		i--
	}

	return jumps
}

func newArgument(a, b string) CommandArgument {
	return CommandArgument{a, b}
}
