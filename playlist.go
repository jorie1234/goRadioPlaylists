package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)


type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{httpClient: httpClient}
	return c
}

func (c *Client) SetBaseURL(u string) *Client {
	c.BaseURL, _ = url.Parse(u)
	return c
}

func (c *Client) GetPlaylist(path string ) (*PlayData, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data PlayData
	err = json.NewDecoder(resp.Body).Decode(&data)
	return &data, err
}
