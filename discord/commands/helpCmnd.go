package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

func helpCmd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		maincmd *CommandRegistrant
		reg     = cmd.Registrar()
		authlvl = 0
		err     error
	)

	if len(a[1:]) < 1 && a[0] == "help" {
		return listCmd(s, m, a, c, cmd)
	}

	if a[0] == "help" {
		if maincmd, err = reg.Command(a[1]); err != nil {
			return nil, &ErrInvalidCommand{
				name: a[1],
				cmd:  cmd}
		}
	} else {
		a = a[0:1]
		maincmd = cmd
	}

	// Get the users authorization level
	member, _ := s.GuildMember(reg.GuildID, m.Author.ID)
	for _, r := range member.Roles {
		if l := c.GetRoleAuth(r); l > authlvl {
			authlvl = l
		}
	}

	authreq := c.GetCmndAuth(maincmd.Name())

	if c.CommandDisabled(maincmd.Name()) || authlvl < authreq {
		return nil, &ErrCommandDisabled{cmd: cmd}
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
			return nil, &ErrInvalidSubcommand{
				subname: a[2],
				cmd:     maincmd}
		}
	}

	return maincmd.Help(), nil
}
