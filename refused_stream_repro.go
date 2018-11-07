package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/net/http2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func main() {
	appCtx := context.Background()
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: false,
	}
	tlsConfig := &tls.Config{}
	httpTransport := &http.Transport{
		TLSClientConfig:       tlsConfig,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if err := http2.ConfigureTransport(httpTransport); err != nil {
		panic(fmt.Sprintf("failed to configure http2: %s", err))
	}

	tokenSource, err := google.DefaultTokenSource(appCtx)

	oauth2Transport := &oauth2.Transport{
		Source: tokenSource,
		Base:   httpTransport,
	}
	httpClient := &http.Client{
		Transport: oauth2Transport,
	}

	storageClient, err := storage.NewClient(appCtx, option.WithHTTPClient(httpClient))
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
