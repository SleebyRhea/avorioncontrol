package player

// ErrPlayerNotFound is an error in which a faction ID is empty
type ErrPlayerNotFound struct {
	Ref string
}

func (e *ErrPlayerNotFound) Error() string {
	return "cannot locate player from reference: " + e.Ref
}

// ErrEmptyFactionID is an error in which a faction ID is empty
type ErrEmptyFactionID struct {
}

func (e *ErrEmptyFactionID) Error() string {
	return "faction ID string is empty"
}

// ErrEmptySteam64ID is an error in which a steam ID is empty
type ErrEmptySteam64ID struct {
}

func (e *ErrEmptySteam64ID) Error() string {
	return "Steam64 ID string is empty"
}

// ErrEmptyName is an error in which a player name is empty
type ErrEmptyName struct {
}

func (e *ErrEmptyName) Error() string {
	return "player name is empty"
}

// ErrBadSteam64ID is an error in which a steam64 ID is invalid
type ErrBadSteam64ID struct {
}

func (e *ErrBadSteam64ID) Error() string {
	return "invalid Steam64 ID string"
}

// ErrDiscordMapped is an error in which a discord ID is already mapped
type ErrDiscordMapped struct {
}

func (e *ErrDiscordMapped) Error() string {
	return "Discord user is already mapped to a player"
}
