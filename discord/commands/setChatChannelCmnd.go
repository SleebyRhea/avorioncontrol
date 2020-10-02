package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func setChatChannelCmnd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator) (string, error) {
	var (
		channels []*discordgo.Channel

		reg *CommandRegistrar
		cmd *CommandRegistrant
		out string
		msg string
		err error
	)

	if !HasNumArgs(a, 1, 1) {
		return wrongArgsCmd(s, m, a, c)
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return out, err
	}

	if cmd, err = reg.Command("setchatchannel"); err != nil {
		return out, err
	}

	channels, err = s.GuildChannels(m.GuildID)

	for _, dch := range channels {
		logger.LogDebug(cmd, sprintf("Checking channel ID %s against %s", dch.ID, a[1]))
		if dch.ID == a[1] && dch.Type == discordgo.ChannelTypeGuildText {
			c.SetChatChannel(a[1])
			c.SaveConfiguration()
			err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			return "", err
		}
	}

	msg = sprintf("Invalid channel ID: `%s`", a[1])
	_, err = s.ChannelMessageSend(m.ChannelID, msg)
	return out, err
}
