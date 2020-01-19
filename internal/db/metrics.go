package db

import (
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/zergslaw/users/internal/app"
)

var metric struct { //nolint:gochecknoglobals
	callTotal    *prometheus.CounterVec
	callErrTotal *prometheus.CounterVec
	callDuration *prometheus.HistogramVec
}

const (
	methodLabel = "method"
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics(namespace string) {
	const subsystem = "db"

	metric.callTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "call_total",
			Help:      "Amount of Repo calls.",
		},
		[]string{methodLabel},
	)
	metric.callErrTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "errors_total",
			Help:      "Amount of Repo errors.",
		},
		[]string{methodLabel},
	)
	metric.callDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "call_duration_seconds",
			Help:      "Repo call latency.",
		},
		[]string{methodLabel},
	)

	for _, method := range methodsOf(new(app.Repo)) {
		l := prometheus.Labels{
			methodLabel: method,
		}
		metric.callTotal.With(l)
		metric.callErrTotal.With(l)
		metric.callDuration.With(l)
	}
}

func methodsOf(v interface{}) []string {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Interface {
		panic("require pointer to interface")
	}
	typ = typ.Elem()
	methods := make([]string, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		methods[i] = typ.Method(i).Name
	}
	return methods
}

// Usage:
//	func (…) SomeMethod(…) (err error) {
//		methodName, methodDone := methodMetrics(0)
//		defer methodDone(&err)
//		…
//	}
func methodMetrics(skip int) (name string, done func(*error)) {
	pc, _, _, _ := runtime.Caller(1 + skip)
	names := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	method := names[len(names)-1]
	start := time.Now()
	l := prometheus.Labels{methodLabel: method}
	return method, func(refErr *error) {
		metric.callTotal.With(l).Inc()
		metric.callDuration.With(l).Observe(time.Since(start).Seconds())
		if refErr != nil && *refErr != nil {
			metric.callErrTotal.With(l).Inc()
		}
	}
}
