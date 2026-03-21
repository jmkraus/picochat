package backend

func newOllamaClient(baseURL string) Client {
	return &ollamaClient{baseURL: baseURL}
}
