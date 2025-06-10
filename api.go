package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
)

// APIClient manages connections to the UniFi Network API
type APIClient struct {
	httpClient *http.Client
	host       string
	token      string
}

// NewAPIClient creates a new API client with the configured host and token
func NewAPIClient() (*APIClient, error) {
	// Get environment variables
	host := os.Getenv("UNIFI_HOST")
	if host == "" {
		return nil, fmt.Errorf("UNIFI_HOST environment variable is not set")
	}
	token := os.Getenv("UNIFI_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("UNIFI_TOKEN environment variable is not set")
	}

	// Create HTTP client
	// TODO: For production, trust the UDM-Pro SE's self-signed certificate instead of InsecureSkipVerify
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return &APIClient{
		httpClient: client,
		host:       host,
		token:      token,
	}, nil
}

// Request executes an HTTP request to the UniFi API
func (c *APIClient) Request(method, path string, body io.Reader) (*http.Response, error) {
	// Construct URL
	url := fmt.Sprintf("https://%s%s", c.host, path)

	// Create request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("X-API-KEY", c.token)
	req.Header.Set("Accept", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return resp, nil
}
