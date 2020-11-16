package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func setprefixCmd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {

	if !HasNumArgs(a, 1, 1) {
		return nil, &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`, cmd.Name()),
			cmd:     cmd}
	}

	//eg: aa!, aa!!, !, !!, or <@!USERID> if mention is used
	var r = "^([a-zA-Z0-9]{0,2}[?!;:>%$#~=+-]{1,2}|mention)$"
	var out = newCommandOutput(cmd, "Set Prefix")
	out.Quoted = true

	if !regexp.MustCompile(r).MatchString(a[1]) {
		return nil, &ErrInvalidArgument{
			message: sprintf("Invalid prefix supplied: `%s`", a[1]),
			cmd:     cmd}
	}

	if a[1] == "mention" {
		c.SetPrefix(sprintf("<@!%s>", s.State.User.ID))
		out.AddLine("Updated prefix to " + s.State.User.Mention())
	} else {
		c.SetPrefix(a[1])
		out.AddLine(sprintf("Updated prefix to `%s`", a[1]))
	}

	c.SaveConfiguration()
	logger.LogInfo(cmd, sprintf("User %s updated the prefix", m.Author.String()))
	out.Construct()
	return out, nil
}
