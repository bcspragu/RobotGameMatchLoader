package main

import (
	"../match"
	"bufio"
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"os"
)

func main() {
	gob.Register(match.Match{})

	db, err := bolt.Open("match.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	f, err := os.Create("rounds.csv")
	writer := bufio.NewWriter(f)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Records"))
		b.ForEach(func(k, v []byte) error {
			var match match.Match
			dec := gob.NewDecoder(bytes.NewReader(v))
			err := dec.Decode(&match)
			if err != nil {
				return err
			}
			if match.Winner == match.R1ID {
				writer.WriteString("-1 ")
			} else {
				writer.WriteString("1 ")
			}
			writer.WriteString(match.GridString())
			return nil
		})
		return nil
	})

	writer.Flush()
}
