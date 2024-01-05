package internal

// Location are the response when a client what know the location of an image.
type Location struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
