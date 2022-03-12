package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {

	start := time.Now()

	var waitGroup sync.WaitGroup

	for i := 0; i < 190; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			_, y, err := GetAllItemNames()

			if err != nil {
				log.Fatalf("Program failed: %v\n", err)
			}

			elapsed := time.Since(start)

			fmt.Printf("Time spent: %v, Web requests made: %v, Data sorted: %v\n", elapsed, y.TotalPages, y.TotalAuctions)
			fmt.Println("Clearing data")
		}()
	}

	waitGroup.Wait()

	elapsed := time.Since(start)
	log.Printf("Time taken to run: %s", elapsed)

}
