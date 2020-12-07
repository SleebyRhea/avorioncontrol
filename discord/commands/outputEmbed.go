package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var embedStatusColors map[int]int

func init() {
	embedStatusColors = make(map[int]int, 0)
	embedStatusColors[ifaces.CommandSuccess] = 6749952
	embedStatusColors[ifaces.CommandFailure] = 16711680
	embedStatusColors[ifaces.CommandWarning] = 16776960
}

// GenerateOutputEmbed creates a Discord embed from a CommandOutput, and outputs
// booleans depicting whether or not the Next or Previous page of output is available
func GenerateOutputEmbed(out *CommandOutput, page *Page) (*discordgo.MessageEmbed, bool, bool) {
	if out == nil || page == nil {
		log.Output(0, "Attempt to use GenerateOutputEmbed without out or page")
		log.Output(0, "Please check channel permissions")
		return nil, false, false
	}

	i, m := out.Index()
	p, n := false, false
	output := "Output"

	if out.Header != "" {
		output = out.Header
	}

	if out.last != nil {
		p = true
	}

	if out.next != nil {
		n = true
	}

	logger.LogDebug(out, sprintf("Generating embed for page %d of %d", i, m))
	if m > 0 {
		output = sprintf("%s (%d of %d)", output, i+1, m+1)
	}

	var embed = &discordgo.MessageEmbed{
		Type:      discordgo.EmbedTypeRich,
		Color:     embedStatusColors[out.Status],
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    make([]*discordgo.MessageEmbedField, 0)}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   out.Title,
		Value:  out.Description,
		Inline: false,
	})

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   output,
		Value:  page.Content,
		Inline: false,
	})

	// Only process this codepath during debug mode
	if out.Loglevel() > 2 {
		logger.LogDebug(out,
			sprintf("Generated Embed with the following data fields:\n"+
				"\tF1.Name:\"%s\"\n"+
				"\tF1.Value:\"%s\"\n"+
				"\tF2.Name:\"%s\"\n"+
				"\tF2.Value:\"%s\"\n",
				out.Title, out.Description, output, page.Content))
	}

	return embed, p, n
}

// CreatePagedEmbed is used to create paginated embeds under a goroutine that
// will eventually expire and return. The operates by paging the provided
// *CommandOutput linked list, so long as the correct react is added to the
// initial embed message. To detect which reacts need to be added, the function
// generateOutputEmbed returns two boolean values.
//
// The first boolean denotes previous, the second denotes next. These variables
// are doP and doN respectively.
func CreatePagedEmbed(out *CommandOutput, s *discordgo.Session,
	m *discordgo.MessageCreate, expirech chan struct{}, exitch chan struct{}) {

	nextReact := "▶️"
	prevReact := "◀️"

	// Inactivity timer
	inactive := time.NewTimer(time.Minute)

	// Initial embed, and reactions
	embed, doP, doN := GenerateOutputEmbed(out, out.ThisPage())
	u, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		logger.LogError(out, "discordgo: "+err.Error())
		return
	}

	cid := u.ChannelID
	uid := u.ID

	// Output some logging information and add an expired footer to the embed,
	// before updating one final time.
	defer func() {
		logger.LogInfo(out, "Multi-page embed has expired")
		m, _ := s.ChannelMessage(cid, uid)
		if m != nil {
			embed.Footer = &discordgo.MessageEmbedFooter{Text: "(expired)"}
			_, err := s.ChannelMessageEditEmbed(cid, uid, embed)
			if err != nil {
				logger.LogError(out, "discordgo: "+err.Error())
			}

			if err = s.MessageReactionsRemoveAll(cid, uid); err != nil {
				logger.LogError(out, "discordgo: "+err.Error())
			}

			out = nil
		}
	}()

	if doN {
		s.MessageReactionAdd(cid, uid, nextReact)
	}

	if err != nil {
		logger.LogError(out, "discordgo: "+err.Error())
		return
	}

	for {
		// 10 minute timeout, no matter what.
		select {
		case <-time.After(time.Minute * 10):
			return

		// Timeout if out inactivity timer ends up completing
		case <-inactive.C:
			return

		// External expiration channel
		case <-expirech:
			return

		// Exit with the rest of the application
		case <-exitch:
			return

		// Check the embed every second for new reacts
		case <-time.After(time.Second / 1):
			logger.LogDebug(out, "Checking for update on multi-page embed")
			m, err := s.ChannelMessage(cid, uid)

			// Return if there was an error
			if err != nil {
				logger.LogError(out, err.Error())
				return
			}

			// Return if the embed is nil
			if m == nil {
				return
			}

			for _, r := range m.Reactions {
				logger.LogDebug(out, "Found emoji: "+r.Emoji.Name)
				if r.Emoji.MessageFormat() == nextReact && r.Count > 1 && r.Me {
					embed, doP, doN = GenerateOutputEmbed(out, out.NextPage())
					_, err := s.ChannelMessageEditEmbed(cid, uid, embed)
					if err != nil {
						logger.LogError(out, "discordgo: "+err.Error())
					}

					err = s.MessageReactionsRemoveAll(cid, uid)
					if err != nil {
						logger.LogError(out, "discordgo: "+err.Error())
					}

					if doP {
						s.MessageReactionAdd(cid, uid, prevReact)
					}

					if doN {
						s.MessageReactionAdd(cid, uid, nextReact)
					}

					inactive.Reset(time.Minute)
				} else if r.Emoji.MessageFormat() == prevReact && r.Count > 1 && r.Me {
					embed, doP, doN = GenerateOutputEmbed(out, out.PreviousPage())
					_, err := s.ChannelMessageEditEmbed(cid, uid, embed)
					if err != nil {
						logger.LogError(out, "discordgo: "+err.Error())
					}

					err = s.MessageReactionsRemoveAll(cid, uid)
					if err != nil {
						logger.LogError(out, "discordgo: "+err.Error())
					}

					if doP {
						s.MessageReactionAdd(cid, uid, prevReact)
					}

					if doN {
						s.MessageReactionAdd(cid, uid, nextReact)
					}

					inactive.Reset(time.Minute)
				}
			}
		}
	}
}
