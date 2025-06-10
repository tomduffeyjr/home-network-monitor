package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Represents a Unifi device given from the API response
type Device struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Model      string `json:"model"`
	MacAddress string `json:"macAddress"`
}

// Response represents the API response structure given by the API
type Response struct {
	Offset     int32    `json:"offset"`
	Limit      int32    `json:"limit"`
	Count      int32    `json:"count"`
	TotalCount int32    `json:"totalCount"`
	Data       []Device `json:"data"`
}

func main() {

	//The host <unifi dream machine> IP address needs to be set to the  emvironmental variable UNIFI_HOST
	//The API Token from Unifi needs to be set to the enviornmental variable UNIFI_TOKEN

	//Gets host from the environmental variable
	host := os.Getenv("UNIFI_HOST")
	if host == "" {
		fmt.Println("Error: UNIFI_HOST enviroment variable is not set")
		return
	}

	//Gets site ID from environment variable
	siteID := os.Getenv("UNIFI_SITE_ID")
	if siteID == "" {
		fmt.Println("Error: UNIFI_SITE_ID environment variable is not set")
		return
	}

	// Construct URL with site ID
	url := fmt.Sprintf("https://%s/proxy/network/integration/v1/sites/%s/devices", host, siteID)

	token := os.Getenv("UNIFI_TOKEN")
	if token == "" {
		fmt.Println("Error: UNIFI_TOKEN environment variable is not set")
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			//this tells the Go program to skip the Certificate verification even though it isn't trusted.
			//The reason from this is due to unifi using a self signed certificate
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	//Get Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	// Set headers
	req.Header.Set("X-API-KEY", token)
	req.Header.Set("Accept", "application/json")

	//Send Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	//Prints the Http Status code
	fmt.Println("HTTP Status:", resp.Status)

	defer resp.Body.Close()

	//read the body response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading resonse:", err)
		return
	}

	//Print Raw API Response
	fmt.Println("Raw response body:", string(body))

	// Parse JSON Response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Check for API error
	if resp.StatusCode != http.StatusOK {
		fmt.Println("API error: Status", resp.Status)
		return
	}

	for _, device := range response.Data {
		fmt.Printf("ID: %s, Name: %s, Model: %s, MAC: %s\n", device.ID, device.Name, device.Model, device.MacAddress)
	}

}
