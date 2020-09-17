package commands

import (
	"avorioncontrol/ifaces"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func setprefixCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator) (string, error) {
	var (
		msg string
		out string
		p   string
	)

	if !HasNumArgs(a, 1, 1) {
		return wrongArgsCmd(s, m, a, c)
	}

	//eg: aa!, aa!!, !, !!, or <@!USERID> if mention is used
	r := "^([a-zA-Z0-9]{0,2}[?!;:>%$#~=+-]{1,2}|mention)$"
	author := m.Author.String()

	if !regexp.MustCompile(r).MatchString(a[1]) {
		msg = sprintf("Invalid prefix supplied: `%s`", a[1])
		out = "User " + author + " attempted to set an invalid prefix"
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		return out, err
	}

	if a[1] == "mention" {
		c.SetPrefix(sprintf("<@!%s>", s.State.User.ID))
		msg = sprintf("Setting prefix to %s", p)
	} else {
		c.SetPrefix(a[1])
		msg = sprintf("Setting prefix to `%s`", a[1])
	}

	err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return out, err
}
