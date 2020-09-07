package commands

import (
	"AvorionControl/gameserver"
)

// IBotCommandableServer describes an IGameServer that can be commanded by a
//	discord.Bot
type IBotCommandableServer interface {
	gameserver.IGameServer
	gameserver.IPlayableServer
	gameserver.ICommandableServer
	gameserver.IDiscordIntegratedServer
}
