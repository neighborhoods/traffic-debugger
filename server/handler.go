package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var contentTypeHeader = http.CanonicalHeaderKey("Content-Type")
var connectionHeader = http.CanonicalHeaderKey("Connection")

type handler struct {
	options     Options
	connTracker *connectionTracker
}

func (h handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Printf("got request to %v\n", r.URL.Path)
		return
	}

	if h.options.MaxRequests > 0 {
		requestCount := h.connTracker.GetRequestCount(r)
		if requestCount >= h.options.MaxRequests {
			w.Header()[connectionHeader] = []string{"close"}
			fmt.Printf("requested connection close (%v requests) for %v\n", requestCount, getRequestConnectionID(r))
		}
	}

	stats := generateStats(h.options.EnvVars)
	var sleptFor = randomSleep(h.options.MinLatency, h.options.MaxLatency)
	stats.SleepDuration = sleptFor
	stats.SleepDurationString = sleptFor.String()

	data, err := json.MarshalIndent(stats, "", "\t")
	if err != nil {
		w.Header()[contentTypeHeader] = []string{"text/plain"}
		r.Close = true
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	// fmt.Printf("Handled url: %s\n", r.URL.Path)
	// fmt.Printf("Handled request: %s\n", stats.RandomString)

	w.Header()[contentTypeHeader] = []string{"application/json"}
	w.Write(data)
}

func okResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
