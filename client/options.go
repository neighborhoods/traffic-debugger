package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"
)

// Options is a struct of command line options
type Options struct {
	AbortOnError             bool
	MaxConnections           int
	MaxInFlight              int
	MaxRequestsPerConnection int
	URL                      url.URL
	RequestDelay             time.Duration
}

func parseFlags(options *Options) {
	flag.BoolVar(&options.AbortOnError, "abort-on-error", false, "Abort the client if an error is encountered")
	flag.IntVar(&options.MaxConnections, "max-connections", 1, "Maximum number of active connections")
	flag.IntVar(&options.MaxInFlight, "max-in-flight", 1, "Maximum number of active requests in flight")
	flag.IntVar(&options.MaxRequestsPerConnection, "max-requests-per-connection", 0, "Maximum number of requests each connection can have; 0 for unlimited")
	flag.DurationVar(&options.RequestDelay, "request-delay", 0, "Time to wait between requests over the same connection")

	var rawURL string
	flag.StringVar(&rawURL, "url", "", "(required) URL to request")

	flag.Parse()

	if rawURL == "" {
		printUsage()
		panic("missing -url")
	}

	url, err := url.Parse(rawURL)
	if err != nil {
		printUsage()
		panic(err)
	}

	options.URL = *url
}

func printUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
