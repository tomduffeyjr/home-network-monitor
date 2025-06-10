package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Device represents a UniFi device from the API response
type Device struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Model      string `json:"model"`
	MacAddress string `json:"macAddress"`
}

// Response represents the API response structure
type Response struct {
	Offset     int32    `json:"offset"`
	Limit      int32    `json:"limit"`
	Count      int32    `json:"count"`
	TotalCount int32    `json:"totalCount"`
	Data       []Device `json:"data"`
}

// ListDevices queries the UniFi API for a list of devices
func ListDevices(client *APIClient) ([]Device, error) {
	// Get site ID
	siteID := os.Getenv("UNIFI_SITE_ID")
	if siteID == "" {
		return nil, fmt.Errorf("UNIFI_SITE_ID environment variable is not set")
	}

	// Construct endpoint
	path := fmt.Sprintf("/proxy/network/integration/v1/sites/%s/devices", siteID)

	// Send request
	resp, err := client.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Log HTTP status
	fmt.Println("HTTP Status:", resp.Status)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	fmt.Println("Raw response body:", string(body))

	// Check for API error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: Status %s", resp.Status)
	}

	// Parse JSON
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return response.Data, nil
}
