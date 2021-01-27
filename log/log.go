package log

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	level  = DEBUG
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type Event struct {
	LogLevel LogLevel
	Payload  string
}

func (e *Event) Type() string {
	return e.LogLevel.String()
}

func Infof(format string, v ...interface{}) {
	event := newLog(INFO, format, v...)
	print(event)
}

func Warnf(format string, v ...interface{}) {
	event := newLog(WARNING, format, v...)
	print(event)
}

func Errorf(format string, v ...interface{}) {
	event := newLog(ERROR, format, v...)
	print(event)
}

func Debugf(format string, v ...interface{}) {
	event := newLog(DEBUG, format, v...)
	print(event)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func Level() LogLevel {
	return level
}

func SetLevel(newLevel LogLevel) {
	level = newLevel
}

func print(data *Event) {
	if data.LogLevel < level {
		return
	}

	switch data.LogLevel {
	case INFO:
		log.Infoln(data.Payload)
	case WARNING:
		log.Warnln(data.Payload)
	case ERROR:
		log.Errorln(data.Payload)
	case DEBUG:
		log.Debugln(data.Payload)
	}
}

func newLog(logLevel LogLevel, format string, v ...interface{}) *Event {
	return &Event{
		LogLevel: logLevel,
		Payload:  fmt.Sprintf(format, v...),
	}
}