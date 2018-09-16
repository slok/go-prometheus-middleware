package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Config is the configuration for the middleware factory.
type Config struct {
	// Prefix is the prefix that will be set on the metrics, by default it will be empty.
	Prefix string
	// Buckets are the buckets used by Prometheus for the HTTP request metrics, by default
	// Uses Prometheus default buckets (from 5ms to 10s).
	Buckets []float64
	// GroupedStatus will group the status label in the form of `\dxx`, for example,
	// 200, 201, and 203 will have the label `code="2xx"`. This impacts on the cardinality
	// of the metrics and also improves the performance of queries that are grouped by
	// status code because there are already aggregated in the metric.
	// By default will be false.
	GroupedStatus bool
}

func (c *Config) validate() {
	if len(c.Buckets) == 0 {
		c.Buckets = prometheus.DefBuckets
	}
}

// Middleware is a factory that creates middlewares or wrappers that
// measure requests to the wrapped handler using Prometheus metrics.
type Middleware interface {
	// Handler wraps the received handler with the Prometheus middleware.
	// The first argument receives the handlerID, all the metrics will have
	// that handler ID as the handler label on the metrics, if an empty
	// string is passed then it will get the handlerID from the request
	// path.
	Handler(handlerID string, h http.Handler) http.Handler
}

// middelware is the prometheus middleware instance.
type middleware struct {
	httpRequestHistogram *prometheus.HistogramVec

	cfg Config
	reg prometheus.Registerer
}

// NewDefault returns the default Prometheus middleware factory
// that will wrap the handlers using the default middleware values.
func NewDefault() Middleware {
	return New(Config{}, nil)
}

// New returns the a Prometheus middleware factory that will
// that will wrap the handlers using the customized middleware values.
func New(cfg Config, reg prometheus.Registerer) Middleware {
	// If no registerer then set the default one.
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}

	// Validate the configuration.
	cfg.validate()

	// Create our middleware with all the configuration options.
	m := &middleware{
		httpRequestHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Prefix,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "The latency of the HTTP requests.",
			Buckets:   cfg.Buckets,
		}, []string{"handler", "method", "code"}),

		cfg: cfg,
		reg: reg,
	}

	// Register all the middleware metrics on prometheus registerer.
	m.registerMetrics()

	return m
}

func (m *middleware) registerMetrics() {
	m.reg.MustRegister(
		m.httpRequestHistogram,
	)
}

// Handler satisfies Middlware interface.
func (m *middleware) Handler(handlerID string, h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intercept the writer so we can retrieve data afterwards.
		wi := &responseWriterInterceptor{
			statusCode:     http.StatusOK,
			ResponseWriter: w,
		}

		// If there isn't predefined handler ID we
		// set that ID as the URL path.
		hid := handlerID
		if handlerID == "" {
			hid = r.URL.Path
		}

		// Start the timer and when finishing measure the duration.
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()

			// If we need to group the status code, it uses the
			// first number of the status code because is the least
			// required identification way.
			var code string
			if m.cfg.GroupedStatus {
				code = fmt.Sprintf("%dxx", wi.statusCode/100)
			} else {
				code = strconv.Itoa(wi.statusCode)
			}

			m.httpRequestHistogram.WithLabelValues(hid, r.Method, code).Observe(duration)
		}()

		h.ServeHTTP(wi, r)
	})
}

// responseWriterInterceptor is a simple wrapper to incercept set data on a
// ResponseWriter.
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
