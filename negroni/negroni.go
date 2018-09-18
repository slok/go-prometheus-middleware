// Package negroni is a helper package to get a negroni compatible
// handler/middleware from the standatd net/http Middleware factory
// (from github.com/slok/go-prometheus-middleware).
package negroni

import (
	"net/http"

	"github.com/urfave/negroni"

	prommiddleware "github.com/slok/go-prometheus-middleware"
)

// Handler returns a Negroni compatible middleware from a Middleware factory instance.
// The first HandlerID argument is the same argument passed on Middleware.Handler method.
func Handler(handlerID string, m prommiddleware.Middleware) negroni.Handler {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		m.Handler(handlerID, next).ServeHTTP(rw, r)
	})
}
