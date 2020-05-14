package web

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zergslaw/boilerplate/internal/log"
	"github.com/zergslaw/boilerplate/internal/metrics"
	"go.uber.org/zap"
)

type middlewareFunc func(http.Handler) http.Handler

// go-swagger responders panic on error while writing response to client,
// this shouldn't result in crash - unlike a real, reasonable panic.
//
// Usually it should be second middlewareFunc (after createLogger).
func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			const code = http.StatusInternalServerError
			switch err := recover(); err := err.(type) {
			default:
				metrics.PanicsTotal.Inc()
				logger := log.FromContext(r.Context())
				logger.With(
					zap.Any(log.Error, err),
				).Error("panic")

				w.WriteHeader(code)
			case nil:
			case net.Error:
				logger := log.FromContext(r.Context())
				logger.With(
					zap.Error(err),
				).Error("recovered")

				w.WriteHeader(code)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func createLogger(basePath string, logger *zap.Logger) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newLogger := logger.With(
				zap.String(log.Remote, r.RemoteAddr),
				zap.String(log.HTTPMethod, r.Method),
				zap.String(log.Func, strings.TrimPrefix(r.URL.Path, basePath)),
			)

			r = r.WithContext(log.SetContext(r.Context(), newLogger))

			next.ServeHTTP(w, r)
		})
	}
}

func accessLog(basePath string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			metric.reqInFlight.Inc()
			defer metric.reqInFlight.Dec()

			m := httpsnoop.CaptureMetrics(next, w, r)
			l := prometheus.Labels{
				resourceLabel: strings.TrimPrefix(r.URL.Path, basePath),
				methodLabel:   r.Method,
				codeLabel:     strconv.Itoa(m.Code),
			}
			metric.reqTotal.With(l).Inc()
			metric.reqDuration.With(l).Observe(m.Duration.Seconds())

			logger := log.FromContext(r.Context())
			if m.Code < http.StatusInternalServerError {
				logger.With(zap.Int(log.HTTPStatus, m.Code)).Info("handled")
			} else {
				logger.With(zap.Int(log.HTTPStatus, m.Code)).Warn("failed to handle")
			}
		})
	}
}
