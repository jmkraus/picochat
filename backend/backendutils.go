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

// buildProviderURL constructs a provider endpoint URL. Empty paths and the
// legacy /api and /v1 roots are accepted for backward compatibility.
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
	if root == "" {
		return "", fmt.Errorf("api root is empty")
	}

	switch basePath {
	case "", "api", "v1":
		// empty path is preferred; /api and /v1 are accepted as legacy roots
	default:
		return "", fmt.Errorf("invalid base url path %q: only empty, /api or /v1 are allowed", u.Path)
	}

	u.Path = "/" + path.Join(root, cleanEndpoint)
	return u.String(), nil
}

// buildOllamaURL builds an Ollama endpoint URL.
func buildOllamaURL(baseURL, endPoint string) (string, error) {
	return buildProviderURL(baseURL, "api", endPoint)
}

// buildOpenAIURL builds an OpenAI-compatible endpoint URL.
func buildOpenAIURL(baseURL, endPoint string) (string, error) {
	return buildProviderURL(baseURL, "v1", endPoint)
}

type openAIModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

// fetchOpenAIModels fetches model IDs from an OpenAI-compatible /models endpoint.
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
