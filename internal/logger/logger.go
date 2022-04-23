package logger

import (
	"fmt"
	"strings"
)

type Logger struct {
	color map[string]string
}

// New returns a new Logger instance
func New() *Logger {
	return &Logger{
		color: map[string]string{
			"black":   "\033[30m",
			"red":     "\033[31m",
			"green":   "\033[32m",
			"yellow":  "\033[33m",
			"blue":    "\033[34m",
			"magenta": "\033[35m",
			"cyan":    "\033[36m",
			"white":   "\033[37m",
			"reset":   "\033[0m",
		},
	}
}

func (l *Logger) HandleError(e error) {
	if e != nil {
		l.Error(e.Error())
		panic(e)
	}
}

func (l *Logger) Color(color string) string {
	return l.color[color]
}

func (l *Logger) Info(message string) {
	fmt.Printf("%s%s%s\n", l.Color("reset"), message, l.Color("reset"))
}

func (l *Logger) Success(message string) {
	l.Colorize("green", message)
}

func (l *Logger) Error(message string) {
	l.Colorize("red", message)
}

func (l *Logger) Warn(message string) {
	l.Colorize("yellow", message)
}

func (l *Logger) Colorize(color string, message string) {
	color = strings.ToLower(color)
	fmt.Printf("%s%s%s\n", l.Color(color), message, l.Color("reset"))
}
