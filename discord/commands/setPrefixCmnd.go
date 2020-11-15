package commands

import (
	"avorioncontrol/ifaces"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func setprefixCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator, cmd *CommandRegistrant) (string, ICommandError) {

	if !HasNumArgs(a, 1, 1) {
		return "", &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`, cmd.Name()),
			cmd:     cmd}
	}

	//eg: aa!, aa!!, !, !!, or <@!USERID> if mention is used
	var r = "^([a-zA-Z0-9]{0,2}[?!;:>%$#~=+-]{1,2}|mention)$"
	var msg string

	if !regexp.MustCompile(r).MatchString(a[1]) {
		return "", &ErrInvalidArgument{
			message: sprintf("Invalid prefix supplied: `%s`", a[1]),
			cmd:     cmd}
	}

	if a[1] == "mention" {
		c.SetPrefix(sprintf("<@!%s>", s.State.User.ID))
		msg = "Updated prefix to " + s.State.User.Mention()
	} else {
		c.SetPrefix(a[1])
		msg = sprintf("Updated prefix to `%s`", a[1])
	}

	c.SaveConfiguration()
	s.ChannelMessageSend(m.ChannelID, msg)
	return sprintf("User %s updated the prefix", m.Author.String()), nil
}
