package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/jmcdade11/pokedexcli/internal/pokeapi"
)

type config struct {
	apiClient pokeapi.Client
	nextURL   *string
	prevURL   *string
	pokedex   pokedex
}

type pokedex struct {
	entry map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas.  Each subsequent call displays the next 20 locations.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations areas.  Returns an error if used on the first page of results.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <location_name>",
			description: "List all of the Pokemon in a given area.",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempt to catch the given pokemon.",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "Inspect the given pokemon.",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Review the Pokedex.",
			callback:    commandPokedex,
		},
	}
}

func commandHelp(config *config, args ...string) error {
	fmt.Printf(`
Welcome to the Pokedex!
Usage:

`)
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandExit(config *config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *config, args ...string) error {
	resp, err := config.apiClient.GetLocations(config.nextURL)
	if err != nil {
		return err
	}
	config.nextURL = resp.Next
	config.prevURL = resp.Previous

	for _, result := range resp.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapb(config *config, args ...string) error {
	if config.prevURL == nil {
		return errors.New("cannot go prior to the first page")
	}

	resp, err := config.apiClient.GetLocations(config.prevURL)
	if err != nil {
		return err
	}
	config.nextURL = resp.Next
	config.prevURL = resp.Previous

	for _, result := range resp.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(config *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a location name")
	}
	name := args[0]
	fmt.Println("Exploring area...")
	resp, err := config.apiClient.GetLocation(name)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range resp.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	name := args[0]
	pokemon, err := config.apiClient.GetPokemon(name)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a ball at %s...\n", name)

	failureChance := rand.Intn(pokemon.BaseExperience)
	if failureChance > 50 {
		fmt.Printf("%s escaped!\n", name)
		return nil
	}
	fmt.Printf("%s was caught!\n", name)
	fmt.Println("You may now inspect it with the inspect command.")
	config.pokedex.entry[name] = pokemon
	return nil
}

func commandInspect(config *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	name := args[0]
	pokemon, ok := config.pokedex.entry[name]
	if !ok {
		return errors.New("you have not caught this pokemon yet")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stats := range pokemon.Stats {
		fmt.Printf("    -%s: %v\n", stats.Stat.Name, stats.BaseStat)
	}
	fmt.Println("Types:")
	for _, types := range pokemon.Types {
		fmt.Printf("    - %s", types.Type.Name)
	}
	fmt.Println()
	return nil
}

func commandPokedex(config *config, args ...string) error {
	if len(args) > 0 {
		return errors.New("this command takes no arguments")
	}
	if len(config.pokedex.entry) == 0 {
		fmt.Println("Your Pokedex is empty.")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, entry := range config.pokedex.entry {
		fmt.Printf("    - %s\n", entry.Name)
	}
	fmt.Println()
	return nil
}
