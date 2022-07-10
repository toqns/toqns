package p2p

type LogLevel int

const (
	Debug LogLevel = 1 << iota
	Warning
	Info
	Error
	Fatal
	Panic
)

type Logger func(LogLevel, string, ...any)
