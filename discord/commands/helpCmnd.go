package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func helpCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator) (string, error) {
	var (
		maincmd *CommandRegistrant
		reg     *CommandRegistrar
		err     error
		out     string
	)

	if reg, err = Registrar(m.GuildID); err != nil {
		return "", err
	}

	if len(a[1:]) < 1 {
		return listCmd(s, m, a, c)
	}

	if maincmd, err = reg.Command(a[1]); err != nil {
		msg := sprintf("Command `%s` doesn't exist or isn't registered", a[1])
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		return "", &ErrInvalidCommand{msg}
	}

	if c.CommandDisabled(maincmd.Name()) {
		msg := sprintf("`%s` is not a valid command", a[1])
		return "", &ErrCommandDisabled{msg}
	}

	if len(a[1:]) > 1 {
		if c, cmdlets := maincmd.Subcommands(); c > 0 {
			for _, sub := range cmdlets {
				if a[2] == sub.Name() {
					maincmd = sub
					break
				}
			}
		}

		if maincmd.Name() == a[1] {
			msg := sprintf("Subcommand `%s` doesn't exist under `%s`", a[2],
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
