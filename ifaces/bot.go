package ifaces

// IDiscordBot describes an interface to a discord bot
type IDiscordBot interface {
	IBotMentioner
	IBotChatter
	IBotStarter
}

// IBotMentioner describes an interface to a Bot that can provide a
// 	Discord mention string (@Username#Descriminator)
type IBotMentioner interface {
	Mention() string
}

// IBotStarter describes a bot that can start
type IBotStarter interface {
	Start(IGameServer, IMessageBus, IGalaxyCache)
}

// IBotChatter describes an interface to a bot that can chat
type IBotChatter interface {
	SetChatPipe(chan ChatData)
	ChatPipe() chan ChatData
}
