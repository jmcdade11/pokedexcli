package main

import (
	"time"

	"github.com/jmcdade11/pokedexcli/internal/pokeapi"
)

func main() {
	apiClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	pokedex := pokedex{
		entry: map[string]pokeapi.Pokemon{},
	}
	config := &config{
		apiClient: apiClient,
		pokedex:   pokedex,
	}
	startRepl(config)
}
