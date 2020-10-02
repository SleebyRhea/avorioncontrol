package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func listCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		reg *CommandRegistrar
		err error
		msg string
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	_, cmnds := reg.AllCommands()
	for _, n := range cmnds {
		cmd, _ := reg.Command(n)
		if c.CommandDisabled(cmd.Name()) {
			continue
		}
		msg = sprintf("%s\n%s - %s", msg, cmd.Name(), cmd.description)
	}

	if msg != "" {
		_, err = s.ChannelMessageSend(m.ChannelID,
			sprintf("**Available Commands:**\n```\n%s\n```", msg))
		return "", err
	}

	_, err = s.ChannelMessageSend(m.ChannelID,
		sprintf("No commands available"))
	return "", err
}
