package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadFileToURL(filePath, uploadURL string) error {
	// 1. Open the zip file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close() // Ensure the file is closed

	// 2. Create a new PUT request
	req, err := http.NewRequest(http.MethodPut, uploadURL, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Set the Content-Type header
	// For zip files, common content types are "application/zip" or "application/octet-stream".
	// Ensure this matches what the upload URL expects.
	req.Header.Set("Content-Type", "application/zip") // Or "application/octet-stream"

	// 4. Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// 5. Check the response status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read response body for error details
		return fmt.Errorf("upload failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	fmt.Printf("File '%s' successfully uploaded to %s\n", filePath, uploadURL)
	return nil
}

func main() {
	// Replace with your actual file path and the generated upload URL
	zipFilePath := "my_local_archive.zip"
	// This uploadURL would typically be a pre-signed URL from a service like GCP Cloud Storage
	uploadURL := "YOUR_GENERATED_UPLOAD_URL_HERE"

	// Create a dummy zip file for demonstration if it doesn't exist
	if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
		fmt.Printf("Creating a dummy zip file '%s' for demonstration...\n", zipFilePath)
		dummyContent := []byte("This is some dummy content for the zip file.")
		err = os.WriteFile(zipFilePath, dummyContent, 0644) // Write directly as "zip" for this demo
		if err != nil {
			fmt.Printf("Error creating dummy file: %v\n", err)
			return
		}
		// Note: For a real zip file, you'd use archive/zip to create it.
		// For simplicity, we just write some bytes for testing the upload mechanism.
		fmt.Println("Dummy file created. For a real scenario, ensure it's a valid zip archive.")
	}

	err := uploadFileToURL(zipFilePath, uploadURL)
	if err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
	}
}
