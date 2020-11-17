package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

// ICommandError describes an error producable by a bot command
type ICommandError interface {
	Command() *CommandRegistrant
	Subcommand() *CommandRegistrant
	Emit(*discordgo.Session, string)
	error
}

// ICommandOutput descibes an interface to a commands output
type ICommandOutput interface {
	Pages() []string
	Next() string
	Last() string
}

// BotArgs - Botarguments type (for BotCommand)
type BotArgs []string

// BotCommand - Function signature for a bots primary function
type BotCommand = func(*discordgo.Session, *discordgo.MessageCreate, BotArgs,
	ifaces.IConfigurator, *CommandRegistrant) (*CommandOutput, ICommandError)

// CommandArgument - Define an argument for a command
//  @0    Argument's invokation
//  @1    A description of it's effect on the command
type CommandArgument [2]string
