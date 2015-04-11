package main

import (
	"./match"
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Loading new matches")
		c, err := match.LoadNewMatches()
		if err != nil {
			fmt.Println("Error loading matches:", err)
			break
		}
		fmt.Println("Loaded", c, "new matches")
		// Wait a day before loading matches again
		<-time.Tick(24 * time.Hour)
	}
}
