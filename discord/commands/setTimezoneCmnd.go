package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"time"

	"github.com/bwmarrin/discordgo"
)

func setTimezoneCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		cmd *CommandRegistrant
		err error
	)

	if !HasNumArgs(a, 1, 1) {
		return wrongArgsCmd(s, m, a, c)
	}

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if cmd, err = reg.Command(a[0]); err != nil {
		return "", err
	}

	if _, err := time.LoadLocation(a[1]); err != nil {
		logger.LogWarning(cmd, sprintf("%s attempted to set an incorrect TZ: %s",
			m.Author.Mention(), a[1]))
		s.ChannelMessageSend(m.ChannelID, sprintf("Incorrect timezone: `%s`", a[1]))
		return "", nil
	}

	c.SetTimeZone(a[1])
	c.SaveConfiguration()
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return "", err
}
