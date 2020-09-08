package logger

// ILogger describes an interface to an object that can log
type ILogger interface {
	UUID() string
	Loglevel() int
	SetLoglevel(int)
}
