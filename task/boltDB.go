package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func openDB(filepath string) (*bolt.DB, error) {
	db, err := bolt.Open(filepath, 0600, nil)
	return db, err
}

func createBucket(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
}

func putValue(db *bolt.DB, bucketName, value string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}
		var maxID int
		err := bucket.ForEach(func(k, v []byte) error {
			var id int
			_, err := fmt.Sscanf(string(k), "%d", &id)
			if err == nil && id > maxID {
				maxID = id
			}
			return nil
		})
		if err != nil {
			return err
		}

		newID := maxID + 1
		return bucket.Put([]byte(fmt.Sprintf("%d", newID)), []byte(value))
	})
}

func getValue(db *bolt.DB, bucketName, key string) (string, error) {
	var result string
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		value := bucket.Get([]byte(key))
		if value != nil {
			result = string(value)
		}
		return nil
	})
	return result, err
}

func listAllValue(db *bolt.DB, bucketName string) error {
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}
		return bucket.ForEach(func(k, v []byte) error {
			fmt.Printf("ID: %s, Description: %s\n", k, v)
			return nil
		})
	})
}

func removeValue(db *bolt.DB, bucketName, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}

		err := bucket.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("could not delete key %s: %v", key, err)
		}

		return nil
	})
}
