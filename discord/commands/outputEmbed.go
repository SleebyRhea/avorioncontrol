package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func generateOutputEmbed(out *CommandOutput, page *Page) (*discordgo.MessageEmbed, bool, bool) {
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
		Color:     6749952,
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
