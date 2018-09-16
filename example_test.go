package middleware_test

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	prommiddleware "github.com/slok/go-prometheus-middleware"
)

// DefaultMiddleware shows how you would create a default middleware factory and wrap a
// handler to measure with the default settings. DEfault settings will act on prometheus
// default registry.
func ExampleMiddleware_defaultMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := prommiddleware.NewDefault()

	// Create our handler.
	myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world!"))
	})

	// Wrap our handler with the middleware.
	h := mdlw.Handler("", myHandler)

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go http.ListenAndServe(":801", promhttp.Handler())

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
