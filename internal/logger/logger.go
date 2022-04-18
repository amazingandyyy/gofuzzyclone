package logger

import (
	"fmt"
	"strings"
)

// Println prints a message with a color with \n
func Println(color string, message string) {
	color = strings.ToLower(color)
	var color_palettes = map[string]string{
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"reset":  "\033[0m",
	}

	fmt.Printf("%s%s%s\n", string(color_palettes[color]), message, color_palettes["reset"])
}

// Printf prints a message with a color
func Printf(color string, message string) {
	color = strings.ToLower(color)
	var color_palettes = map[string]string{
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"reset":  "\033[0m",
	}

	fmt.Printf("%s%s%s", string(color_palettes[color]), message, color_palettes["reset"])
}
