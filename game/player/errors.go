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

// ErrFailedKick is an error in which the server failed to kick a user
type ErrFailedKick struct {
	Err error
}

func (e *ErrFailedKick) Error() string {
	if e.Err == nil {
		return "failed to kick player due to an internal error"
	}
	return "failed to kick player due to an internal error: " + e.Err.Error()
}

// ErrFailedBan is an error in which the server failed to ban a user
type ErrFailedBan struct {
	Err error
}

func (e *ErrFailedBan) Error() string {
	if e.Err == nil {
		return "failed to ban player due to an internal error"
	}
	return "failed to ban player due to an internal error: " + e.Err.Error()
}
