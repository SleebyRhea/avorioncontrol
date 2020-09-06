// Sections of this file were adapted from the code present under the following
// repositoty: https://github.com/bwmarrin/discordgo/blob/master/examples/
// As such, this file is subject to the terms provided by the following license:
// https://github.com/bwmarrin/discordgo/blob/master/LICENSE
// Copyright (c) 2015, Bruce Marriner
// All rights reserved.

package discord

import (
	"AvorionControl/discord/botconfig"
	"AvorionControl/gameserver"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"AvorionControl/discord/commands"
	"AvorionControl/logger"
)

// Bot is an object representing a Discord bot
type Bot struct {
	loglevel int
}

/*****************/
/* logger.Logger */
/*****************/

// SetLoglevel sets the current loglevel for the object
func (b *Bot) SetLoglevel(l int) {
	b.loglevel = l
}

// Loglevel returns the current loglevel for the object
func (b *Bot) Loglevel() int {
	return b.loglevel
}

// UUID returns the UUID for the Logger
func (b *Bot) UUID() string {
	return "Bot"
}

// Init initializes the discordgo backend
func Init(core *Bot, config *botconfig.Config, gs gameserver.Server) {
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	// Default to a user mention as the prefix
	if config.Prefix == "" {
		config.Prefix = "<@!" + dg.State.User.ID + ">"
	}

	for {
		if !dg.DataReady {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	for _, g := range dg.State.Guilds {
		onGuildJoin(g.ID, core, gs)
	}

	logger.LogInit(core, "DISCORD USER:   "+dg.State.User.String())
	logger.LogInit(core, "DISCORD PREFIX: "+config.Prefix)

	// Setup our message handler for processing commands
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		var (
			reg *commands.CommandRegistrar
			err error
		)

		if reg, err = commands.Registrar(m.GuildID); err != nil {
			onGuildJoin(m.GuildID, core, gs)
			if reg, err = commands.Registrar(m.GuildID); err != nil {
				log.Fatal(err)
			}
		}

		// Dont do anything if the User is this bot
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Disallow other bots from commanding this one
		if strings.HasPrefix(m.Author.Token, "Bot ") && !config.BotsAllowed {
			return
		}

		// Process a command if the prefix is used
		if strings.HasPrefix(m.Content, config.Prefix) {
			if err = reg.ProcessCommand(s, m, config); err != nil {
				log.Fatal(err)
			}
		}
	})
}

// onGuildJoin handler
func onGuildJoin(gid string, b *Bot, gs gameserver.Server) {
	cmdregistry := commands.NewRegistrar(gid, gs)
	cmdregistry.SetLoglevel(b.Loglevel())
	commands.InitializeCommandRegistry(cmdregistry)
	logger.LogDebug(cmdregistry, "Initialized new command registrar")
}
