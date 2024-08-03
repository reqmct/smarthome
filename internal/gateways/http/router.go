package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"homework/internal/gateways/http/handlers"
	"homework/internal/gateways/http/models"
	"net/http"
	"strconv"
	"time"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	httpRequestErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Total number of HTTP request errors",
		},
		[]string{"method", "endpoint"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpResponseStatusCodes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_codes_total",
			Help: "Total number of HTTP response status codes",
		},
		[]string{"status_code", "method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestErrorsTotal, httpRequestDuration, httpResponseStatusCodes)
}

func metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		endpoint := c.FullPath()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
		httpRequestsTotal.WithLabelValues(method, endpoint).Inc()
		httpResponseStatusCodes.WithLabelValues(strconv.Itoa(statusCode), method, endpoint).Inc()

		if len(c.Errors) > 0 {
			httpRequestErrorsTotal.WithLabelValues(method, endpoint).Inc()
		}
	}
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func setupRouter(r *gin.Engine, cases UseCases, wsHandler *WebSocketHandler) {
	r.Use(metricsMiddleware())

	endpoints := []handlers.Handler{
		handlers.NewUsersHandler(cases.User),
		handlers.NewSensorsHandler(cases.Sensor),
		handlers.NewSensorHandler(cases.Sensor),
		handlers.NewEventsHandler(cases.Event),
		handlers.NewSensorOwnerHandler(cases.User),
		handlers.NewSensorHistoryHandler(cases.Event),
	}

	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, e := range endpoints {
		for _, method := range methods {
			if !contains(e.GetAvailableMethods(), method) {
				r.Handle(
					method,
					e.GetPath(),
					func(ctx *gin.Context) {
						ctx.AbortWithStatus(http.StatusMethodNotAllowed)
					},
				)
			}
		}
		e.SetupRouterGroup(r)
	}

	r.GET("/sensors/:sensor_id/events",
		func(ctx *gin.Context) {
			v := &models.SensorIDParam{}
			if err := ctx.ShouldBindUri(v); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "Error in the URI parameters of the request"})
				return
			}

			if err := v.Validate(nil); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{"reason": "URI parameters validation error: " + err.Error()})
				return
			}
			err := wsHandler.Handle(ctx, *v.SensorID)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"reason": err.Error()})
			}
		},
	)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
