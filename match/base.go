package match

import (
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
)

var db *bolt.DB
var loadedMatches map[uint64]bool

func init() {
	loadedMatches = make(map[uint64]bool)
	var err error
	db, err = bolt.Open("match.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Because the records are not linearly increasing by ID, we need to load all
	// of the ids we've looked at. We store them in a map in memory for quick
	// access
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Records"))
		if err != nil {
			return err
		}

		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			id, _ := binary.Uvarint(k)
			loadedMatches[id] = true
		}

		fmt.Println("Loaded", len(loadedMatches), "matches from DB")

		return err
	})

	if err != nil {
		panic(err)
	}
}
