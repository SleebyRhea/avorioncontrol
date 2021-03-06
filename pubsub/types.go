package pubsub

// RconCommand describes an RCON command to be passed to the server, with a
// return channel for errors
type RconCommand struct {
	Command   string
	Arguments []string
	Return    struct {
		Out chan string
		Err chan error
	}
}

// Close closes and deletes the channels that are created with the RconCommand
// object, and then returns itself. This is primarily used for cases where having
// channels is undesired.
//
// TODO: This is inefficient. Consider reworking this object to instead not
// create the channels at all.
func (r *RconCommand) Close() *RconCommand {
	close(r.Return.Out)
	close(r.Return.Err)
	r.Return.Out = nil
	r.Return.Err = nil
	return r
}

// NewRconCommand returns a new RconCommand from a set of strings
func NewRconCommand(command string, arguments ...string) *RconCommand {
	return &RconCommand{
		Command:   command,
		Arguments: arguments,
		Return: struct {
			Out chan string
			Err chan error
		}{
			Out: make(chan string),
			Err: make(chan error),
		}}
}

// ChatData describes datapassed between Discord and the Server
type ChatData struct {
	Name string
	UID  string
	Msg  string
}

// NewChatData returns a new ChatData from a set of strings
func NewChatData(name, uid, msg string) *ChatData {
	return &ChatData{
		Name: name,
		UID:  uid,
		Msg:  msg}
}
