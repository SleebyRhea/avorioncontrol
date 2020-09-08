package configuration

import (
	"net"
	"strconv"
)

func isPortAvailable(p int) bool {
	// https://coolaj86.com/articles/how-to-test-if-a-port-is-available-in-go/
	ps := strconv.Itoa(p)
	l, err := net.Listen("tcp", ":"+ps)
	if err != nil {
		return false
	}
	_ = l.Close()
	return true
}
