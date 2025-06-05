package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TheHawk24/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(conf *pokemon_location_area, args ...string) error
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

//type encounter_method struct {
//	Name string `json: "name"`
//	Url  string `json: "url"`
//}
//
//type version_details struct {
//	VersionDetails []struct {
//		Rate    int `json: "rate"`
//		Version struct {
//			Name string `json: "diamond"`
//			Url  string `json: "url"`
//		} `json: "version"`
//	} `json: "version_details"`
//}

//	type encounter_details struct {
//		Chance          int   `json: "chance"`
//		ConditionValues []any `json: "condition_values"`
//		MaxLevel        int   `json: "max_level"`
//		Method          struct {
//			Name string `json: "super-rod"`
//			Url  string `json: "url"`
//		} `json: "method"`
//		MinLevel int `json: "min_level"`
//	}

type explore_area struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	Id        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					Url  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height                 int    `json:"height"`
	HeldItems              []any  `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        any `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  any  `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []any `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

var config = &pokemon_location_area{
	Next: "https://pokeapi.co/api/v2/location-area",
}

var pokedex = make(map[string]Pokemon)

var cache = pokecache.NewCache(5 * time.Second)

func cleanInput(text string) []string {
	low_string := strings.ToLower(text)
	trimmed := strings.TrimSpace(low_string)
	words := strings.Split(trimmed, " ")
	return words
}

func commandExit(conf *pokemon_location_area, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *pokemon_location_area, args ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage: ")
	for key, value := range registry() {
		fmt.Printf("%v: %v\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandMap(conf *pokemon_location_area, args ...string) error {
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

func commandMapb(conf *pokemon_location_area, args ...string) error {
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

func commandExplore(conf *pokemon_location_area, args ...string) error {

	if len(args) == 0 {
		return errors.New("Provide name or id of a location-area")
	}

	url := "https://pokeapi.co/api/v2/location-area/" + args[0]
	fmt.Println(url)

	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		cache.Add(url, data)
	}
	var pokemon_area explore_area
	//var pokemon_area AutoGenerated
	err := json.Unmarshal(data, &pokemon_area)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %v\n", pokemon_area.Name)
	fmt.Println("Found Pokemon:")
	for _, v := range pokemon_area.PokemonEncounters {
		pokemon_name := v.Pokemon.Name
		fmt.Printf("- %v\n", pokemon_name)
	}
	return nil
}

func catch_pokemon(base_experience int) bool {

	if base_experience < 0 {
		base_experience = 5
	} else if base_experience > 100 {
		base_experience = 100
	}
	//Calculate probabilty
	result := 1 - (float64(base_experience) / 100)
	prob := 0.0
	if result > 0 {
		prob = result
	}
	random := rand.Float64()
	return random <= prob
}

func commandCatch(conf *pokemon_location_area, args ...string) error {

	if len(args) == 0 {
		return errors.New("Provide name or id of the pokemon")
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + args[0]
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var pokemon Pokemon
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)

	_, already_caught := pokedex[pokemon.Name]
	if already_caught {
		fmt.Printf("%v already caught\n", pokemon.Name)
	} else {
		ok := catch_pokemon(pokemon.BaseExperience)
		if ok {
			pokedex[pokemon.Name] = pokemon
			fmt.Printf("%v was caught!\n", pokemon.Name)
		} else {
			fmt.Printf("%v escaped!\n", pokemon.Name)
		}
	}

	return nil
}

func commandInspect(conf *pokemon_location_area, args ...string) error {
	if len(args) == 0 {
		return errors.New("Provide name or id of the pokemon")
	}

	info, ok := pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon\n")
	} else {
		fmt.Printf("Name: %v\n", info.Name)
		fmt.Printf("Height: %v\n", info.Height)
		fmt.Printf("Weight: %v\n", info.Weight)
		fmt.Printf("Stats:\n")
		for _, v := range info.Stats {
			fmt.Printf(" -%v: %v\n", v.Stat.Name, v.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, v := range info.Types {
			fmt.Printf(" - %v\n", v.Type.Name)
		}
	}

	return nil
}

func commandPokedex() error {
	fmt.Printf("Your Pokedex:\n")
	for k, _ := range pokedex {
		fmt.Printf(" - %v\n", k)
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
		"explore": {
			name:        "explore",
			description: "List pokemons located at a certain area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Get information about a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View all the pokemon caught",
		},
	}
	return commands
}

func exec_command(command string, args ...string) {
	value, ok := registry()[command]
	if ok {
		if command == "pokedex" {
			commandPokedex()
		} else {
			err := value.callback(config, args...)
			if err != nil {
				fmt.Println(fmt.Errorf("Error: %v", err))
			}
		}

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
		args := cleaned_output[1:]
		exec_command(first_word, args...)
	}
}
