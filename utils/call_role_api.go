package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Andrewalifb/alpha-pos-system-product-service/dto"
)

func GetPosRole(id string, token string) (*dto.RoleApiResponse, error) {
	// Construct the URL for the request
	url := fmt.Sprintf("http://localhost:8080/api/v1/roles/pos_role/%s", id)

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request using a new http Client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is not 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response body into the response DTO
	var respDto dto.RoleApiResponse
	if err := json.Unmarshal(respBody, &respDto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response DTO: %w", err)
	}

	// Debugging line - Normally this would be handled by a logger in production code
	fmt.Printf("ROLE DATA: %+v\n", respDto.Data.PosRole)

	return &respDto, nil
}
