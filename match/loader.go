package match

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"net/http"
	"strconv"
)

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: true,
	}
	client = http.Client{Transport: tr}
)

func init() {
	gob.Register(Match{})
}

type MatchResult struct {
	MatchID uint64
	Match   Match
	Err     error
}

func LoadMatch(matchIDs <-chan uint64, loaded chan<- MatchResult) {
	var matchRes MatchResult
	for matchID := range matchIDs {
		fmt.Println("Loading", matchID)
		matchRes.MatchID = matchID
		req, err := http.NewRequest("GET", matchURL(matchID), nil)
		if err != nil {
			matchRes.Err = err
			loaded <- matchRes
			continue
		}
		req.Close = true

		matchResp, err := client.Do(req)
		if err != nil {
			matchRes.Err = err
			loaded <- matchRes
			continue
		}
		defer matchResp.Body.Close()

		dec := json.NewDecoder(matchResp.Body)
		var match Match
		err = dec.Decode(&match)
		matchRes.Match = match

		if err != nil {
			matchRes.Err = err
			loaded <- matchRes
			continue
		}
		loaded <- matchRes
	}
}

func LoadNewMatches() (int, error) {
	var err error
	db, err = bolt.Open("match.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	records, err := LoadNewRecords()
	if err != nil {
		return 0, err
	}
	newMatchCount := 0
	for _, record := range records {
		// If it's a new record
		if !loadedMatches[record.MatchID] {
			newMatchCount++
		}
	}
	newMatchIDs := make(chan uint64, newMatchCount)
	loaded := make(chan MatchResult, 50)
	for _, record := range records {
		// If it's a new record
		if !loadedMatches[record.MatchID] {
			newMatchIDs <- record.MatchID
		}
	}

	for i := 0; i < 20; i++ {
		go LoadMatch(newMatchIDs, loaded)
	}

	loadedCount := 0
	for {
		// If we've loaded everything
		if len(newMatchIDs) == 0 && loadedCount == newMatchCount {
			break
		}

		// Load from the worker pool
		matchRes := <-loaded
		if matchRes.Err != nil {
			fmt.Printf("Error loading %d: %q, retrying\n", matchRes.MatchID, matchRes.Err)
			newMatchIDs <- matchRes.MatchID
			continue
		}

		// Keep our in-memory representation up to date
		loadedMatches[matchRes.MatchID] = true
		if err := toDB(matchRes); err != nil {
			fmt.Printf("Error persisting %d to BoltDB: %q, skipping\n", matchRes.MatchID, matchRes.Err)
		}

		loadedCount++
	}

	close(loaded)

	return loadedCount, nil
}

func toDB(matchRes MatchResult) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Records"))

		// We make a []byte representation of our matchID
		buf := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(buf, matchRes.MatchID)
		// Now we encode our Match as a []byte, in compact, gobby format
		var serialMatch bytes.Buffer
		enc := gob.NewEncoder(&serialMatch)
		err := enc.Encode(matchRes.Match)
		if err != nil {
			return err
		}

		err = b.Put(buf, serialMatch.Bytes())
		return err
	})
}

func LoadNewRecords() ([]MatchRecord, error) {
	histResp, err := http.Get("https://robotgame.net/api/match/history")
	if err != nil {
		return nil, err
	}
	defer histResp.Body.Close()

	dec := json.NewDecoder(histResp.Body)
	var records []MatchRecord
	err = dec.Decode(&records)

	if err != nil {
		return nil, err
	}
	return records, nil
}

func matchURL(id uint64) string {
	return "https://robotgame.net/api/match/" + strconv.FormatInt(int64(id), 10)
}
