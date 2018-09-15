# go-prometheus-middleware

This middleware will measure the request latency of the different handlers.

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
	go http.ListenAndServe(":9090", promhttp.Handler())

	// Serve our handler.
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
```

For more examples check the the [examples]

## Options

One of the options that you need to pass when wrapping the handler with the middleware is `handlerID`, this has 2 working ways.

- If you pass a empty string `mdwr.Handler("", h)` it will get the `handler` label from the url path. This will create very high cardnialty on the metrics because `/p/123/dashboard/1`, `/p/123/dashboard/2` and `/p/9821/dashboard/1` would have different `handler` labels. **This method is only recomended when the URLs are fixed (not dynamic or don't have parameters on the path)**.

- If pass a predefined handler ID, `mdwr.Handler("/p/:userID/dashboard/:page", h)` and this will keep cardinalty low because `/p/123/dashboard/1`, `/p/123/dashboard/2` and `/p/9821/dashboard/1` would have the same `handler` label on the metrics.

There are different parameters to set up your middleware factory, you can check everything on the [docs] and see the usage in the [examples].

## Benchmarks

[docs]: http://godoc.org/github.com/slok/go-prometheus-middleware
[examples]: examples/
