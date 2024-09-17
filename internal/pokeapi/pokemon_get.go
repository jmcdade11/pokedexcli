package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetPokemon(name string) (Pokemon, error) {

	url := baseURL + "/pokemon/" + name

	if val, ok := c.cache.Get(url); ok {
		pokemonResp := Pokemon{}
		err := json.Unmarshal(val, &pokemonResp)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemonResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return Pokemon{}, errors.New("pokemon not found")
	}

	if resp.StatusCode != 200 {
		return Pokemon{}, fmt.Errorf("error: GetPokemon - %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	pokemonResp := Pokemon{}
	err = json.Unmarshal(body, &pokemonResp)
	if err != nil {
		return Pokemon{}, err
	}
	c.cache.Add(url, body)
	return pokemonResp, nil
}
