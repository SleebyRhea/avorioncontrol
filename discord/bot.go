// Sections of this file were adapted from the code present under the following
// repository: https://github.com/bwmarrin/discordgo/blob/master/examples/
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
	"sync"
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
	exit  chan struct{}
	wg    *sync.WaitGroup
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
func New(c ifaces.IConfigurator, wg *sync.WaitGroup, exit chan struct{}) *Bot {
	b := &Bot{
		config: c,
		wg:     wg,
		exit:   exit}
	b.SetLoglevel(c.Loglevel())
	return b
}

/****************************/
/* IFace ifaces.IBotStarter */
/****************************/

// Start initializes the discordgo backend
func (b *Bot) Start(gs ifaces.IGameServer) {
	logger.LogInit(b, "Initialized Discord bot")
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
	if b.config.Prefix() == "" || b.config.Prefix() == "mention" {
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

	b.processDirectMsg = func(s *discordgo.Session, m *discordgo.MessageCreate) {
		v := regexp.MustCompile("^[0-9]+:[0-9]{10}$")
		in := strings.TrimSpace(m.Content)
		if v.MatchString(in) {
			if gs.ValidateIntegrationPin(in, m.Author.ID) {
				s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
				s.ChannelMessageSend(m.ID, "Thanks for validating!")
			}
		}
	}

	// Setup our message handler for processing commands
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		var (
			reg    *commands.CommandRegistrar
			cmdhlp *commands.CommandOutput
			cmderr commands.ICommandError
			err    error
		)

		// Stop DMs, and handle Discord integrations with Avorion
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
			if _, cmderr = reg.ProcessCommand(s, m, b.config, b.exit); cmderr != nil {
				s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
				cmderr.Emit(s, m.ChannelID)
				switch cmderr.(type) {
				case *commands.ErrInvalidArgument:
					if cmderr.Command() != nil {
						cmdhlp = cmderr.Command().Help()
					}
					if cmdhlp != nil {
						embed, _, _ := commands.GenerateOutputEmbed(cmdhlp, cmdhlp.ThisPage())
						s.ChannelMessageSendEmbed(m.ChannelID, embed)
					}

				case *commands.ErrUnauthorizedUsage:
					if cmderr.Command() != nil {
						logger.LogWarning(b, fmt.Sprintf(
							"%s attempted to run [%s], but wasn't authorized do so",
							m.Author.String(), cmderr.Command().Name()))
					}

				case *commands.ErrInvalidTimezone:
					logger.LogError(reg, cmderr.Error())

				case *commands.ErrInvalidCommand:

				case *commands.ErrCommandDisabled:

				case *commands.ErrInvalidAlias:

				case *commands.ErrCommandError:

				case *commands.ErrInvalidSubcommand:
					if cmderr.Subcommand() != nil {
						cmdhlp = cmderr.Subcommand().Help()
					} else if cmderr.Command() != nil {
						cmdhlp = cmderr.Command().Help()
					}
					if cmdhlp != nil {
						embed, _, _ := commands.GenerateOutputEmbed(cmdhlp, cmdhlp.ThisPage())
						s.ChannelMessageSendEmbed(m.ChannelID, embed)
					}

				default:
					logger.LogError(reg, cmderr.Error())
				}
			} else {
				s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
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

func (b *Bot) updateServerStatus(guild string, s *discordgo.Session,
	gs ifaces.IGameServer) {
	b.wg.Add(1)
	defer b.wg.Done()

	var (
		ok              bool
		cid             string
		lastcid         string
		laststatus      ifaces.ServerStatus
		statusmessageid string
	)

	updatechan := func(stat ifaces.ServerStatus) {
		_, err := s.Channel(cid)
		if err != nil {
			logger.LogError(b, "Discordgo: "+err.Error())
			return
		}

		tz, err := time.LoadLocation(b.config.TimeZone())
		if err != nil {
			logger.LogError(b, "Failed to load timezone from config")
			return
		}

		_, err = s.ChannelMessageEditEmbed(cid, statusmessageid,
			generateEmbedStatus(gs.Status(), tz))
		if err != nil {
			logger.LogError(b, "Discordgo: "+err.Error())
		}
	}

	setupchan := func(stat ifaces.ServerStatus, clear bool) {
		cid, ok = b.config.StatusChannel()
		if ok {
			logger.LogInit(b, "Setting up server status on channel: "+cid)

			if clear {
				messages, err := s.ChannelMessages(cid, 100, "", "", "")
				if err != nil {
					logger.LogError(b, "Discordgo: "+err.Error())
					return
				}

				messageids := make([]string, 0)
				for _, m := range messages {
					messageids = append(messageids, m.ID)
				}
				s.ChannelMessagesBulkDelete(cid, messageids)

				if err != nil {
					logger.LogError(b, "Discordgo: "+err.Error())
					return
				}
			}

			tz, err := time.LoadLocation(b.config.TimeZone())
			if err != nil {
				logger.LogError(b, "Failed to load timezone from config")
				return
			}

			m, err := s.ChannelMessageSendEmbed(cid, generateEmbedStatus(
				stat, tz))
			if err != nil {
				logger.LogError(b, "Discordgo: "+err.Error())
				return
			}

			statusmessageid = m.ID
			lastcid = cid
		}
	}

	// TODO: Once autostart configuration is in place, we should modify
	//	this so that it *doesnt* change the status to initializing
	logger.LogInit(b, "Starting server updater for "+guild)
	laststatus = gs.Status()
	laststatus.Status = ifaces.ServerStarting
	setupchan(laststatus, b.config.StatusChannelClear())

	// Make sure that we change the embed to state that the server is
	// offline when we exit
	defer func() {
		laststatus.Status = ifaces.ServerOffline
		laststatus.PlayersOnline = 0
		updatechan(laststatus)
	}()

	for {
		select {
		case <-b.exit:
			laststatus.Status = ifaces.ServerStopping
			updatechan(laststatus)
			for gs.IsUp() {
				time.Sleep(time.Second)
			}
			return

		case <-time.After(time.Second * 5):
			// No point in continuing if the server status hasn't changed
			if stat := gs.Status(); gs.CompareStatus(stat, laststatus) {
				continue
			} else {
				logger.LogInfo(b, "Server status updated")
				laststatus = stat
			}

			if cid, ok := b.config.StatusChannel(); ok {
				if lastcid != cid {
					setupchan(gs.Status(), b.config.StatusChannelClear())
					continue
				}

				updatechan(laststatus)
			}
		}
	}
}

/*********/
/* OTHER */
/*********/

// onGuildJoin handler
func onGuildJoin(gid string, s *discordgo.Session, b *Bot,
	gs ifaces.IGameServer) {

	reg := commands.NewRegistrar(gid, gs)
	reg.SetLoglevel(b.Loglevel())
	commands.InitializeCommandRegistry(reg)

	go func() {
		logger.LogInit(b, "Started bot chat supervisor")
		b.wg.Add(1)

		defer func() {
			b.wg.Done()
			logger.LogInfo(b, "Stopped bot chat supervisor")
		}()

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
			case <-b.exit:
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	go b.updateServerStatus(gid, s, gs)

	logger.LogDebug(reg, "Initialized new command registrar")
}
