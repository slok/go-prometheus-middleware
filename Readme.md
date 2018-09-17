# go-prometheus-middleware [![Build Status][travis-image]][travis-url] [![Go Report Card][goreport-image]][goreport-url] [![GoDoc][godoc-image]][godoc-url]

This middleware will measure metrics of a Go net/http handler in Prometheus format. The metrics measured are based on [RED] and/or [Four golden signals], follow standards and try to be measured in a efficent way.

If you are using a framework that isn't directly compatible with go's `http.Handler` interface from the std library, do not worry, there are multiple helpers available to get middlewares fo the most used http Go frameworks.

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

For more examples check the the [examples]. [default][default-example] and [custom][custom-example] are the examples for Go net/http std library users.

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

### Factory options

The factory options are the ones that are passed in the moment of creating the middleware factory using the `Config` object.

#### Prefix

This option will make exposed metrics have a `{PREFIX}_` in fornt of the metric. For example if a regular exposed metric is `http_request_duration_seconds_count` and I use `Prefix: batman` my exposed metric will be `batman_http_request_duration_seconds_count`. By default this will be disabled or empty, but can be useful if all the metrics of the app are prefixed with the app name.

#### Buckets

Buckets are the buckets used for the histogram metric, by default it will use Prometheus defaults, this is from 5ms to 10s, on a regular HTTP service this is very common and in most cases this default works perfect, but on some cases where the latency is very low or very high due the nature of the service, this could be changed to measure a different range of time. Example, from 500ms to 320s `Buckets: []float64{.5, 1, 2.5, 5, 10, 20, 40, 80, 160, 320}`. Is not adviced to use more than 10 buckets.

#### GroupedStatus

Storing all the status codes could increase the cardinality of the metrics, usually this is not a common case because the used status codes by a service are not too much and are finite, but some services use a lot of different status codes, grouping the status on the `\dxx` form could impact the performance (in a good way) of the queries on Prometheus (as they are already aggregated), on the other hand it losses detail. For example the metrics code `code="401"`, `code="404"`, `code="403"` with this enabled option would end being `code="4xx"` label. By default is disabled.

### Wrapper options

The wrapper options are the ones passed in the moment of creating the wrapper middleware using the factory `Middleware`.

#### handlerID

One of the options that you need to pass when wrapping the handler with the middleware is `handlerID`, this has 2 working ways.

- If an empty string is passed `mdwr.Handler("", h)` it will get the `handler` label from the url path. This will create very high cardnialty on the metrics because `/p/123/dashboard/1`, `/p/123/dashboard/2` and `/p/9821/dashboard/1` would have different `handler` labels. **This method is only recomended when the URLs are fixed (not dynamic or don't have parameters on the path)**.

- If a predefined handler ID is passed, `mdwr.Handler("/p/:userID/dashboard/:page", h)` this will keep cardinalty low because `/p/123/dashboard/1`, `/p/123/dashboard/2` and `/p/9821/dashboard/1` would have the same `handler` label on the metrics.

There are different parameters to set up your middleware factory, you can check everything on the [docs] and see the usage in the [examples].

## Frameworks

The middleware is mainly focused to be compatible with Go std library using http.Handler, but it comes with helpers to get middlewares for other frameworks or libraries.

**The different helpers are on separate packages so when you import the project it doesn't import other framework packages and dependencies, for example if I don't use Negroni and instead I use std go net/http, it wouldn't be nice to import Negroni on my project.**

- [Negroni][negroni-example]
- [Gin][gin-example]

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
[default-example]: examples/default
[custom-example]: examples/custom
[negroni-example]: examples/negroni
[gin-example]: examples/gin
