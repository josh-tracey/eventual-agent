package logging

import (
	"fmt"
	"log"
	"os"
)

var logLevel = os.Getenv("LOG_LEVEL")
var disableTimestamp = os.Getenv("DISABLE_LOG_TIMESTAMPS")

type Logger struct {
	LogLevel int8
	trace    chan string
	debug    chan string
	info     chan string
	error    chan string
	warn     chan string
}

func (l *Logger) Trace(msg string, a ...interface{}) {
	l.trace <- fmt.Sprintf(msg, a...)
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
		return 10
	case "warn":
		return 20
	case "info":
		return 30
	case "debug":
		return 40
	case "trace":
		return 50
	default:
		return 30
	}
}

func NewLogger() *Logger {
	return &Logger{
		LogLevel: getLogLevel(),
		trace:    make(chan string, 4),
		debug:    make(chan string, 4),
		info:     make(chan string, 4),
		error:    make(chan string, 4),
		warn:     make(chan string, 4),
	}
}

func (l *Logger) Start() {
	for {
		select {
		case msg := <-l.trace:
			if l.LogLevel >= 50 {
				if disableTimestamp == "true" {
					fmt.Printf("%strace: %s%s\n", FgGrey, Reset, msg)
				} else {
					log.Printf("%strace: %s%s\n", FgGrey, Reset, msg)
				}
			}

		case msg := <-l.debug:
			if l.LogLevel >= 40 {
				if disableTimestamp == "true" {
					fmt.Printf("%sdebug: %s%s\n", FgBlue, Reset, msg)
				} else {
					log.Printf("%sdebug: %s%s\n", FgGreen, Reset, msg)
				}
			}

		case msg := <-l.info:
			if l.LogLevel >= 30 {
				if disableTimestamp == "true" {
					fmt.Printf("%sinfo: %s%s\n", FgGreen, Reset, msg)
				} else {
					log.Printf("%sinfo: %s%s\n", FgGreen, Reset, msg)
				}
			}

		case msg := <-l.warn:
			if l.LogLevel >= 20 {
				if disableTimestamp == "true" {
					fmt.Printf("%swarn: %s%s\n", FgYellow, Reset, msg)
				} else {
					log.Printf("%swarn: %s%s\n", FgYellow, Reset, msg)
				}
			}

		case msg := <-l.error:
			if l.LogLevel >= 10 {
				if disableTimestamp == "true" {
					fmt.Printf("%serror: %s%s\n", FgRed, Reset, msg)
				} else {
					log.Printf("%serror: %s%s\n", FgRed, Reset, msg)
				}
			}
		}
	}
}
