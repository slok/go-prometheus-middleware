package middleware_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"

	prommiddleware "github.com/slok/go-prometheus-middleware"
)

func getFakeHandler(statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	})
}

func TestMiddlewareHandler(t *testing.T) {
	tests := []struct {
		name       string
		config     prommiddleware.Config
		requests   func(h http.Handler)
		handlerID  string
		statusCode int
		expMetrics []string
	}{
		{
			name:       "default configuration without handlerID  should measure without prefix with the default buckets and the URL as handler",
			config:     prommiddleware.Config{},
			handlerID:  "",
			statusCode: 403,
			requests: func(h http.Handler) {
				r := httptest.NewRequest("GET", "/test", nil)
				r2 := httptest.NewRequest("POST", "/test2", nil)
				h.ServeHTTP(httptest.NewRecorder(), r)
				h.ServeHTTP(httptest.NewRecorder(), r2)
			},
			expMetrics: []string{
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="0.005"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="0.01"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="0.025"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="0.05"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="0.1"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="0.25"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="0.5"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="1"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="5"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="10"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test",method="GET",le="+Inf"} 1`,
				`http_request_duration_seconds_count{code="403",handler="/test",method="GET"} 1`,

				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="0.005"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="0.01"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="0.025"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="0.05"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="0.1"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="0.25"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="0.5"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="1"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="2.5"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="5"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="10"} 1`,
				`http_request_duration_seconds_bucket{code="403",handler="/test2",method="POST",le="+Inf"} 1`,
				`http_request_duration_seconds_count{code="403",handler="/test2",method="POST"} 1`,
			},
		},
		{
			name: "custom configuration with handlerID should measure with prefix, with the default buckets and the the same handlerID",
			config: prommiddleware.Config{
				Prefix:  "batman",
				Buckets: []float64{.5, 1, 2.5, 5, 10, 20, 40, 80, 160, 320},
			},
			handlerID:  "bruceWayne",
			statusCode: 201,
			requests: func(h http.Handler) {
				r := httptest.NewRequest("GET", "/test", nil)
				r2 := httptest.NewRequest("POST", "/test2", nil)
				h.ServeHTTP(httptest.NewRecorder(), r)
				h.ServeHTTP(httptest.NewRecorder(), r2)
			},
			expMetrics: []string{
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="0.5"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="1"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="2.5"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="5"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="10"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="20"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="40"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="80"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="160"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="320"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="GET",le="+Inf"} 1`,

				`batman_http_request_duration_seconds_count{code="201",handler="bruceWayne",method="GET"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="0.5"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="1"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="2.5"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="5"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="10"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="20"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="40"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="80"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="160"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="320"} 1`,
				`batman_http_request_duration_seconds_bucket{code="201",handler="bruceWayne",method="POST",le="+Inf"} 1`,
				`batman_http_request_duration_seconds_count{code="201",handler="bruceWayne",method="POST"} 1`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			reg := prometheus.NewRegistry()
			m := prommiddleware.New(test.config, reg)
			h := m.Handler(test.handlerID, getFakeHandler(test.statusCode))

			// Make the calls to our handler.
			test.requests(h)

			// Get the metrics handler and serve.
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/metrics", nil)
			promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(rec, req)

			resp := rec.Result()

			// Check all metrics are present.
			if assert.Equal(http.StatusOK, resp.StatusCode) {
				body, _ := ioutil.ReadAll(resp.Body)
				for _, expMetric := range test.expMetrics {
					assert.Contains(string(body), expMetric, "metric not present on the result")
				}
			}
		})
	}
}

func BenchmarkMiddlewareHandler(b *testing.B) {
	b.StopTimer()

	benchs := []struct {
		name      string
		handlerID string
	}{
		{
			name:      "benchmark with URL",
			handlerID: "",
		},
		{
			name:      "benchmark with predefined handler ID",
			handlerID: "benchmark1",
		},
	}

	for _, bench := range benchs {
		b.Run(bench.name, func(b *testing.B) {
			// Prepare.
			reg := prometheus.NewRegistry()
			m := prommiddleware.New(prommiddleware.Config{}, reg)
			h := m.Handler(bench.handlerID, getFakeHandler(200))
			r := httptest.NewRequest("GET", "/test", nil)

			// Make the requests.
			for n := 0; n < b.N; n++ {
				b.StartTimer()
				h.ServeHTTP(httptest.NewRecorder(), r)
				b.StopTimer()
			}
		})
	}
}
