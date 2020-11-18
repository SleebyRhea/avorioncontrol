package logger

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

// var logger log.Logger
var spf = fmt.Sprintf

// func init() {
// 	log.SetOutput(ioutil.Discard)
// }

// formatResponseCode return string formatting for the given response code
func formatResponseHeader(r int, m string) string {
	white := color.FgWhite.Render
	black := color.FgBlack.Render
	greenbg := color.BgGreen.Render
	redbg := color.BgRed.Render
	bluebg := color.BgBlue.Render
	yellowbg := color.BgYellow.Render

	out := "[" + m + " " + strconv.Itoa(r) + "]"
	switch {
	case r >= 200 && r < 300:
		return greenbg(black(out))
	case r >= 300 && r < 400:
		return bluebg(black(out))
	case r >= 400 && r < 500:
		return yellowbg(black(out))
	default:
		return white(redbg(out))
	}
}

// LogOutput logs the given string with a timestamp and no prefix. Logging does
// not depend on the current loglevel of an object
func LogOutput(l ILogger, m string, chs ...chan []byte) {
	log.Output(1, spf("[%s] %s", l.UUID(), m))
}

// LogError logs an error.
func LogError(l ILogger, m string, chs ...chan []byte) {
	m = spf("[%s] [%s] %s", errorPrefix, l.UUID(), m)
	log.Output(1, m)
	sendToChans(m, chs)
}

// LogWarning logs a warning
func LogWarning(l ILogger, m string, chs ...chan []byte) {
	m = spf("[%s] [%s] %s", warnPrefix, l.UUID(), m)
	log.Output(1, m)
	sendToChans(m, chs)
}

// LogDebug logs a debug message if the loglevel of the given object is three
// or greater
func LogDebug(l ILogger, m string, chs ...chan []byte) {
	if l.Loglevel() >= debugLevel {
		log.Output(1, spf("[%s] [%s] %s", debugPrefix, l.UUID(), m))
	}
}

// LogInfo logs an informational notification if the loglevel is one or greater
func LogInfo(l ILogger, m string, chs ...chan []byte) {
	if l.Loglevel() >= infoLevel {
		log.Output(1, spf("[%s] [%s] %s", infoPrefix, l.UUID(), m))
	}
	sendToChans(spf("[%s] [%s] %s", infoPrefix, l.UUID(), m), chs)
}

// LogInit logs an initialization message
func LogInit(l ILogger, m string, chs ...chan []byte) {
	m = spf("[%s] [%s] %s", initPrefix, l.UUID(), m)
	log.Output(1, m)
	sendToChans(m, chs)
}

// LogVerbose logs a message only when the loglevel of an object is 2 or greater
func LogVerbose(l ILogger, m string, chs ...chan []byte) {
	if l.Loglevel() >= verboseLevel {
		log.Output(1, spf("[%s] [%s] %s", verbosePrefix, l.UUID(), m))
	}
}

// LogChat logs server chat
func LogChat(l ILogger, m string, chs ...chan []byte) {
	if l.Loglevel() >= infoLevel {
		log.Output(1, spf("[%s] [%s] %s", chatPrefix, l.UUID(), m))
	}
	sendToChans(spf("[%s] [%s] %s", chatPrefix, l.UUID(), m), chs)
}

// LogHTTP logs an HTTP response code and string. Provides formatting for the
// response, and will output if the loglevel of the object is 1 or greater
func LogHTTP(l ILogger, rc int, r *http.Request, chs ...chan []byte) {
	if l.Loglevel() >= infoLevel {
		rcs := formatResponseHeader(rc, r.Method)
		rinfo := spf("%s - %s %s",
			r.RemoteAddr,
			r.Host,
			r.RequestURI)
		log.Output(1, spf("[%s] %s %s", l.UUID(), rcs, rinfo))
	}
}
