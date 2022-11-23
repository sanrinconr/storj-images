// Package internal have the internal definitions of mongo db documents.
package internal

import "time"

// Image are the returned document.
type Image struct {
	ID               string    `bson:"_id"`
	Name             string    `bson:"name"`
	ObjectStorageKey string    `bson:"object_storage_key"`
	CreatedAt        time.Time `bson:"created_at"`
}
