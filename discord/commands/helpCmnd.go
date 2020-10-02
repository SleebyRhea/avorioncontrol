package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func helpCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		maincmd *CommandRegistrant //Primary command being checked
		reg     *CommandRegistrar
		err     error
		out     string
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if len(a[1:]) < 1 {
		_, err = s.ChannelMessageSend(m.ChannelID, "Please provide a command")
		return "", err
	}

	if maincmd, err = reg.Command(a[1]); err != nil {
		msg := sprintf("Command `%s` doesn't exist or isn't registered", a[1])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return "", err
	}

	if c.CommandDisabled(maincmd.Name()) {
		msg := sprintf("Command `%s` is not enabled", a[1])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return "", err
	}

	if c, cmdlets := maincmd.Subcommands(); c > 0 {
		for _, cmd := range a[2:] {
		cmdletloop:
			for _, sub := range cmdlets {
				if string(cmd[0]) == sub.Name() {
					maincmd = sub
					break cmdletloop
				}
			}

			msg := sprintf("Subcommand `%s` doesn't exist under `%s`", cmd[0],
				maincmd.Name())
			_, err := s.ChannelMessageSend(m.ChannelID, msg)
			return "", err
		}
	}

	if out, err = maincmd.Help(); err != nil {
		return "", err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, out)

	return "", err
}
