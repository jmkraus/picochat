package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// buildProviderURL constructs a full endpoint URL for a provider.
// If the base URL already contains a path, it is trusted as-is and
// the endpoint is appended. Otherwise, the provider root path is added.
//
// Parameters:
//
//	baseURL (string) - server base URL
//	apiRoot (string) - default API root path (e.g. "api", "v1")
//	endPoint (string) - endpoint path without leading slash
//
// Returns:
//
//	string - full endpoint URL
//	error  - error if URL is invalid or endpoint is empty
func buildProviderURL(baseURL, apiRoot, endPoint string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", fmt.Errorf("invalid http url string %w", err)
	}

	cleanEndpoint := strings.Trim(strings.TrimSpace(endPoint), "/")
	if cleanEndpoint == "" {
		return "", fmt.Errorf("endpoint is empty")
	}

	basePath := strings.Trim(strings.TrimSpace(u.Path), "/")
	root := strings.Trim(strings.TrimSpace(apiRoot), "/")

	// If no base path is provided, add provider root (api/v1).
	// If a path already exists, trust it and only append endpoint.
	finalBase := basePath
	if finalBase == "" {
		finalBase = root
	}

	u.Path = "/" + path.Join(finalBase, cleanEndpoint)
	return u.String(), nil
}

// buildOllamaURL builds a full Ollama endpoint URL.
//
// Parameters:
//
//	baseURL (string) - Ollama server base URL
//	endPoint (string) - endpoint path without leading slash
//
// Returns:
//
//	string - full endpoint URL
//	error  - error if URL or endpoint is invalid
func buildOllamaURL(baseURL, endPoint string) (string, error) {
	return buildProviderURL(baseURL, "api", endPoint)
}

// buildOpenAIURL builds a full OpenAI-compatible endpoint URL.
//
// Parameters:
//
//	baseURL (string) - OpenAI-compatible server base URL
//	endPoint (string) - endpoint path without leading slash
//
// Returns:
//
//	string - full endpoint URL
//	error  - error if URL or endpoint is invalid
func buildOpenAIURL(baseURL, endPoint string) (string, error) {
	return buildProviderURL(baseURL, "v1", endPoint)
}

type openAIModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

// fetchOpenAIModels fetches model IDs from an OpenAI-compatible /models endpoint.
//
// Parameters:
//
//	baseURL (string) - OpenAI-compatible server base URL
//	apiKey (string)  - bearer token for authorization
//
// Returns:
//
//	[]string - list of model IDs
//	error    - error if request/decoding fails or API key is missing
func fetchOpenAIModels(baseURL, apiKey string) ([]string, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("missing OpenAI API key")
	}

	endpoint, err := buildOpenAIURL(baseURL, "models")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch models failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %d - %s", resp.StatusCode, string(msg))
	}

	var result openAIModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	models := make([]string, 0, len(result.Data))
	for _, v := range result.Data {
		if v.ID != "" {
			models = append(models, v.ID)
		}
	}
	return models, nil
}
