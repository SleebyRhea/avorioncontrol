package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
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
	_, m := out.Index()
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

	if m != 1 {
		output = sprintf("%s (%d of %d)", output, out.current+1, m)
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
	m *discordgo.MessageCreate, exitch chan struct{}) {

	defer func() {
		logger.LogInfo(out, "Multi-page embed has expired")
		out = nil
	}()

	nextReact := "▶️"
	prevReact := "◀️"

	// Inactivity timer
	inactive := time.NewTimer(time.Minute)

	// Initial embed, and reactions
	embed, doP, doN := GenerateOutputEmbed(out, out.ThisPage())
	u, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)

	cid := u.ChannelID
	uid := u.ID

	if doN {
		s.MessageReactionAdd(cid, uid, nextReact)
	}

	if err != nil {
		logger.LogError(out, "discordgo: "+err.Error())
		return
	}

	for {
		select {
		case <-time.After(time.Minute * 10):
			return

		case <-inactive.C:
			return

		// Exit with the rest of the application
		case <-exitch:
			return

		case <-time.After(time.Second / 2):
			logger.LogDebug(out, "Checking for update on multi-page embed")
			m, _ := s.ChannelMessage(cid, uid)
			for _, r := range m.Reactions {
				logger.LogDebug(out, "Found emoji: "+r.Emoji.ID)
				if r.Emoji.MessageFormat() == nextReact && r.Count > 1 && r.Me {
					embed, doP, doN = GenerateOutputEmbed(out, out.NextPage())
					s.ChannelMessageEditEmbed(cid, uid, embed)
					s.MessageReactionsRemoveAll(cid, uid)

					if doP {
						s.MessageReactionAdd(cid, uid, prevReact)
					}

					if doN {
						s.MessageReactionAdd(cid, uid, nextReact)
					}

					inactive.Reset(time.Minute)
				} else if r.Emoji.MessageFormat() == prevReact && r.Count > 1 && r.Me {
					embed, doP, doN = GenerateOutputEmbed(out, out.PreviousPage())
					s.ChannelMessageEditEmbed(cid, uid, embed)
					s.MessageReactionsRemoveAll(cid, uid)

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