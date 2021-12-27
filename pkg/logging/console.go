package logging

import (
	"fmt"
	"os"
)

var logLevel = os.Getenv("LOG_LEVEL")

type Logger struct {
	LogLevel int8
	debug    chan string
	info     chan string
	error    chan string
	warn     chan string
}

func (l *Logger) Debug(msg string, a ...interface{}) {
	l.debug <- fmt.Sprintf(msg, a...)
}

func (l *Logger) Info(msg string, a ...interface{}) {
	l.info <- fmt.Sprintf(msg, a...)
}

func (l *Logger) Warn(msg string, a ...interface{}) {
	l.warn <- fmt.Sprintf(msg, a...)
}

func (l *Logger) Error(msg string, a ...interface{}) {
	l.error <- fmt.Sprintf(msg, a...)
}

func getLogLevel() int8 {
	switch logLevel {
	case "error":
		return 50
	case "warn":
		return 40
	case "info":
		return 30
	case "debug":
		return 20
	default:
		return 30
	}
}

func NewLogger() *Logger {

	return &Logger{
		LogLevel: getLogLevel(),
		debug:    make(chan string, 20),
		info:     make(chan string, 20),
		error:    make(chan string, 20),
		warn:     make(chan string, 20),
	}
}

func (l *Logger) Start() {
	for {
		select {
		case msg := <-l.debug:
      if l.LogLevel == 20  {
				fmt.Printf("%sdebug: %s%s\n", FgBlue, Reset, msg)
			}

		case msg := <-l.info:
			if l.LogLevel >= 20 && l.LogLevel <= 30 {
				fmt.Printf("%sinfo: %s%s\n", FgGreen, Reset, msg)
			}

		case msg := <-l.warn:
			if l.LogLevel >= 20 && l.LogLevel <= 40 {
				fmt.Printf("%swarn: %s%s\n", FgYellow, Reset, msg)
			}

		case msg := <-l.error:
			if l.LogLevel >= 20 && l.LogLevel <= 50 {
				fmt.Printf("%serror: %s%s\n", FgRed, Reset, msg)
      }
		}
	}
}
