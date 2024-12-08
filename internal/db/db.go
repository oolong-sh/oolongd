package db

import (
	"go.etcd.io/bbolt"
)

var Database *bbolt.DB

const PinnedBucket = "PinnedNotes"
const db_path = "pinned_notes.db"

// Initialize the database and ensure the bucket exists
func InitializeDB() error {
	var err error
	Database, err = bbolt.Open(db_path, 0666, nil)
	if err != nil {
		return err
	}

	return Database.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(PinnedBucket))
		return err
	})
}

// Close the database when the server shuts down
func CloseDB() {
	if Database != nil {
		Database.Close()
	}
}
