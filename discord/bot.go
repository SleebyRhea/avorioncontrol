// Sections of this file were adapted from the code present under the following
// repositoty: https://github.com/bwmarrin/discordgo/blob/master/examples/
// As such, this file is subject to the terms provided by the following license:
// https://github.com/bwmarrin/discordgo/blob/master/LICENSE
// Copyright (c) 2015, Bruce Marriner
// All rights reserved.

package discord

import (
	"avorioncontrol/ifaces"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"avorioncontrol/discord/commands"
	"avorioncontrol/logger"
)

// Bot is an object representing a Discord bot
type Bot struct {
	processDirectMsg func(*discordgo.Session, *discordgo.MessageCreate)

	config   ifaces.IConfigurator
	session  *discordgo.Session
	chatpipe chan ifaces.ChatData
	loglevel int

	// Close goroutines
	close chan struct{}
	stop  chan struct{}
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

/****************************/
/* IFace ifaces.IBotChatter */
/****************************/

// SetChatPipe sets the current channel to pipe chats into
func (b *Bot) SetChatPipe(cd chan ifaces.ChatData) {
	b.chatpipe = cd
}

// ChatPipe returns the current channel to pipe chats into
func (b *Bot) ChatPipe() chan ifaces.ChatData {
	return b.chatpipe
}

// New returns a new instance of discord.Bot
func New(c ifaces.IConfigurator) *Bot {
	b := &Bot{
		config: c}
	b.SetLoglevel(c.Loglevel())
	return b
}

/****************************/
/* IFace ifaces.IBotStarter */
/****************************/

// Start initializes the discordgo backend
func (b *Bot) Start(gs ifaces.IGameServer) {
	logger.LogInit(b, "Initialized Discord bot")
	defer logger.LogInfo(b, "Stopped Discord bot")
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
		if gs.IsUp() && b.config.ChatChannel() == m.ChannelID {
			_, err = gs.RunCommand(fmt.Sprintf("discordsay \"%s\" \"%s\"",
				m.Author.String(), m.Content))
			if err != nil {
				s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
			} else {
				s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			}
			return
		}
	})

	go func() {
		for {
			time.Sleep(10 * time.Minute)
			gs.RunCommand(fmt.Sprintf("setdiscorddata \"%s\" \"%s\"",
				dg.State.User.String(), b.config.DiscordLink()))
		}
	}()

	dg.UpdateStatus(0, "Avorion")
	logger.LogInit(b, "DISCORD USER:   "+dg.State.User.String())
	logger.LogInit(b, "DISCORD PREFIX: "+b.config.Prefix())
}

/******************************/
/* IFace ifaces.IBotMentioner */
/******************************/

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
		defer logger.LogInfo(b, "Stopped bot chat supervisor")
		logger.LogInit(b, "Started bot chat supervisor")

		for {
			select {
			case cm := <-b.config.ChatPipe():
				logger.LogDebug(b, "Processing chat data from server")
				if b.config.ChatChannel() != "" {
					// Don't bother with empty messages
					if len(cm.Msg) == 0 {
						continue
					}

					// Default to Avorion
					if cm.Name == "" {
						cm.Name = "Avorion"
					}

					// Truncate messages larger than 1900 to make sure we have enough room
					//	for the rest of the message
					msg := string(cm.Msg)
					if len(msg) > 1900 {
						msg = msg[0:1900]
						msg += "...(truncated)"
					}

					if cm.UID != "" {
						msg = fmt.Sprintf("<@%s>: %s", cm.UID, msg)
					} else {
						msg = fmt.Sprintf("▫️ **%s**: %s", cm.Name, msg)
					}
					s.ChannelMessageSend(b.config.ChatChannel(), msg)
				}
			default:
				logger.LogInfo(b, "New channel selected")
			}
		}
	}()

	go func() {
		for {

		}
	}()

	logger.LogDebug(reg, "Initialized new command registrar")
}
