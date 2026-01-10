package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GetAvailableModels fetches model names from the /tags endpoint.
//
// Parameters:
//
//	apiBaseURL string - Base URL of the API.
//
// Returns:
//
//	[]string - Slice of model names.
//	error - Error if fetching or decoding fails.
func GetAvailableModels(apiBaseURL string) ([]string, error) {
	tagsURL, err := BuildCleanUrl(apiBaseURL, "tags")
	if err != nil {
		return nil, fmt.Errorf("error fetching models: %w", err)
	}

	resp, err := http.Get(tagsURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %d – %s", resp.StatusCode, string(body))
	}

	var result ModelTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var names []string
	for _, model := range result.Models {
		names = append(names, model.Name)
	}

	return names, nil
}

// GetServerVersion retrieves the server version from the /version endpoint.
//
// Parameters:
//
//	apiBaseURL string - Base URL of the API.
//
// Returns:
//
//	string - Server version.
//	error - Error if request or decoding fails.
func GetServerVersion(apiBaseURL string) (string, error) {
	versionURL, err := BuildCleanUrl(apiBaseURL, "version")
	if err != nil {
		return "", fmt.Errorf("error fetching version: %w", err)
	}

	resp, err := http.Get(versionURL)
	if err != nil {
		return "", fmt.Errorf("error fetching version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("non-200 response: %d – %s", resp.StatusCode, string(body))
	}

	var v ServerVersion
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", fmt.Errorf("invalid server version response: %w", err)
	}

	if v.Version == "" {
		return "", fmt.Errorf("server response did not include a version")
	}

	return v.Version, nil
}

// BuildCleanUrl constructs a full API URL by ensuring the baseURL
// includes the /api path and appending the specified endpoint.
//
// Parameters:
//
//	apiBaseURL string - Base API URL.
//	endPoint string - API endpoint to append.
//
// Returns:
//
//	string - Full URL string.
//	error - Error if base URL is invalid.
func BuildCleanUrl(apiBaseURL, endPoint string) (string, error) {
	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid http url string %w", err)
	}

	// Remove trailing slash for unified processing
	path := strings.TrimSuffix(u.Path, "/")

	// Check if /api is already part of the path
	if !strings.HasSuffix(path, "/api") {
		// /api does not exist and must be added
		path = path + "/api"
	}

	// Add endpoint
	u.Path = fmt.Sprintf("%s/%s", path, endPoint)
	apiFullURL := u.String()

	return apiFullURL, nil
}
