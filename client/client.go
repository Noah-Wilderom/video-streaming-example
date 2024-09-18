package client

import (
	"os"
)

type Client struct {
	BaseUrl string
}

func NewClient() *Client {
	return &Client{
		BaseUrl: os.Getenv("API_HOST"),
	}
}
