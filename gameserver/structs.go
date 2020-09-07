package gameserver

// ChatData describes datapassed between Discord and the Server
type ChatData struct {
	Name string
	UID  string
	Msg  string
}
