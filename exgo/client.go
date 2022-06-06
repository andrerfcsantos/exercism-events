package exgo

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const API_BASE_URL = "https://exercism.org/api/v2"

var ErrNoToken = errors.New("no token provided")

type Client struct {
	token      string
	httpClient *http.Client
}

func New() (*Client, error) {
	token, err := findToken()
	if err != nil {
		return nil, err
	}
	return NewWithToken(token)
}

func NewWithToken(token string) (*Client, error) {
	httpClient := &http.Client{}

	if token == "" {
		return nil, ErrNoToken
	}

	return &Client{
		httpClient: httpClient,
		token:      token,
	}, nil
}

func (c *Client) performRequest(method, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(
		method,
		API_BASE_URL+path,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("getting new request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer res.Body.Close()

	replyBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return replyBody, nil
}
