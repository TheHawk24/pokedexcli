package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	low_string := strings.ToLower(text)
	trimmed := strings.TrimSpace(low_string)
	words := strings.Split(trimmed, " ")
	return words
}

func main() {
	fmt.Println("Hello, World!")
}
