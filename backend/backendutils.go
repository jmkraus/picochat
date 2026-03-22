package backend

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

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

func buildOllamaURL(baseURL, endPoint string) (string, error) {
	return buildProviderURL(baseURL, "api", endPoint)
}

func buildOpenAIURL(baseURL, endPoint string) (string, error) {
	return buildProviderURL(baseURL, "v1", endPoint)
}
