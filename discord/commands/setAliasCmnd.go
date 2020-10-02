package commands

import (
	"avorioncontrol/ifaces"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func setaliasCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		err error
	)

	if !HasNumArgs(a, 2, 2) {
		return wrongArgsCmd(s, m, a, c)
	}

	author := m.Author.String()
	out := ""
	v := "^[a-zA-Z]{1,10}$"

	if reg, err = Registrar(m.GuildID); err != nil {
		return out, err
	}

	if !regexp.MustCompile(v).MatchString(a[2]) {
		out = "User " + author + " attempted to set an improper alias"
		msg := sprintf("Invalid alias supplied: `%s`", a[2])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	if reg.IsRegistered(a[1]) == false {
		msg := sprintf("Command supplied is not valid: `%s`", a[1])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	if err = c.SetAliasCommand(a[1], a[2]); err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Failed to configure Alias!")
		return out, err
	}

	c.SaveConfiguration()
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return out, err
}
