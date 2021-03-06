package ifaces

// IDiscordBot describes an interface to a discord bot
type IDiscordBot interface {
	Start(IGameServer, IMessageBus, IGalaxyCache)
}
