package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/a0dotrun/a0ctl/internal/api"
)

type DeployResponse struct {
	DeploymentId string `json:"deploymentId"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
}

func CallDeploy(client *api.Client, artifactKey string) (*DeployResponse, error) {
	body := map[string]string{
		"serverId":    "01K0GPFKGP01TS8PZ376P827VZ",
		"target":      "DEVELOPMENT",
		"artifactKey": artifactKey,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := client.Post("http://localhost:8080/v1/deployments", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Response code: %v", resp.StatusCode)

	var deployResp DeployResponse
	if err := json.NewDecoder(resp.Body).Decode(&deployResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("response status: %s", resp.Status)
	return &deployResp, nil
}
