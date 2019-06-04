package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
	connectionHeader  = "Connection"
	clientIDHeader    = "X-Client-ID"
)

type handler struct {
	options     Options
	connTracker *connectionTracker
}

func (h handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Printf("Handled unknown path: %v\n", r.URL.Path)
		return
	}

	clientID := r.Header.Get(clientIDHeader)
	if clientID != "" {
		h.connTracker.TrackRequestClient(r, clientID)
	}

	if h.options.MaxRequests > 0 {
		requestCount := h.connTracker.GetRequestCount(r)
		if requestCount >= h.options.MaxRequests {
			w.Header().Set(connectionHeader, "close")
		}
	}

	stats := generateStats(h.options.EnvVars)
	var sleptFor = randomSleep(h.options.MinLatency, h.options.MaxLatency)
	stats.SleepDuration = sleptFor
	stats.SleepDurationString = sleptFor.String()

	data, err := json.MarshalIndent(stats, "", "\t")
	if err != nil {
		w.Header().Set(contentTypeHeader, "text/plain")
		r.Close = true
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	w.Header().Set(contentTypeHeader, "application/json")
	w.Write(data)
}

func okResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
