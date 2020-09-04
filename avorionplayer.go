package main

import (
	"net"
)

// AvorionPlayer - Defines a player that has connected to the server at some point
type AvorionPlayer struct {
	ip     net.IP
	name   string
	server *AvorionServer
}

// SetIP sets or updates a players IP address
func (p AvorionPlayer) SetIP(ips string) {
	p.ip = net.ParseIP(ips)
}

// IP returns the IP address that the player used to connect this session
func (p AvorionPlayer) IP() net.IP {
	return p.ip
}

// Name returns the name of the player
func (p AvorionPlayer) Name() string {
	return p.name
}

// Kick kicks the player
func (p AvorionPlayer) Kick(r string) {
	p.server.RunCommand("kick" + p.Name())
}

// Ban bans the player
func (p AvorionPlayer) Ban(r string) {
	p.server.RunCommand("ban " + p.Name())
}
