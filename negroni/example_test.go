package negroni_test

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	prommiddleware "github.com/slok/go-prometheus-middleware"
	promnegroni "github.com/slok/go-prometheus-middleware/negroni"
	"github.com/urfave/negroni"
)

// NegroniMiddleware shows how you would create a default middleware factory and use it
// to create a Negroni compatible middleware.
func Example_negroniMiddleware() {
	// Create our handler.
	myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world!"))
	})

	// Create our negroni instance.
	n := negroni.Classic()

	// Create our middleware factory with the default settings.
	mdlw := prommiddleware.NewDefault()
	// Add the middleware to negroni.
	n.Use(promnegroni.Handler("", mdlw))

	// Finally set our router on negroni.
	n.UseHandler(myHandler)

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go http.ListenAndServe(":8081", promhttp.Handler())

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := http.ListenAndServe(":8080", n); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
