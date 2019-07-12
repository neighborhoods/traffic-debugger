# Traffic Debugger
Traffic Debugger is an HTTP tool that assists with debugging of HTTP traffic flows.

## Client

### Usage

```
$ /client --help
Usage of /client:
  -abort-on-error
        Abort the client if an error is encountered
  -max-connections int
        Maximum number of active connections (default 1)
  -max-in-flight int
        Maximum number of active requests in flight (default 1)
  -max-requests-per-connection int
        Maximum number of requests each connection can have; 0 for unlimited
  -request-delay duration
        Time to wait between requests over the same connection
  -url string
        (required) URL to request
```

## Server

### Usage

```
$ /server --help
Usage of /server:
  -env-vars string
        Names of environment variables to return in responses
  -max-latency duration
        Max latency per request
  -max-requests int
        Maximum number of HTTP requests a single connection is allowed to handle before being closed
  -metrics-port int
        Optional port to listen for metrics and health check requests
  -min-latency duration
        Minimum latency per request
  -port int
        Port to listen on (default 8080)
```
