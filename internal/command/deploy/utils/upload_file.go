package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/a0dotrun/a0ctl/internal/api"
)

type ArtifactResponse struct {
	UploadURL   string `json:"uploadUrl"`
	Expires     int64  `json:"expires"`
	ArtifactKey string `json:"artifactKey"`
}

func CallArtifact(client *api.Client) (*ArtifactResponse, error) {
	body := map[string]string{
		"serverId": "01K0GPFKGP01TS8PZ376P827VZ",
		"version":  "latest",
		"suffix":   "tar.gz",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := client.Post("http://localhost:8080/v1/artifacts", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var artifactResp ArtifactResponse
	if err := json.NewDecoder(resp.Body).Decode(&artifactResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("response status: %s", resp.Status)
	return &artifactResp, nil
}

func UploadArtifact(filePath string, uploadURL string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open artifact file: %w", err)
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: status %s, response: %s", resp.Status, string(body))
	}

	log.Println("Artifact uploaded successfully")
	return nil
}
