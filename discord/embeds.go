package discord

import (
	"avorioncontrol/ifaces"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	embedStatusStrings     map[int]string
	embedStatusColors      map[int]int
	galaxyFieldTemplate    string
	configOneFieldTemplate string
	configTwoFieldTemplate string
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

	configOneFieldTemplate = "> • **Version**: _%s_\n" +
		"> • **Seed**: _%s_\n" +
		"> \n" +
		"> • **Difficulty**: _%s_\n" +
		"> • **Collision**: _%s_\n" +
		"> • **PVP**: _%s_\n" +
		"> \n" +
		"> • **Block Limit**: _%d_\n" +
		"> • **Volume Limit**: _%d_\n"

	configTwoFieldTemplate = "> **_Players_**\n" +
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
		color          int
		ok             bool
		stat           string
		plrs           string
		statusField    *discordgo.MessageEmbedField
		galaxyField    *discordgo.MessageEmbedField
		configOneField *discordgo.MessageEmbedField
		configTwoField *discordgo.MessageEmbedField

		version      = "1"
		collision    = "1"
		name         = "Avorion Server"
		pvp          = true
		pvpString    = "Enabled"
		seed         = ""
		difLevel     = 0
		blkLimit     = int64(0)
		volLimit     = int64(0)
		pMaxSlots    = int64(0)
		pMaxShips    = int64(0)
		pMaxStations = int64(0)
		aMaxSlots    = int64(0)
		aMaxShips    = int64(0)
		aMaxStations = int64(0)
	)

	if color, ok = embedStatusColors[s.Status]; !ok {
		color = embedStatusColors[ifaces.ServerCrashed]
	}

	if stat, ok = embedStatusStrings[s.Status]; !ok {
		stat = "Invalid Server Status"
	}

	embed := discordgo.MessageEmbed{
		Type:      discordgo.EmbedTypeRich,
		Color:     color,
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    make([]*discordgo.MessageEmbedField, 0)}

	statusField = &discordgo.MessageEmbedField{
		Inline: false, Value: stat, Name: "State"}

	configOneField = &discordgo.MessageEmbedField{
		Inline: true, Name: "Server Config", Value: configOneFieldTemplate}

	configTwoField = &discordgo.MessageEmbedField{
		Inline: true, Name: "Player Config", Value: configTwoFieldTemplate}

	if s.INI != nil {
		version = s.INI.Version
		collision = s.INI.Collision
		name = s.INI.Name
		pvp = s.INI.PVP
		seed = s.INI.Seed
		difLevel = s.INI.Difficulty
		blkLimit = s.INI.BlockLimit
		volLimit = s.INI.VolumeLimit
		pMaxSlots = s.INI.MaxPlayerShips
		pMaxShips = s.INI.MaxPlayerSlots
		pMaxStations = s.INI.MaxPlayerStations
		aMaxSlots = s.INI.MaxAllianceSlots
		aMaxShips = s.INI.MaxAllianceShips
		aMaxStations = s.INI.MaxAllianceStations
	}

	if !pvp {
		pvpString = "disabled"
	}

	embed.Title = name + " Status"

	configOneField.Value = fmt.Sprintf(configOneField.Value, version,
		seed, ifaces.Difficulty(difLevel), collision, pvpString, blkLimit,
		volLimit)

	configTwoField.Value = fmt.Sprintf(configTwoField.Value, pMaxSlots,
		pMaxShips, pMaxStations, aMaxSlots, aMaxShips, aMaxStations)

	if s.PlayersOnline > 0 {
		plrs = strings.TrimSuffix(s.Players, "\n")
		plrs = strings.TrimPrefix(plrs, "\n")
		plrs = strings.ReplaceAll("\n"+plrs, "\n", "\n> • ")
	} else {
		plrs = "\n> _No players online_"
	}

	galaxyField = &discordgo.MessageEmbedField{
		Inline: false, Name: "Galaxy Information", Value: galaxyFieldTemplate}

	galaxyField.Value = fmt.Sprintf(galaxyField.Value, s.Alliances,
		s.TotalPlayers, plrs)

	embed.Fields = append(embed.Fields, statusField, configOneField,
		configTwoField, galaxyField)
	return &embed
}
