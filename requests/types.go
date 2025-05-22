package requests

type ModelTag struct {
	Name string `json:"name"`
}

type ModelTagsResponse struct {
	Models []ModelTag `json:"models"`
}
