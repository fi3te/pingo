package logging

import (
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelError
)

type LogSetup struct {
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
}

func New(level LogLevel) *LogSetup {
	var ls LogSetup
	ls.Debug = log.New(io.Discard, "", 0)
	if level <= LevelDebug {
		ls.Debug.SetOutput(os.Stdout)
	}

	ls.Info = log.New(io.Discard, "", 0)
	if level <= LevelInfo {
		ls.Info.SetOutput(os.Stdout)
	}

	ls.Error = log.New(os.Stderr, "", 0)

	return &ls
}
