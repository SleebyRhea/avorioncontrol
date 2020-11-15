package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

// BotArgs - Botarguments type (for BotCommand)
type BotArgs []string

// BotCommand - Function signature for a bots primary function
type BotCommand = func(*discordgo.Session, *discordgo.MessageCreate, BotArgs,
	ifaces.IConfigurator, *CommandRegistrant) (string, ICommandError)

// CommandArgument - Define an argument for a command
//  @0    Argument's invokation
//  @1    A description of it's effect on the command
type CommandArgument [2]string

// ICommandError describes an error producable by a bot command
type ICommandError interface {
	Command() *CommandRegistrant
	Subcommand() *CommandRegistrant
	error
}
