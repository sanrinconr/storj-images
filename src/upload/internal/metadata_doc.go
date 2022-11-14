// Package internal manage all the representation of data to be saved
package internal

import "time"

// MetadataDoc are the document saved into the document storage.
type MetadataDoc struct {
	ID               string    `bson:"_id"`
	Name             string    `bson:"name"`
	ObjectStorageKey string    `bson:"object_storage_key"`
	CreatedAt        time.Time `bson:"created_at"`
}
