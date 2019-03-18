# go-prometheus-middleware

This library has been deprecated in favor of [go-http-metrics], this new library has the same features as this one, plus the extensibility to new metrics formats (not only Prometheus).

The migration is simple, example:

From this:

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

To this:

```golang
package main

import (
    "log"
    "net/http"

    "github.com/prometheus/client_golang/prometheus/promhttp"
    metrics "github.com/slok/go-http-metrics/metrics/prometheus"
    "github.com/slok/go-http-metrics/middleware"
)

func main() {
    // Create our middleware.
    mdlw := middleware.New(middleware.Config{
        Recorder: metrics.New(metrics.Config{}),
    })

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

[go-http-metrics]: https://github.com/slok/go-http-metrics
