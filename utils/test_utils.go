//go:build testing

package utils

import (
	"log"
	"net/http"
	"pdrive/pipedrive-test-api/config"
)

// DeleteCreatedResourceInTests is a helper method to delete the created deals in tests when testing POST/PUT requests
func DeleteCreatedResourceInTests(id string) error {
	client := &http.Client{}
	targetURL := config.BaseURL + id + config.APITokenParam

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
