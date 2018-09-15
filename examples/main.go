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
