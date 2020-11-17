package commands

import (
	"avorioncontrol/ifaces"
	"avorioncontrol/logger"
	"avorioncontrol/randstring"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/gabriel-vasile/mimetype"
)

func sendBroadcastCmnd(s *discordgo.Session, m *discordgo.MessageCreate, a BotArgs,
	c ifaces.IConfigurator, cmd *CommandRegistrant) (*CommandOutput, ICommandError) {
	if len(m.Attachments) < 1 {
		return nil, &ErrCommandError{
			message: "Please attach a text file to process", cmd: cmd}
	}

	if !HasNumArgs(a, 1, -1) {
		return nil, &ErrInvalidArgument{
			message: "Please provide a message subject",
			cmd:     cmd}
	}

	var (
		dir  = c.DataPath() + c.Galaxy() + "/messages/"
		sub  = strings.Join(a[1:], " ")
		name = m.Attachments[0].Filename
		size = m.Attachments[0].Size
		url  = m.Attachments[0].URL
		out  = newCommandOutput(cmd, "Mass In-Game Email")
		srv  = cmd.Registrar().server
	)

	if utf8.RuneCountInString(sub) > 48 {
		return nil, &ErrInvalidArgument{
			message: "Subject is too long. Must be 48 characters or less.",
			cmd:     cmd}
	}

	if size >= 32768 {
		return nil, &ErrCommandError{
			message: "Attachment is too large! Please send a text file <32Kb in size",
			cmd:     cmd}
	}

	errout := &ErrCommandError{
		message: "Failed to send email, please check the logs",
		cmd:     cmd}

	logger.LogDebug(cmd, "Ensuring that messages directory exists")
	_, err := os.Stat(c.DataPath() + c.Galaxy() + "/messages")
	if os.IsNotExist(err) {
		if os.Mkdir(c.DataPath()+c.Galaxy()+"/messages", 0700) != nil {
			logger.LogError(cmd, "Failed to create messages directory: "+err.Error())
			return nil, errout
		}
	} else if err != nil {
		logger.LogError(cmd, "Failed to create tmp directory: "+err.Error())
		return nil, errout
	}

	logger.LogDebug(cmd, "Creating random temporary file")
	tmp := randstring.New(16)

	logger.LogDebug(cmd, "Acquiring attachment: "+url)
	resp, err := http.Get(url)
	if err != nil {
		logger.LogError(cmd, err.Error())
		return nil, errout
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(cmd, err.Error())
		return nil, errout
	}

	err = ioutil.WriteFile(dir+tmp, content, 0644)
	defer os.Remove(dir + tmp)
	if err != nil {
		logger.LogError(cmd, err.Error())
		return nil, errout
	}

	mime, err := mimetype.DetectFile(dir + tmp)
	if err != nil {
		logger.LogError(cmd, err.Error())
		return nil, errout
	}

	if !mime.Is("text/plain") {
		return nil, &ErrCommandError{
			message: sprintf("Invalid content type. Expected _text/plain_, got _%s_",
				mime.String()),
			cmd: cmd}
	}

	ret, err := srv.RunCommand(sprintf("sendmail -b -h \"%s\" -s \"%s\" -f \"%s\"",
		sub, m.Author.String(), tmp))
	if err != nil {
		return nil, &ErrCommandError{
			message: "Sendmail returned an error: " + err.Error(),
			cmd:     cmd}
	}

	out.Header = "Result"
	out.AddLine("> " + ret)
	out.AddLine("**Message Information**")
	out.AddLine(sprintf("> Message Size: _%s_", humanize.Bytes(uint64(size))))
	out.AddLine(sprintf("> Message URL: [_%s_](%s)", name, url))
	out.Construct()

	return out, nil
}
