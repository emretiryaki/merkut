package log

import "github.com/inconshreveable/log15"

type Level int


const(
	LevelCrit Level=iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)


var Root log15.Logger

type Logger interface {

	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})

}

func New(logger string, ctx ...interface{}) Logger {
	params := append([]interface{}{"logger", logger}, ctx...)
	return Root.New(params...)
}