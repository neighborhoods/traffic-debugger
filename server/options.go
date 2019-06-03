package main

import (
	"flag"
	"strings"
	"time"
)

// Options is a struct of command line options
type Options struct {
	Port        int
	MetricsPort int
	EnvVars     []string
	MinLatency  time.Duration
	MaxLatency  time.Duration
	MaxRequests int
}

func parseFlags(options *Options) {
	flag.IntVar(&options.Port, "port", 8080, "Port to listen on")
	flag.IntVar(&options.MetricsPort, "metrics-port", 0, "Optional port to listen for metrics and health check requests")
	flag.DurationVar(&options.MinLatency, "min-latency", 0, "Minimum latency per request")
	flag.DurationVar(&options.MaxLatency, "max-latency", 0, "Max latency per request")
	flag.IntVar(&options.MaxRequests, "max-requests", 0, "Maximum number of HTTP requests a single connection is allowed to handle before being closed")

	var envVars string
	flag.StringVar(&envVars, "env-vars", "", "Names of environment variables to return in responses")

	flag.Parse()
	options.EnvVars = parseEnvVarOption(envVars)
}

func parseEnvVarOption(raw string) (envVars []string) {
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return
	}

	envVars = strings.Split(raw, ",")

	for i, name := range envVars {
		envVars[i] = strings.TrimSpace(name)
	}

	return
}
