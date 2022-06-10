package exgo

import (
	"encoding/json"
	"fmt"
	"time"
)

type NoticiationsResults struct {
	Results []Notification    `json:"results"`
	Meta    NotificationsMeta `json:"meta"`
}

type Notification struct {
	UUID      string    `json:"uuid"`
	URL       string    `json:"url"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	ImageType string    `json:"image_type"`
	ImageURL  string    `json:"image_url"`
	IsRead    bool      `json:"is_read"`
}

type Links struct {
	All string `json:"all"`
}

type NotificationsMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalCount  int   `json:"total_count"`
	TotalPages  int   `json:"total_pages"`
	Links       Links `json:"links"`
	UnreadCount int   `json:"unread_count"`
}

func (c *Client) GetNotifications(page int, perPage int) (*NoticiationsResults, error) {
	req, err := c.performRequest(
		"GET",
		fmt.Sprintf("/notifications?page=%d&per_page=%d&order=unread_first", page, perPage),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("getting notifications: %w", err)
	}

	var notifications NoticiationsResults
	err = json.Unmarshal(req, &notifications)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling notifications: %w", err)
	}

	return &notifications, nil
}

func (c *Client) GetAllUnreadNotifications() (*NoticiationsResults, error) {
	var notifications NoticiationsResults

	current_page := 1
outer:
	for {
		req, err := c.GetNotifications(current_page, 25)
		if err != nil {
			return nil, fmt.Errorf("getting notifications: %w", err)
		}

		for _, notification := range req.Results {
			if notification.IsRead {
				break outer
			}

			notifications.Results = append(notifications.Results, notification)
		}

		if req.Meta.CurrentPage >= req.Meta.TotalPages {
			notifications.Meta = req.Meta
			break
		}
		current_page += 1
	}

	return &notifications, nil
}
