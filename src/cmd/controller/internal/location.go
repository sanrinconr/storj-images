package internal

import "time"

// Location are the response when a client what know the location of an image.
type Location struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
