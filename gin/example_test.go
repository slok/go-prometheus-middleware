package gin_test

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	prommiddleware "github.com/slok/go-prometheus-middleware"
	promgin "github.com/slok/go-prometheus-middleware/gin"
)

// GinMiddleware shows how you would create a default middleware factory and use it
// to create a Gin compatible middleware.
func Example_ginMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := prommiddleware.NewDefault()

	// Create our gin instance.
	r := gin.New()

	// Add the middlewares to all gin routes.
	r.Use(promgin.Handler("", mdlw))

	// Add our handler
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go http.ListenAndServe(":8081", promhttp.Handler())

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := r.Run(":8080"); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
