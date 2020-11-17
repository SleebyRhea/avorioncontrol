package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var modURLBase = `https://steamcommunity.com/sharedfiles/filedetails/?id=`

func modAddSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {

	var (
		out     = newCommandOutput(cmd, "Add Server Mods")
		failed  = make([]string, 0)
		reason  = make([]string, 0)
		success = make([]int64, 0)
	)

	if !HasNumArgs(a[1:], 1, -1) {
		return nil, &ErrInvalidArgument{
			message: sprintf("`%s` was passed the wrong number of arguments", cmd.Name()),
			cmd:     cmd}
	}

	for _, arg := range a[2:] {
		if !regexp.MustCompile(`[0-9]{10}`).MatchString(arg) {
			return nil, &ErrInvalidArgument{
				message: sprintf("`%s` is not a valid workshop id", arg),
				cmd:     cmd}
		}
	}

	for _, mod := range a[2:] {
		if id, err := strconv.ParseInt(mod, 10, 64); err == nil {
			if err := c.AddServerMod(id); err != nil {
				failed = append(failed, mod)
				reason = append(reason, err.Error())
			} else {
				logger.LogInfo(cmd, sprintf("%s added %d to the mod configuration",
					m.Author.String(), id))
				success = append(success, id)
			}
		} else {
			failed = append(failed, mod)
			reason = append(reason, "Not a valid workshop id")
		}
	}

	out.Header = "Mods Added"
	if len(success) > 0 {
		for i, id := range success {
			out.AddLine(sprintf("> %d. %s%d", i+1, modURLBase, id))
		}
	}

	if len(failed) > 0 {
		out.Status = ifaces.CommandWarning
		out.AddLine("**Failed to Add**")
		for i, id := range failed {
			if len(reason) >= i {
				out.AddLine(sprintf("> %s: %s", id, reason[i]))
			} else {
				out.AddLine(sprintf("> %s: %s", id, "Unspecified error"))
			}
		}
	}

	if len(success) < 1 {
		return nil, &ErrCommandError{
			message: sprintf("Failed to add %d mod[s]", len(failed)),
			cmd:     cmd}
	}

	out.Construct()
	return out, nil
}

func modRemoveSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		out     = newCommandOutput(cmd, "Remove Server Mods")
		failed  = make([]string, 0)
		reason  = make([]string, 0)
		success = make([]int64, 0)
	)

	if !HasNumArgs(a[1:], 1, -1) {
		return nil, &ErrInvalidArgument{
			message: sprintf("`%s` was passed the wrong number of arguments", cmd.Name()),
			cmd:     cmd}
	}

	for _, arg := range a[2:] {
		if !regexp.MustCompile(`[0-9]{10}`).MatchString(arg) {
			return nil, &ErrInvalidArgument{
				message: sprintf("`%s` is not a valid workshop id", arg),
				cmd:     cmd}
		}
	}

	for _, mod := range a[2:] {
		if id, err := strconv.ParseInt(mod, 10, 64); err == nil {
			if err := c.RemoveServerMod(id); err != nil {
				failed = append(failed, mod)
				reason = append(reason, err.Error())
			} else {
				logger.LogInfo(cmd, sprintf("> %s removed %d from the mod configuration",
					m.Author.String(), id))
				success = append(success, id)
			}
		} else {
			failed = append(failed, mod)
			reason = append(reason, "Not a valid workshop id")
		}
	}

	out.Header = "Mods Removed"
	if len(success) > 0 {
		for i, id := range success {
			out.AddLine(sprintf("> %d. %s%d", i+1, modURLBase, id))
		}
	}

	if len(failed) > 0 {
		out.Status = ifaces.CommandWarning
		out.AddLine("**Failed to Remove**")
		for i, id := range failed {
			if len(reason) >= i {
				out.AddLine(sprintf("> %s: %s", id, reason[i]))
			} else {
				out.AddLine(sprintf("> %s: %s", id, "Unspecified error"))
			}
		}
	}

	if len(success) < 1 {
		return nil, &ErrCommandError{
			message: sprintf("Failed to remove %d mod[s]", len(failed)),
			cmd:     cmd}
	}

	out.Construct()
	return out, nil
}

func listModsSubCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	var (
		out  = newCommandOutput(cmd, "Mod List")
		mods = c.ListServerMods()
	)

	out.Quoted = true

	if len(mods) < 1 {
		out.AddLine("No mods currently configured")
		out.Construct()
		return out, nil
	}

	out.Header = "Mods Installed"
	for i, id := range mods {
		out.AddLine(sprintf("%d. %s%d", i+1, modURLBase, id))
	}

	out.Construct()
	return out, nil
}
