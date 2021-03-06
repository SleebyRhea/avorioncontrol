package ifaces

// IPlayerCache describes the fields required to be considered a player data cache
type IPlayerCache interface {
	NewPlayer(string, string, string, string) (IPlayer, error)

	FromName(string) IPlayer
	FromFactionID(string) IPlayer
	FromSteam64ID(string) IPlayer
	FromDiscordID(string) IPlayer
}
