package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/storage"
)

func main() {
	appCtx := context.Background()
	storageClient, err := storage.NewClient(appCtx)
	if err != nil {
		panic(err)
	}

	counter := 0

	for {
		loopCtx, cancel := context.WithCancel(appCtx)
		myCounter := counter
		counter++
		fmt.Printf("[%d] Spawning read\n", myCounter)
		go func() {
			shouldClose := rand.Int31n(2)
			if shouldClose > 0 {
				sleepyTime := rand.Int31n(2000)
				go func() { time.Sleep(time.Duration(sleepyTime) * time.Microsecond); cancel() }()
			}
			w := storageClient.Bucket("kjs_cool_vault_bucket").Object(fmt.Sprintf("test1/%d", myCounter)).NewWriter(loopCtx)

			defer func() {
				closeErr := w.Close()
				if closeErr != nil {
					fmt.Printf("[%d] We got a close error: %s\n", myCounter, closeErr)
				} else {
					fmt.Printf("[%d] Closed OK\n", myCounter)
				}
			}()
			if _, err := w.Write([]byte("hello i am the data")); err != nil {
				fmt.Printf("[%d] Write failed: %s\n", myCounter, err)
			}
		}()
		time.Sleep(10000 * time.Microsecond)
	}
}
