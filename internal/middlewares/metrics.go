package middlewares

import (
	"strconv"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type FiberPrometheus struct {
	operationsProcessed *prometheus.CounterVec
	operationsDuration  *prometheus.HistogramVec
	defaultURL          string
}

func create(registry prometheus.Registerer, serviceName, namespace, subsystem string, labels map[string]string) *FiberPrometheus {
	constLabels := make(prometheus.Labels)
	if serviceName != "" {
		constLabels["service"] = serviceName
	}
	for label, value := range labels {
		constLabels[label] = value
	}

	counter := promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName(namespace, subsystem, "requests_total"),
			Help:        "Count all http requests by status code, method and path.",
			ConstLabels: constLabels,
		},
		[]string{"status_code", "method", "path"},
	)
	histogram := promauto.With(registry).NewHistogramVec(prometheus.HistogramOpts{
		Name:        prometheus.BuildFQName(namespace, subsystem, "request_duration_seconds"),
		Help:        "Duration of all HTTP requests by status code, method and path.",
		ConstLabels: constLabels,
		Buckets: []float64{
			0.000000001, // 1ns
			0.000000002,
			0.000000005,
			0.00000001, // 10ns
			0.00000002,
			0.00000005,
			0.0000001, // 100ns
			0.0000002,
			0.0000005,
			0.000001, // 1µs
			0.000002,
			0.000005,
			0.00001, // 10µs
			0.00002,
			0.00005,
			0.0001, // 100µs
			0.0002,
			0.0005,
			0.001, // 1ms
			0.002,
			0.005,
			0.01, // 10ms
			0.02,
			0.05,
			0.1, // 100 ms
			0.2,
			0.5,
			1.0, // 1s
			2.0,
			5.0,
			10.0, // 10s
			15.0,
			20.0,
			30.0,
		},
	},
		[]string{"status_code", "method", "path"},
	)

	return &FiberPrometheus{
		operationsProcessed: counter,
		operationsDuration:  histogram,
		defaultURL:          "/metrics",
	}
}

// New Creates a new instance of FiberPrometheus middleware
func New(serviceName string) *FiberPrometheus {
	return create(prometheus.DefaultRegisterer, serviceName, "http", "", nil)
}

// RegisterAt will register the prometheus handler at a given URL
func (ps *FiberPrometheus) RegisterAt(app *fiber.App, url string, handlers ...fiber.Handler) {
	ps.defaultURL = url

	h := append(handlers, adaptor.HTTPHandler(promhttp.Handler()))
	app.Get(ps.defaultURL, h...)
}

func (ps *FiberPrometheus) Middleware(ctx *fiber.Ctx) error {
	start := time.Now()
	method := ctx.Route().Method

	if ctx.Route().Path == ps.defaultURL {
		return ctx.Next()
	}

	err := ctx.Next()

	status := fiber.StatusInternalServerError

	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			status = e.Code
		}
	} else {
		status = ctx.Response().StatusCode()
	}

	path := ctx.Route().Path

	statusCode := strconv.Itoa(status)
	ps.operationsProcessed.WithLabelValues(statusCode, method, path).Inc()

	elapsed := float64(time.Since(start).Nanoseconds()) / 1e9
	ps.operationsDuration.WithLabelValues(statusCode, method, path).Observe(elapsed)

	return err
}
