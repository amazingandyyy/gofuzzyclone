package prompter

import (
	"fmt"
	"strings"
)

type Prompter struct {
	color map[string]string
}

// New returns a new Prompter instance
func New() *Prompter {
	return &Prompter{
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

func (l *Prompter) Color(color string) string {
	return l.color[color]
}

func (l *Prompter) Gather(key string, variable *string) {
	fmt.Printf("%v ", key)
	_, err := fmt.Scanf("%s", variable)
	if err != nil {
		panic(err)
	}
}

func (l *Prompter) Info(message string) {
	fmt.Printf("%s%s%s", l.Color("reset"), message, l.Color("reset"))
}

func (l *Prompter) Success(message string) {
	l.Colorize("green", message)
}

func (l *Prompter) Highlight(message string) {
	l.Colorize("yellow", message)
}

func (l *Prompter) Colorize(color string, message string) {
	color = strings.ToLower(color)
	fmt.Printf("%s%s%s", l.Color(color), message, l.Color("reset"))
}
