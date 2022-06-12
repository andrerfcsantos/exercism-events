package exgo

import (
	"encoding/json"
	"fmt"
)

type TracksResult struct {
	Tracks []Track `json:"tracks"`
}

type TracksLinks struct {
	Exercises string `json:"exercises"`
}

type Track struct {
	Slug               string      `json:"slug"`
	Title              string      `json:"title"`
	IconURL            string      `json:"icon_url"`
	NumSolutionsQueued int         `json:"num_solutions_queued"`
	MedianWaitTime     int         `json:"median_wait_time"`
	Links              TracksLinks `json:"links"`
}

func (c *Client) GetAllTracks() (*TracksResult, error) {
	req, err := c.performRequest(
		"GET",
		"/tracks?page=1&per_page=100&order=unread_first",
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("getting requests: %w", err)
	}

	var result TracksResult
	err = json.Unmarshal(req, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling requests: %w", err)
	}

	return &result, nil
}
