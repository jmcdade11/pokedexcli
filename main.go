package main

import (
	"time"

	"github.com/jmcdade11/pokedexcli/internal/pokeapi"
)

func main() {
	apiClient := pokeapi.NewClient(5 * time.Second)
	config := &config{
		apiClient: apiClient,
	}
	startRepl(config)
}
