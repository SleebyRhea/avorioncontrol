package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"

	"github.com/bwmarrin/discordgo"
)

func setChatChannelCmnd(s *discordgo.Session, m *discordgo.MessageCreate,
	a BotArgs, c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		channels []*discordgo.Channel
		err      error
	)

	if !HasNumArgs(a, 1, 1) {
		return nil, &ErrInvalidArgument{
			message: sprintf(`%s was passed the wrong number of arguments`, cmd.Name()),
			cmd:     cmd}
	}

	if channels, err = s.GuildChannels(m.GuildID); err != nil {
		logger.LogError(cmd, err.Error())
		return nil, &ErrCommandError{
			message: "Server error getting channels",
			cmd:     cmd}
	}

	for _, dch := range channels {
		logger.LogDebug(cmd, sprintf("Checking channel ID %s against %s", dch.ID, a[1]))
		if dch.ID == a[1] && dch.Type == discordgo.ChannelTypeGuildText {
			c.SetChatChannel(a[1])
			c.SaveConfiguration()
			logger.LogInfo(cmd, sprintf(
				"%s set the chat channel to %s", m.Author.String(), dch.ID))
			return nil, nil
		}
	}

	return nil, &ErrInvalidArgument{
		message: sprintf("Invalid channel ID: `%s`", a[1]),
		cmd:     cmd}
}
