package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type connectionTracker struct {
	connections map[string]connectionStats
	mux         sync.RWMutex

	metricsConnectionsOpened    prometheus.Counter
	metricsConnectionsClosed    prometheus.Counter
	metricsConnectionsRequests  prometheus.Histogram
	metricsConnectionsDurations prometheus.Histogram
}

type connectionStats struct {
	requestCount int
	openTime     time.Time
}

func newConnectionTracker() *connectionTracker {
	return &connectionTracker{
		connections: make(map[string]connectionStats),
	}
}

func (t *connectionTracker) Attach(server *http.Server) {
	server.ConnState = t.handleConnStateEvent
}

func (t *connectionTracker) SetupMetrics(namespace string) {
	if namespace == "" {
		namespace = "http"
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.metricsConnectionsOpened = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "opened",
		Help:      "Number of connections that were opened.",
	})
	t.metricsConnectionsOpened.Add(0)

	t.metricsConnectionsClosed = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "closed",
		Help:      "Number of connections that were closed.",
	})
	t.metricsConnectionsClosed.Add(0)

	t.metricsConnectionsRequests = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "requests",
		Help:      "Number of requests made per connection.",

		Buckets: []float64{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000},
	})

	t.metricsConnectionsDurations = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "duration_seconds",
		Help:      "Length of time the connections were opened.",

		Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 15, 30, 60, 120, 300},
	})
}

func (t *connectionTracker) GetRequestCount(req *http.Request) int {
	id := getRequestConnectionID(req)

	t.mux.RLock()
	defer t.mux.RUnlock()

	stats, ok := t.connections[id]
	if !ok {
		panic(fmt.Errorf("cannot get connection stats for untracked connection %v", id))
	}

	return stats.requestCount
}

func (t *connectionTracker) handleConnStateEvent(conn net.Conn, event http.ConnState) {
	if conn == nil {
		return
	}

	switch event {
	case http.StateNew:
		t.newConnection(conn)
	case http.StateActive:
		t.newRequest(conn)
	case http.StateHijacked, http.StateClosed:
		t.closedConnection(conn)
	}
}

func getConnectionID(conn net.Conn) string {
	return conn.RemoteAddr().String()
}

func getRequestConnectionID(req *http.Request) string {
	return req.RemoteAddr
}

func (t *connectionTracker) newConnection(conn net.Conn) {
	id := getConnectionID(conn)

	t.mux.Lock()
	defer t.mux.Unlock()

	_, ok := t.connections[id]
	if ok {
		panic(fmt.Errorf("unexpected already has connection for %v", id))
	}

	stats := connectionStats{
		openTime: time.Now(),
	}

	t.connections[id] = stats

	if t.metricsConnectionsOpened != nil {
		t.metricsConnectionsOpened.Inc()
	}
}

func (t *connectionTracker) newRequest(conn net.Conn) {
	id := getConnectionID(conn)

	t.mux.Lock()
	defer t.mux.Unlock()

	stats, ok := t.connections[id]
	if !ok {
		panic(fmt.Errorf("missing connection stats for %v", id))
	}

	stats.requestCount++
	t.connections[id] = stats
}

func (t *connectionTracker) closedConnection(conn net.Conn) {
	id := getConnectionID(conn)

	t.mux.Lock()
	defer t.mux.Unlock()

	stats, ok := t.connections[id]
	if !ok {
		panic(fmt.Errorf("missing connection stats for %v", id))
	}

	delete(t.connections, id)

	if t.metricsConnectionsClosed != nil {
		t.metricsConnectionsClosed.Inc()
	}

	if t.metricsConnectionsDurations != nil {
		connDuration := time.Since(stats.openTime)
		t.metricsConnectionsDurations.Observe(connDuration.Seconds())
	}

	if t.metricsConnectionsRequests != nil {
		t.metricsConnectionsRequests.Observe(float64(stats.requestCount))
	}
}
