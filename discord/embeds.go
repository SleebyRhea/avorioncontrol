package discord

import (
	"avorioncontrol/ifaces"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	embedStatusStrings map[int]string
	embedStatusColors  map[int]int
)

const (
	authorName = "@SleepyFugu#3611"
	authorIcon = "https://avatars2.githubusercontent.com/u/17704274?s=400&u=3897048ff3956501c2850214d235f5ac6520dd40&v=4"
	authoURL   = "https://github.com/SleepyFugu"
)

func init() {
	// https://www.spycolor.com/web-safe-colors
	embedStatusColors = make(map[int]int, 0)
	embedStatusStrings = make(map[int]string, 0)

	embedStatusColors[ifaces.ServerOffline] = 16711680
	embedStatusColors[ifaces.ServerOnline] = 6749952
	embedStatusColors[ifaces.ServerStarting] = 3381759
	embedStatusColors[ifaces.ServerStopping] = 16776960
	embedStatusColors[ifaces.ServerCrashed] = 15158332
	embedStatusColors[ifaces.ServerRestarting] = 3381759

	embedStatusStrings[ifaces.ServerOffline] = "Offline"
	embedStatusStrings[ifaces.ServerOnline] = "Online"
	embedStatusStrings[ifaces.ServerStarting] = "Initializing"
	embedStatusStrings[ifaces.ServerStopping] = "Stopping"
	embedStatusStrings[ifaces.ServerCrashed] = "Crashed"
	embedStatusStrings[ifaces.ServerRestarting] = "Re-Initializing"
}

func generateEmbedStatus(s ifaces.ServerStatus, tz *time.Location) *discordgo.MessageEmbed {
	var (
		color       int
		ok          bool
		stat        string
		informField *discordgo.MessageEmbedField
		onlineField *discordgo.MessageEmbedField
	)

	if color, ok = embedStatusColors[s.Status]; !ok {
		color = embedStatusColors[ifaces.ServerCrashed]
	}

	if stat, ok = embedStatusStrings[s.Status]; !ok {
		stat = "Invalid Server Status"
	}

	embed := discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Current State",
		Color:       color,
		Timestamp:   time.Now().Format(time.RFC3339),
		Description: stat,
		Fields:      make([]*discordgo.MessageEmbedField, 0)}

	informField = &discordgo.MessageEmbedField{}
	informField.Inline = false
	informField.Name = "Galaxy Information"
	informField.Value = "```"
	informField.Value = fmt.Sprintf("%s\nTotal Players:   %d",
		informField.Value, s.TotalPlayers)
	informField.Value = fmt.Sprintf("%s\nTotal Alliances: %d",
		informField.Value, s.Alliances)
	informField.Value = fmt.Sprintf("%s\nOnline Players:  %d",
		informField.Value, s.PlayersOnline)
	informField.Value = fmt.Sprintf("%s\n```", informField.Value)
	embed.Fields = append(embed.Fields, informField)

	if s.PlayersOnline > 0 {
		onlineField = &discordgo.MessageEmbedField{}
		onlineField.Inline = false
		onlineField.Name = "Players Online"
		plrs := strings.TrimSuffix(s.Players, "\n")
		plrs = strings.TrimPrefix(plrs, "\n")
		onlineField.Value = strings.ReplaceAll("\n"+plrs, "\n", "\n• ")
		embed.Fields = append(embed.Fields, onlineField)
	}

	return &embed
}
