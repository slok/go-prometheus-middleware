# go-prometheus-middleware [![Build Status][travis-image]][travis-url] [![Go Report Card][goreport-image]][goreport-url] [![GoDoc][godoc-image]][godoc-url]

This middleware will measure metrics of a Go net/http handler in Prometheus format. The metrics measured are based on [RED] and/or [Four golden signals] and try to be measured in a efficent way.

## Getting Started

```golang
package main

import (
    "log"
    "net/http"

    "github.com/prometheus/client_golang/prometheus/promhttp"
    prommiddleware "github.com/slok/go-prometheus-middleware"
)

func main() {
    // Create our middleware.
    mdlw := prommiddleware.NewDefault()

    // Our handler.
    myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("hello world!"))
    })
    h := mdlw.Handler("", myHandler)

    // Serve metrics.
    log.Printf("serving metrics at: %s", ":9090")
    go http.ListenAndServe(":9090", promhttp.Handler())

    // Serve our handler.
    log.Printf("listening at: %s", ":8080")
    if err := http.ListenAndServe(":8080", h); err != nil {
        log.Panicf("error while serving: %s", err)
    }
}
```

For more examples check the the [examples]

## Metrics

The metrics obtained with this middleware are the [most important ones][red] for a HTTP service.

The middleware will measure the latency seconds of the requests using a histogram (latency), this will give us also the number of requests (rate), and the metric has the status codes (error rate):

### Query examples

Get the request rate by handler:

```text
sum(
    rate(http_request_duration_seconds_count[30s])
) by (handler)
```

Get the request error rate:

```text
rate(http_request_duration_seconds_count{code=~"5.."}[30s])
```

Get percentile 99 of the whole service:

```text
histogram_quantile(0.99,
    rate(http_request_duration_seconds_bucket[5m]))
```

Get percentile 90 of each handler:

```text
histogram_quantile(0.9,
    sum(
        rate(http_request_duration_seconds_bucket[10m])
    ) by (handler, le)
)
```

## Options

One of the options that you need to pass when wrapping the handler with the middleware is `handlerID`, this has 2 working ways.

- If an empty string is passed `mdwr.Handler("", h)` it will get the `handler` label from the url path. This will create very high cardnialty on the metrics because `/p/123/dashboard/1`, `/p/123/dashboard/2` and `/p/9821/dashboard/1` would have different `handler` labels. **This method is only recomended when the URLs are fixed (not dynamic or don't have parameters on the path)**.

- If a predefined handler ID is passed, `mdwr.Handler("/p/:userID/dashboard/:page", h)` this will keep cardinalty low because `/p/123/dashboard/1`, `/p/123/dashboard/2` and `/p/9821/dashboard/1` would have the same `handler` label on the metrics.

There are different parameters to set up your middleware factory, you can check everything on the [docs] and see the usage in the [examples].

## Benchmarks

```text
BenchmarkMiddlewareHandler/benchmark_with_URL-4                     1000000     1689 ns/op  320 B/op    6 allocs/op
BenchmarkMiddlewareHandler/benchmark_with_predefined_handler_ID-4   1000000     1849 ns/op  320 B/op    6 allocs/op
```

[travis-image]: https://travis-ci.org/slok/go-prometheus-middleware.svg?branch=master
[travis-url]: https://travis-ci.org/slok/go-prometheus-middleware
[goreport-image]: https://goreportcard.com/badge/github.com/slok/go-prometheus-middleware
[goreport-url]: https://goreportcard.com/report/github.com/slok/go-prometheus-middleware
[godoc-image]: https://godoc.org/github.com/slok/go-prometheus-middleware?status.svg
[godoc-url]: https://godoc.org/github.com/slok/go-prometheus-middleware
[docs]: https://godoc.org/github.com/slok/go-prometheus-middleware
[examples]: examples/
[red]: https://www.weave.works/blog/the-red-method-key-metrics-for-microservices-architecture/
[four golden signals]: https://landing.google.com/sre/book/chapters/monitoring-distributed-systems.html#xref_monitoring_golden-signals
