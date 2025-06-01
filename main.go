package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TheHawk24/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(conf *pokemon_location_area) error
}

type pokemon_location_area struct {
	Count    int    `json: "count"`
	Next     string `json: "next"`
	Previous string `json: "previous"`
	Results  []struct {
		Name string `json: "name"`
		Url  string `json: "url"`
	} `json: "results"`
}

var config = &pokemon_location_area{
	Next: "https://pokeapi.co/api/v2/location-area",
}

var cache = pokecache.NewCache(5 * time.Second)

func cleanInput(text string) []string {
	low_string := strings.ToLower(text)
	trimmed := strings.TrimSpace(low_string)
	words := strings.Split(trimmed, " ")
	return words
}

func commandExit(conf *pokemon_location_area) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *pokemon_location_area) error {
	fmt.Println("Welcome to the Pokedex!\nUsage: \n")
	for key, value := range registry() {
		fmt.Printf("%v: %v\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandMap(conf *pokemon_location_area) error {
	//Check if data is cached
	data, ok := cache.Get(conf.Next)

	if !ok {
		resp, err := http.Get(conf.Next)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		cache.Add(conf.Next, data)
	}
	err := json.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	for _, v := range conf.Results {
		fmt.Println(v.Name)
	}
	return nil
}

func commandMapb(conf *pokemon_location_area) error {
	data, ok := cache.Get(conf.Next)
	if !ok {
		resp, err := http.Get(conf.Previous)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		cache.Add(conf.Previous, data)
	}

	err := json.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	for _, v := range conf.Results {
		fmt.Println(v.Name)
	}
	return nil
}

func registry() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays available commands and their usage",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next location areas which are sections of areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous location areas which ares sections of areas",
			callback:    commandMapb,
		},
	}
	return commands
}

func exec_command(command string) {
	value, ok := registry()[command]
	if ok {
		value.callback(config)
	} else {
		fmt.Println("Unknown command")
	}
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
		exec_command(first_word)
	}
}
