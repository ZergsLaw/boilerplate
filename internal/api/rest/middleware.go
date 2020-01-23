package rest

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/zergslaw/users/internal/log"
	"github.com/zergslaw/users/internal/metrics"
)

type middlewareFunc func(http.Handler) http.Handler

// go-swagger responders panic on error while writing response to client,
// this shouldn't result in crash - unlike a real, reasonable panic.
//
// Usually it should be second middlewareFunc (after logger).
func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			const code = http.StatusInternalServerError
			switch err := recover(); err := err.(type) {
			default:
				metrics.PanicsTotal.Inc()
				logger := log.FromContext(r.Context())
				logger.WithFields(logrus.Fields{
					log.HTTPStatus: code,
					log.Error:      "panic",
				}).Error(err)
				w.WriteHeader(code)
			case nil:
			case net.Error:
				logger := log.FromContext(r.Context())
				logger.WithFields(logrus.Fields{
					log.HTTPStatus: code,
					log.Error:      "recovered",
				}).Error(err)
				w.WriteHeader(code)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func logger(basePath string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logrus.New().WithFields(logrus.Fields{
				log.Remote:     r.RemoteAddr,
				log.HTTPStatus: "",
				log.HTTPMethod: r.Method,
				log.Func:       strings.TrimPrefix(r.URL.Path, basePath),
				log.API:        "rest",
			})

			r = r.WithContext(log.SetContext(r.Context(), logger))

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
			if m.Code < 500 {
				logger.WithField(log.HTTPStatus, m.Code).Info("handled")
			} else {
				logger.WithField(log.HTTPStatus, m.Code).Warn("failed to handle")
			}
		})
	}
}
