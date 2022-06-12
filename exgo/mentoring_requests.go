package exgo

import (
	"encoding/json"
	"fmt"
	"time"
)

type MentoringRequestsResults struct {
	Results []MentoringRequest `json:"results"`
	Meta    Meta               `json:"meta"`
}

type MentoringRequest struct {
	UUID                   string    `json:"uuid"`
	TrackTitle             string    `json:"track_title"`
	ExerciseIconURL        string    `json:"exercise_icon_url"`
	ExerciseTitle          string    `json:"exercise_title"`
	StudentHandle          string    `json:"student_handle"`
	StudentAvatarURL       string    `json:"student_avatar_url"`
	UpdatedAt              time.Time `json:"updated_at"`
	HaveMentoredPreviously bool      `json:"have_mentored_previously"`
	IsFavorited            bool      `json:"is_favorited"`
	Status                 string    `json:"status"`
	TooltipURL             string    `json:"tooltip_url"`
	URL                    string    `json:"url"`
}

type Meta struct {
	CurrentPage   int `json:"current_page"`
	TotalCount    int `json:"total_count"`
	TotalPages    int `json:"total_pages"`
	UnscopedTotal int `json:"unscoped_total"`
}

func (c *Client) GetMentoringRequests(track_slug string, page int) (*MentoringRequestsResults, error) {
	req, err := c.performRequest(
		"GET",
		fmt.Sprintf("/mentoring/requests?track_slug=%s&page=%d&order=recent", track_slug, page),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("getting mentoring requests: %w", err)
	}

	var mentoringRequests MentoringRequestsResults
	err = json.Unmarshal(req, &mentoringRequests)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling mentoring requests: %w", err)
	}

	return &mentoringRequests, nil
}

func (c *Client) GetAllMentoringRequests(track_slug string) (*MentoringRequestsResults, error) {
	var mentoringRequests MentoringRequestsResults

	current_page := 1
	for {
		req, err := c.GetMentoringRequests(track_slug, current_page)
		if err != nil {
			return nil, fmt.Errorf("getting mentoring requests: %w", err)
		}

		mentoringRequests.Results = append(mentoringRequests.Results, req.Results...)
		if req.Meta.CurrentPage >= req.Meta.TotalPages {
			mentoringRequests.Meta = req.Meta
			break
		}
		current_page += 1
	}
	return &mentoringRequests, nil
}

type MentoringTracks struct {
	Tracks []MentoringTrack `json:"tracks"`
}

type MentoringTrackLinks struct {
	Exercises string `json:"exercises"`
}

type MentoringTrack struct {
	Slug               string              `json:"slug"`
	Title              string              `json:"title"`
	IconURL            string              `json:"icon_url"`
	NumSolutionsQueued int                 `json:"num_solutions_queued"`
	MedianWaitTime     int                 `json:"median_wait_time"`
	Links              MentoringTrackLinks `json:"links"`
}

func (c *Client) GetAllMentoringTracks() (*MentoringTrack, error) {
	req, err := c.performRequest(
		"GET",
		"/mentoring/tracks?page=1&per_page=100",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("getting tracks: %w", err)
	}

	var result MentoringTrack
	err = json.Unmarshal(req, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling tracks: %w", err)
	}

	return &result, nil
}
