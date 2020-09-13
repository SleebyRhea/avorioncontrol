package ifaces

// ChatData describes datapassed between Discord and the Server
type ChatData struct {
	Name string
	UID  string
	Msg  string
}

// ShipCoordData describes a set of coords for a ship
type ShipCoordData struct {
	x    int
	y    int
	name string
}
