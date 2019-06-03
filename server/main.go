package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var options Options
	parseFlags(&options)

	firstStats := generateStats([]string{})
	fmt.Printf("Running. UUID: %s\n", firstStats.UUID)

	connTracker := newConnectionTracker()

	var handler = handler{
		options:     options,
		connTracker: connTracker,
	}

	var mainRouterHasMetricsEndpoints = true
	if options.MetricsPort == 0 || options.MetricsPort == options.Port {
		mainRouterHasMetricsEndpoints = false
	}

	mainServer := buildMainRouter(&handler, mainRouterHasMetricsEndpoints)
	connTracker.Attach(mainServer)
	connTracker.SetupMetrics("http")

	var metricsServer *http.Server
	if options.MetricsPort > 0 && options.MetricsPort != options.Port {
		metricsServer = buildMetricsRouter()
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if metricsServer == nil {
			return
		}

		addr := generateAddr(options.MetricsPort)
		fmt.Printf("Serving metrics on port %v\n", addr)
		metricsServer.Addr = addr

		errCh := make(chan error)
		go func() {
			errCh <- metricsServer.ListenAndServe()
		}()

		select {
		case <-ctx.Done():
			metricsServer.Close()
		case err := <-errCh:
			panic(err)
		}
	}()

	addr := generateAddr(options.Port)
	fmt.Printf("Listening on %v\n", addr)
	mainServer.Addr = addr
	err := mainServer.ListenAndServe()
	cancel()
	panic(err)
}

func instrumentHandler(handler http.HandlerFunc) http.HandlerFunc {
	var requestCounterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "http",
		Subsystem: "requests",
		Name:      "total",
		Help:      "Count of HTTP requests.",
	}, []string{})

	requestCounterVec.With(nil).Add(0)
	handler = promhttp.InstrumentHandlerCounter(requestCounterVec, handler)

	var requestDurationVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Subsystem: "requests",
		Name:      "duration_seconds",
		Help:      "Duration of HTTP requests in seconds.",
	}, []string{})

	handler = promhttp.InstrumentHandlerDuration(requestDurationVec, handler)

	var requestInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "http",
		Subsystem: "requests",
		Name:      "in_flight",
		Help:      "Number of requests currently being processed.",
	})
	requestInFlight.Set(0)
	handler = promhttp.InstrumentHandlerInFlight(requestInFlight, handler).ServeHTTP

	return handler
}
