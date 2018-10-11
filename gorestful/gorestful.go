// Package gorestful is a helper package to get a gorestful compatible
// handler/middleware from the standard net/http Middleware factory
// (from github.com/slok/go-prometheus-middleware).
package gorestful

import (
	"net/http"

	gorestful "github.com/emicklei/go-restful"

	prommiddleware "github.com/slok/go-prometheus-middleware"
)

// Handler returns a gorestful compatible middleware from a Middleware factory instance.
// The first handlerID argument is the same argument passed on Middleware.Handler method.
func Handler(handlerID string, m prommiddleware.Middleware) gorestful.FilterFunction {
	// Create a dummy handler to wrap the middleware chain of gorestful, this way Middleware
	// interface can wrap the gorestful chain.
	return func(req *gorestful.Request, resp *gorestful.Response, chain *gorestful.FilterChain) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req.Request = r
			resp.ResponseWriter = w
			chain.ProcessFilter(req, resp)
		})

		m.Handler(handlerID, h).ServeHTTP(resp.ResponseWriter, req.Request)
	}
}
