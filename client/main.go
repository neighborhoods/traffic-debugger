package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func requestGenerator(ctx context.Context, requestCh chan<- *http.Request, doneCh <-chan struct{}, maxInFlight int, url url.URL) {
	defer close(requestCh)

	requestsAvailable := maxInFlight
	for {
		outCh := requestCh
		if requestsAvailable <= 0 {
			// nil channels block forever
			outCh = nil
		}

		req := http.Request{
			Method: "GET",
			URL:    &url,
			Header: make(http.Header),
		}

		select {
		case <-ctx.Done():
			return
		case _, ok := <-doneCh:
			if !ok {
				return
			}

			requestsAvailable++
		case outCh <- &req:
			requestsAvailable--
		}
	}
}

func main() {
	var options Options
	parseFlags(&options)

	requestCh := make(chan *http.Request)
	resultCh := make(chan error)

	ctx := context.Background()

	loopOptions := &loopOptions{
		requests:                 requestCh,
		completions:              resultCh,
		maxRequestsPerConnection: options.MaxRequestsPerConnection,
		requestDelay:             options.RequestDelay,
	}

	for i := 0; i < options.MaxConnections; i++ {
		go requestLoop(ctx, loopOptions)
	}

	doneCh := make(chan struct{}, options.MaxInFlight)
	defer close(doneCh)
	go requestGenerator(ctx, requestCh, doneCh, options.MaxInFlight, options.URL)

	for {
		err, ok := <-resultCh
		if !ok {
			break
		}

		doneCh <- struct{}{}

		if err != nil {
			fmt.Printf("Got error: %v\n", err)
			if options.AbortOnError {
				break
			}
		}
	}
}
