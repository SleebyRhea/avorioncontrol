package commands

import (
	"avorioncontrol/ifaces"

	"github.com/bwmarrin/discordgo"
)

var checkingState = false

func checkHangCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		reg = cmd.Registrar()
		srv = reg.server
		out = newCommandOutput(cmd, "Server Hang Check")

		state = srv.Status().Status
	)

	out.Quoted = true

	switch {
	case checkingState:
		out.AddLine("Already checking server state")
		out.Construct()
		return out, nil
	case state == ifaces.ServerOffline:
		out.AddLine("Server is currently offline, and was taken down normally")
		out.Construct()
		return out, nil
	case state >= ifaces.ServerCrashedOffline:
		out.AddLine("Server crash state has already been detected")
		out.Construct()
		return out, nil
	case state == ifaces.ServerStarting:
		out.AddLine("Server is currently being started")
		out.Construct()
		return out, nil
	case state == ifaces.ServerRestarting:
		out.AddLine("Server is currently being restarted")
		out.Construct()
		return out, nil
	case state == ifaces.ServerStopping:
		out.AddLine("Server is being stopped")
		out.Construct()
		return out, nil
	}

	if srv.IsUp() && state == ifaces.ServerOnline {
		checkingState = true
		s.ChannelMessageSend(m.ChannelID, "Checking server state "+
			"(if its hanging this will take some time)")
		_, err := srv.RunCommand(`echo Testing server state`)
		if err != nil {
			go func() { srv.Restart(); checkingState = false }()
			out.AddLine("Server is hanging or is down, starting restart process")
		} else {
			out.AddLine("Server is online")
			checkingState = false
		}
	} else {
		out.AddLine("Server is currently offline")
	}

	out.Construct()
	return out, nil
}
