package main

import (
	"./match"
	"bufio"
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"os"
	"strconv"
)

func main() {
	gob.Register(match.Match{})

	db, err := bolt.Open("match.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	writers := make([]*bufio.Writer, 50)
	for i := 0; i < 50; i++ {
		f, err := os.Create("rounds/out" + strconv.Itoa(i) + ".csv")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		writers[i] = bufio.NewWriter(f)
	}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Records"))
		b.ForEach(func(k, v []byte) error {
			var match match.Match
			dec := gob.NewDecoder(bytes.NewReader(v))
			err := dec.Decode(&match)
			if err != nil {
				return err
			}
			for i := 0; i < 50; i++ {
				writers[i].WriteString(match.GridString(i))
			}
			return nil
		})
		return nil
	})

	for i := 0; i < 50; i++ {
		writers[i].Flush()
	}
}
