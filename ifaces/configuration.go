package ifaces

import "AvorionControl/logger"

// IConfigurator describes an interface to our configuration backend
type IConfigurator interface {
	Validate() error

	IDiscordConfigurator
	ICommandConfigurator
	IGalaxyConfigurator
	IGameConfigurator
	IChatConfigurator
	logger.ILogger
}

// IDiscordConfigurator describes an interface that describes Discord configurations
type IDiscordConfigurator interface {
	Token() string
	BotUID() string
	BotMention() string
	InviteLink() string
	BotsAllowed() bool
	SetBotsAllowed(bool)
}

// IGameConfigurator describes an interface to a games configuration
type IGameConfigurator interface {
	RCONBin() string
	RCONPort() int
	DataPath() string
	RCONAddr() string
	RCONPass() string
	InstallPath() string
}

// IGalaxyConfigurator describes an interface to an object that can configure a
//	galaxy
type IGalaxyConfigurator interface {
	SetGalaxy(string) error
	Galaxy() string
}

// ICommandConfigurator describes an interface to an object that can configure
//	bot commands
type ICommandConfigurator interface {
	DisableCommand(string) error
	CommandDisabled(string) bool

	SetAliasCommand(string, string) error
	GetAliasedCommand(string) (bool, string)
	CommandAliases(string) (bool, []string)

	SetPrefix(string)
	Prefix() string
}

// IChatConfigurator describes an interface to an object that can configure chats
type IChatConfigurator interface {
	SetChatChannel(string) error
	ChatChannel() string
}
