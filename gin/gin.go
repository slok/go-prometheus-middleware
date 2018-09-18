// Package gin is a helper package to get a gin compatible
// handler/middleware from the standard net/http Middleware factory
// (from github.com/slok/go-prometheus-middleware).
package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	prommiddleware "github.com/slok/go-prometheus-middleware"
)

// Handler returns a gin compatible middleware from a Middleware factory instance.
// The first handlerID argument is the same argument passed on Middleware.Handler method.
func Handler(handlerID string, m prommiddleware.Middleware) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// Create a dummy handler to wrap the middleware chain of gin, this way Middleware
		// interface can wrap the gin chain.
		dh := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			ctx.Next()
		})

		m.Handler(handlerID, dh).ServeHTTP(ctx.Writer, ctx.Request)
	})
}
