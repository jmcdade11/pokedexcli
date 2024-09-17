package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetLocation(id string) (LocationPokemon, error) {
	url := baseURL + "/location-area/" + id

	if val, ok := c.cache.Get(url); ok {
		locationPokemonResp := LocationPokemon{}
		err := json.Unmarshal(val, &locationPokemonResp)
		if err != nil {
			return LocationPokemon{}, err
		}
		return locationPokemonResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationPokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationPokemon{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return LocationPokemon{}, errors.New("location not found")
	}

	if resp.StatusCode != 200 {
		return LocationPokemon{}, fmt.Errorf("error: GetPokemon - %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationPokemon{}, err
	}

	locationPokemonResp := LocationPokemon{}
	err = json.Unmarshal(body, &locationPokemonResp)
	if err != nil {
		return LocationPokemon{}, err
	}
	c.cache.Add(url, body)
	return locationPokemonResp, nil
}
