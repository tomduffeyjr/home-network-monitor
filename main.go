package main

import (
	"fmt"
)

func main() {
	// Initialize API client
	client, err := NewAPIClient()
	if err != nil {
		fmt.Println("Error initializing API client:", err)
		return
	}

	// Fetch device list
	devices, err := ListDevices(client)
	if err != nil {
		fmt.Println("Error listing devices:", err)
		return
	}

	// Print device information
	for _, device := range devices {
		fmt.Printf("ID: %s, Name: %s, Model: %s, MAC: %s\n", device.ID, device.Name, device.Model, device.MacAddress)
	}
}
