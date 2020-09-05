package avorion

import (
	"net"
)

// Player - Defines a player that has connected to the server at some point
type Player struct {
	ip     net.IP
	name   string
	server *Server
}

// SetIP sets or updates a players IP address
func (p Player) SetIP(ips string) {
	p.ip = net.ParseIP(ips)
}

// IP returns the IP address that the player used to connect this session
func (p Player) IP() net.IP {
	return p.ip
}

// Name returns the name of the player
func (p Player) Name() string {
	return p.name
}

// Kick kicks the player
func (p Player) Kick(r string) {
	p.server.RunCommand("kick" + p.Name())
}

// Ban bans the player
func (p Player) Ban(r string) {
	p.server.RunCommand("ban " + p.Name())
}
