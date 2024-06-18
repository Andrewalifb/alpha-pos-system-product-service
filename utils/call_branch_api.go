package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Andrewalifb/alpha-pos-system-product-service/dto"
)

func GetPosStoreBranch(id string, token string) (*dto.BranchApiResponse, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/pos_store_branch/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request using a new http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response body to the response DTO
	var respDto dto.BranchApiResponse
	err = json.Unmarshal(respBody, &respDto)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response DTO: %w", err)
	}
	fmt.Println("BRANCH DATA:", respDto.Data.PosStoreBranch)
	return &respDto, nil
}
