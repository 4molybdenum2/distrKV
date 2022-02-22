package db

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("default")

// is a bolt database
type Database struct {
	db *bolt.DB
}

// initialize a new db
func NewDatabase(path string) (db *Database, closeFunc func() error, err error) {

	// db connection
	boltDb, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, nil, err
	}

	db = &Database{db: boltDb}
	closeFunc = boltDb.Close

	// create default bucket for now
	if err := db.createBuckets(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("creating default bucket: %w", err)
	}

	return db, closeFunc, nil
}

// create bucket
func (d *Database) createBuckets() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(defaultBucket); err != nil {
			return err
		}
		return nil
	})
}

// set key or return error
func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value) // return err other wise return nil
	})
}

// get key
func (db *Database) GetKey(key string) ([]byte, error) {
	var result []byte

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		fmt.Printf("The value for the key is: %v\n", string(result))
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
