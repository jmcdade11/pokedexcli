package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jmcdade11/pokedexcli/internal/pokeapi"
)

type config struct {
	apiClient pokeapi.Client
	nextURL   *string
	prevURL   *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}

func commandHelp(config *config) error {
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

func commandExit(config *config) error {
	os.Exit(0)
	return nil
}

func commandMap(config *config) error {
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

func commandMapb(config *config) error {
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
