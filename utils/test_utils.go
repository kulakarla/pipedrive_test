//go:build testing

package utils

import (
	"log"
	"net/http"
)

const (
	apiToken      = "863be942d8456f146e61026f7cf69dc78efda801"
	apiTokenParam = "?api_token=" + apiToken
	baseUrl       = "https://api.pipedrive.com/v1/deals/"
)

func DeleteCreatedResourceInTests(id string) error {
	client := &http.Client{}
	targetURL := baseUrl + id + apiTokenParam

	req, err := http.NewRequest(http.MethodDelete, targetURL, nil)
	if err != nil {
		log.Printf("Failed to create DELETE request: %v", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to execute DELETE request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: got %v, want %v", resp.StatusCode, http.StatusOK)
		return err
	}

	return nil
}
