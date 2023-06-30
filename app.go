package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bgreen/space-traders-go/stapi"
)

func main() {

	// Read Bearer token from token.txt
	token, err := os.ReadFile("token.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "No token file found")
		return
	}

	// Add bearer token to the context
	auth := context.WithValue(context.Background(), stapi.ContextAccessToken, string(token))

	// Create API Client
	configuration := stapi.NewConfiguration()
	apiClient := stapi.NewAPIClient(configuration)

	// Get Agent from API
	resp, r, err := apiClient.AgentsApi.GetMyAgent(auth).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AgentsApi.GetMyAgent``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return
	}

	// Print full response
	fmt.Fprintf(os.Stdout, "Full Response: %v\n", resp)

	// Print specific fields
	fmt.Fprintf(os.Stdout, "Name: %v\nCredits: %v",
		resp.GetData().Symbol,
		resp.GetData().Credits)
}
