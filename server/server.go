package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func buildMainRouter(handler *handler, registerMetrics bool) *http.Server {
	mux := http.NewServeMux()

	if registerMetrics {
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/healthz", okResponse)
	}

	mux.HandleFunc("/favicon.ico", http.NotFound)
	statHandler := http.HandlerFunc(handler.handleRoot)
	statHandler = instrumentHandler(statHandler)
	mux.Handle("/", statHandler)

	server := &http.Server{}
	server.Handler = mux

	return server
}

func buildMetricsRouter() *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", okResponse)
	mux.HandleFunc("/", http.NotFound)

	server := &http.Server{}
	server.Handler = mux

	return server
}

func generateAddr(port int) string {
	return fmt.Sprintf(":%d", port)
}
