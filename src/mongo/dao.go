// Package mongo manage the raw connection with the mongo database to get and insert
// documents.
package mongo

import "time"

//nolint:godox // to think later
// TODO: this document dont allow the return of custom documents different
// at the dao. Think of a way of move this struct
// the problem of only move is this package must return a cursor, and this is bad.

// Document that is saved in the database.
type Document struct {
	ID               string    `bson:"_id"`
	Name             string    `bson:"name"`
	ObjectStorageKey string    `bson:"object_storage_key"`
	CreatedAt        time.Time `bson:"created_at"`
}
