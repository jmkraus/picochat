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
func GetAvailableModels(apiBaseURL string) ([]string, error) {
	tagsURL, err := CleanUrl(apiBaseURL, "tags")
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

func GetServerVersion(apiBaseURL string) (string, error) {
	versionURL, err := CleanUrl(apiBaseURL, "version")
	if err != nil {
		return "", fmt.Errorf("error fetching version: %w", err)
	}

	resp, err := http.Get(versionURL)
	if err != nil {
		return "", fmt.Errorf("could not fetch server version: %w", err)
	}
	defer resp.Body.Close()

	var v ServerVersion
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", fmt.Errorf("invalid server version response: %w", err)
	}

	if v.Version == "" {
		return "", fmt.Errorf("server response did not include a version")
	}

	return v.Version, nil
}

func CleanUrl(apiBaseURL, endPoint string) (string, error) {
	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = strings.TrimSuffix(u.Path, "/") + "/" + endPoint
	apiFullURL := u.String()

	return apiFullURL, nil
}
