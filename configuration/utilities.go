package configuration

import (
	"math/rand"
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

//https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
func makePass() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 128)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
