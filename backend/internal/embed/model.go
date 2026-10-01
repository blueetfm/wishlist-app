package embed

// EmbedRequest is the payload for POST /api/v1/embed.
type EmbedRequest struct {
	URL string `json:"url"`
}

// EmbedResponse is Open Graph metadata scraped from EmbedRequest.URL.
type EmbedResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
