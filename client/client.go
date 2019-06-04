package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type loopOptions struct {
	requests                 <-chan *http.Request
	completions              chan<- error
	maxRequestsPerConnection int
	requestDelay             time.Duration
}

func requestLoop(ctx context.Context, options *loopOptions) {
	var dialer net.Dialer
	transport := &http.Transport{
		MaxConnsPerHost: 1,
		IdleConnTimeout: time.Duration(10 * time.Second),
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			fmt.Printf("Opening new connection\n")
			return dialer.DialContext(ctx, network, address)
		},
	}

	defer transport.CloseIdleConnections()

	var requestCount int

	clientID := mustRandomHex(8)

	for {
		select {
		case <-ctx.Done():
			return
		case req, ok := <-options.requests:
			if !ok {
				return
			}

			requestCount++
			var shouldClose bool
			if options.maxRequestsPerConnection > 0 && requestCount >= options.maxRequestsPerConnection {
				requestCount = 0
				shouldClose = true
			}

			if shouldClose {
				req.Close = true
			}

			req.Header.Set("X-Client-ID", clientID)

			err := sendRequest(ctx, transport, req)

			if shouldClose {
				transport.CloseIdleConnections()
			}

			options.completions <- err
			time.Sleep(options.requestDelay)
		}
	}
}

func sendRequest(ctx context.Context, tripper http.RoundTripper, req *http.Request) error {
	var err error

	req = req.WithContext(ctx)
	fmt.Printf("%v %v\n", req.Method, req.URL.String())
	response, err := tripper.RoundTrip(req)
	fmt.Printf("Response: %v\n", response.StatusCode)
	if err != nil {
		return err
	}

	if response == nil {
		panic("unexpected nil response after nil error")
	}

	err = ReadAndClose(response.Body)
	if err != nil {
		return err
	}

	statusClass := response.StatusCode / 100
	if statusClass != 2 {
		err = fmt.Errorf("got unexpected status code: %d", response.StatusCode)
		return err
	}

	return nil
}
