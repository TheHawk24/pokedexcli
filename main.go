package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	low_string := strings.ToLower(text)
	trimmed := strings.TrimSpace(low_string)
	words := strings.Split(trimmed, " ")
	return words
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		var cleaned_output []string
		token := scanner.Scan()
		if token {
			user_input := scanner.Text()
			output := cleanInput(user_input)
			cleaned_output = output
		} else {
			continue
		}
		first_word := cleaned_output[0]
		fmt.Printf("Your command was: %v\n", first_word)
	}
}
