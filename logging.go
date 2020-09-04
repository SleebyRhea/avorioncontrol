package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gookit/color"
)

const (
	verbosePrefix = "VERBOSE"
	debugPrefix   = "DEBUG"
	errorPrefix   = "ERROR"
	warnPrefix    = "WARN"
	infoPrefix    = "INFO"
	initPrefix    = "INIT"
	chatPrefix    = "CHAT"

	debugLevel   = 2
	verboseLevel = 2
	infoLevel    = 1
	errorLevel   = 0
	warnLevel    = 0
)

var sprintf = fmt.Sprintf

// Loggable - Interface that details an object that can log
type Loggable interface {
	Loglevel() int
	SetLoglevel(int)
	UUID() string
}

// sendToChans sends the given string to the provided list of channels
func sendToChans(out string, chs []chan []byte) {
	for _, ch := range chs {
		select {
		case ch <- []byte(out):
		default:
			log.Output(1, "Unable to send to closed channel!")
		}
	}
}

// formatResponseCode - Provide string formatting for the given response code
func formatResponseHeader(r int, m string) string {
	white := color.FgWhite.Render
	black := color.FgBlack.Render
	greenbg := color.BgGreen.Render
	redbg := color.BgRed.Render
	bluebg := color.BgBlue.Render
	yellowbg := color.BgYellow.Render

	out := "[" + m + " " + strconv.Itoa(r) + "]"
	switch true {
	case 200 <= r && r <= 299:
		return greenbg(black(out))
	case 300 <= r && r <= 399:
		return bluebg(black(out))
	case 400 <= r && r <= 499:
		return yellowbg(black(out))
	default:
		return white(redbg(out))
	}
}

// LogOutput logs the given string with a timestamp and no prefix. Logging does
// not depend on the current loglevel of an object
func LogOutput(l Loggable, m string, chs ...chan []byte) {
	log.Output(1, m)
}

// LogError logs an error.
func LogError(l Loggable, m string, chs ...chan []byte) {
	m = sprintf("[%s] %s", errorPrefix, m)
	log.Output(1, m)
	sendToChans(m, chs)
}

// LogWarning logs a warning
func LogWarning(l Loggable, m string, chs ...chan []byte) {
	m = sprintf("[%s] %s", warnPrefix, m)
	log.Output(1, m)
	sendToChans(m, chs)
}

// LogDebug logs a debug message if the loglevel of the given object is three
// or greater
func LogDebug(l Loggable, m string, chs ...chan []byte) {
	if l.Loglevel() >= debugLevel {
		log.Output(1, sprintf("[%s] %s", debugPrefix, m))
	}
}

// LogInfo logs an informational notification if the loglevel is one or greater
func LogInfo(l Loggable, m string, chs ...chan []byte) {
	if l.Loglevel() >= infoLevel {
		log.Output(1, sprintf("[%s] %s", infoPrefix, m))
	}
	sendToChans(sprintf("[%s] %s", infoPrefix, m), chs)
}

// LogInit logs an initialization message
func LogInit(l Loggable, m string, chs ...chan []byte) {
	m = sprintf("[%s] %s", initPrefix, m)
	log.Output(1, m)
	sendToChans(m, chs)
}

// LogVerbose logs a message only when the loglevel of an object is 2 or greater
func LogVerbose(l Loggable, m string, chs ...chan []byte) {
	if l.Loglevel() >= verboseLevel {
		log.Output(1, sprintf("[%s] %s", verbosePrefix, m))
	}
}

// LogChat logs server chat
func LogChat(l Loggable, m string, chs ...chan []byte) {
	if l.Loglevel() >= infoLevel {
		log.Output(1, sprintf("[%s] %s", chatPrefix, m))
	}
	sendToChans(sprintf("[%s] %s", chatPrefix, m), chs)
}

// LogHTTP logs an HTTP response code and string. Provides formatting for the
// response, and will output if the loglevel of the object is 1 or greater
func LogHTTP(l Loggable, rc int, r *http.Request, chs ...chan []byte) {
	if l.Loglevel() >= infoLevel {
		rcs := formatResponseHeader(rc, r.Method)
		rinfo := sprintf("%s - %s %s",
			r.RemoteAddr,
			r.Host,
			r.RequestURI)
		log.Output(1, sprintf("%s %s", rcs, rinfo))
	}
}
