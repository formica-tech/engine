package internal

import (
	"log"
	"os"
)

type loggers struct {
	warn  *log.Logger
	info  *log.Logger
	error *log.Logger
	debug *log.Logger
}

func (l loggers) Debug() *log.Logger {
	return l.debug
}

func (l loggers) Info() *log.Logger {
	return l.info
}

func (l loggers) Warn() *log.Logger {
	return l.warn
}

func (l loggers) Error() *log.Logger {
	return l.error
}

type Loggers interface {
	Debug() *log.Logger
	Info() *log.Logger
	Warn() *log.Logger
	Error() *log.Logger
}

func InitLogger() Loggers {
	info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warn := log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	err := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debug := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &loggers{
		info:  info,
		debug: debug,
		warn:  warn,
		error: err,
	}
}
