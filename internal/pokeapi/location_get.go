package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) GetLocations(pageURL *string) (LocationAreas, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if val, ok := c.cache.Get(url); ok {
		locationAreasResp := LocationAreas{}
		err := json.Unmarshal(val, &locationAreasResp)
		if err != nil {
			return LocationAreas{}, err
		}
		return locationAreasResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreas{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreas{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreas{}, err
	}

	locationAreasResp := LocationAreas{}
	err = json.Unmarshal(body, &locationAreasResp)
	if err != nil {
		return LocationAreas{}, err
	}
	c.cache.Add(url, body)
	return locationAreasResp, nil
}
