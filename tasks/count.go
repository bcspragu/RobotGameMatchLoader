package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("match.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var count int
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Records"))
		count = b.Stats().KeyN
		return nil
	})

	fmt.Println(count, "match records")
}
