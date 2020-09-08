// Sections of this file were adapted from the code present under the following
// repositoty: https://github.com/bwmarrin/discordgo/blob/master/examples/
// As such, this file is subject to the terms provided by the following license:
// https://github.com/bwmarrin/discordgo/blob/master/LICENSE
// Copyright (c) 2015, Bruce Marriner
// All rights reserved.

package discord

import (
	"AvorionControl/ifaces"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"AvorionControl/discord/commands"
	"AvorionControl/logger"
)

// Bot is an object representing a Discord bot
type Bot struct {
	processDirectMsg func(*discordgo.Session, *discordgo.MessageCreate)

	config   ifaces.IConfigurator
	session  *discordgo.Session
	loglevel int
}

/************************/
/* IFace logger.ILogger */
/************************/

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

// New returns a new instance of discord.Bot
func New(c ifaces.IConfigurator) *Bot {
	b := &Bot{
		config: c}
	b.SetLoglevel(c.Loglevel())
	return b
}

// Start initializes the discordgo backend
func (b *Bot) Start(gs ifaces.IGameServer) {
	dg, err := discordgo.New("Bot " + b.config.Token())
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	// Default to a user mention as the prefix
	if b.config.Prefix() == "" {
		b.config.SetPrefix(fmt.Sprintf("<@!%s>", dg.State.User.ID))
	}

	for {
		if !dg.DataReady {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	for _, g := range dg.State.Guilds {
		onGuildJoin(g.ID, dg, b, gs)
	}

	b.processDirectMsg = func(s *discordgo.Session,
		m *discordgo.MessageCreate) {
		v := regexp.MustCompile("^[0-9]+:[0-9]{6}$")
		in := strings.TrimSpace(m.Content)

		if !v.MatchString(in) {
			return
		}

		if gs.ValidateIntegrationPin(in, m.Author.ID) {
			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
		}
	}

	// Setup our message handler for processing commands
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		var (
			reg *commands.CommandRegistrar
			err error
		)

		// Stop DMs
		if m.GuildID == "" {
			b.processDirectMsg(s, m)
			return
		}

		if reg, err = commands.Registrar(m.GuildID); err != nil {
			onGuildJoin(m.GuildID, dg, b, gs)
			if reg, err = commands.Registrar(m.GuildID); err != nil {
				log.Fatal(err)
			}
		}

		// Dont do anything if the User is this bot
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Disallow other bots from commanding this one
		if strings.HasPrefix(m.Author.Token, "Bot ") && !b.config.BotsAllowed() {
			return
		}

		// Process a command if the prefix is used
		if strings.HasPrefix(m.Content, b.config.Prefix()) {
			if err = reg.ProcessCommand(s, m, b.config); err != nil {
				logger.LogError(reg, err.Error())
			}
			return
		}

		// Send messages from Discord to the ifaces as the user if its available
		if gs.IsUp() && b.config.ChatChannel() != "" {
			_, err = gs.RunCommand(fmt.Sprintf("say [%s] %s", m.Author.String(), m.Content))
			if err != nil {
				s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
			} else {
				s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			}
			return
		}

		logger.LogInit(b, "DISCORD USER:   "+dg.State.User.String())
		logger.LogInit(b, "DISCORD PREFIX: "+b.config.Prefix())
	})
}

// Mention returns the sesions bot mention
func (b *Bot) Mention() string {
	return b.session.State.User.String()
}

// onGuildJoin handler
func onGuildJoin(gid string, s *discordgo.Session, b *Bot,
	gs ifaces.IGameServer) {

	reg := commands.NewRegistrar(gid, gs)
	reg.SetLoglevel(b.Loglevel())
	commands.InitializeCommandRegistry(reg)

	go func() {
		for {
			select {
			case cm := <-gs.DCOutput():
				if b.config.ChatChannel() != "" {
					msg := string(cm.Msg)
					if cm.UID != "" {
						msg = fmt.Sprintf("<@%s>: %s", cm.UID, msg)
					} else {
						msg = fmt.Sprintf("**%s**: %s", cm.Name, msg)
					}
					s.ChannelMessageSend(b.config.ChatChannel(), msg)
				}
			}
		}
	}()

	logger.LogDebug(reg, "Initialized new command registrar")
}
