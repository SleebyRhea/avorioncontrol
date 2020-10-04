package discord

import (
	"avorioncontrol/ifaces"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	embedStatusStrings  map[int]string
	embedStatusColors   map[int]int
	configFieldTemplate string
	galaxyFieldTemplate string
)

const (
	authorName = "@SleepyFugu#3611"
	authorIcon = "https://avatars2.githubusercontent.com/u/17704274?s=400&u=3897048ff3956501c2850214d235f5ac6520dd40&v=4"
	authorURL  = "https://github.com/SleepyFugu"
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

	configFieldTemplate = "> • **Seed**: _%s_\n" +
		"> • **Difficulty**: _%s_\n" +
		"> • **PVP**: _%s_\n" +
		"> • **Block Limit**: _%d_\n" +
		"> • **Volume Limit**: _%d_\n" +
		"> \n" +
		"> **_Players_**\n" +
		"> • **Max Slots**: _%d_\n" +
		"> • **Max Stations**: _%d_\n" +
		"> • **Max Ships**: _%d_\n" +
		"> \n" +
		"> **_Alliances_**\n" +
		"> • **Max Slots**: _%d_\n" +
		"> • **Max Stations**: _%d_\n" +
		"> • **Max Ships**: _%d_\n"

	galaxyFieldTemplate = "> **Alliances**: _%d_\n" +
		"> **Total Players**:  _%d_\n" +
		"> \n" +
		"> **Online Players**%s"
}

func generateEmbedStatus(s ifaces.ServerStatus, tz *time.Location) *discordgo.MessageEmbed {
	var (
		color       int
		ok          bool
		stat        string
		plrs        string
		statusField *discordgo.MessageEmbedField
		configField *discordgo.MessageEmbedField
		galaxyField *discordgo.MessageEmbedField
	)

	if color, ok = embedStatusColors[s.Status]; !ok {
		color = embedStatusColors[ifaces.ServerCrashed]
	}

	if stat, ok = embedStatusStrings[s.Status]; !ok {
		stat = "Invalid Server Status"
	}

	embed := discordgo.MessageEmbed{
		Type:      discordgo.EmbedTypeRich,
		Title:     "Current Server Status",
		Color:     color,
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    make([]*discordgo.MessageEmbedField, 0)}

	statusField = &discordgo.MessageEmbedField{
		Inline: false, Value: stat, Name: "State"}

	configField = &discordgo.MessageEmbedField{
		Inline: true, Name: "Configuration", Value: configFieldTemplate}

	configField.Value = fmt.Sprintf(configField.Value, "Value", "Insane",
		"disabled", 15000, 3000000000, 1000, 10, 10, 1000, 10, 10)

	if s.PlayersOnline > 0 {
		plrs = strings.TrimSuffix(s.Players, "\n")
		plrs = strings.TrimPrefix(plrs, "\n")
		plrs = strings.ReplaceAll("\n"+plrs, "\n", "\n> • ")
	} else {
		plrs = "\n> _No players online_"
	}

	galaxyField = &discordgo.MessageEmbedField{
		Inline: true, Name: "Galaxy Information", Value: galaxyFieldTemplate}

	galaxyField.Value = fmt.Sprintf(galaxyField.Value, s.Alliances,
		s.TotalPlayers, plrs)

	embed.Fields = append(embed.Fields, statusField, configField,
		galaxyField)
	return &embed
}
