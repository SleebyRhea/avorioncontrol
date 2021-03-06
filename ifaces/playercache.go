package ifaces

// IPlayerCache describes the fields required to be considered a player data cache
type IPlayerCache interface {
	FromName(string) IPlayer
	FromFactionID(string) IPlayer
	FromSteam64ID(string) IPlayer
	FromDiscordID(string) IPlayer
}
