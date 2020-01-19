package rest

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
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
				log := logFromCtx(r.Context())
				log.WithFields(logrus.Fields{
					LogHTTPStatus: code,
					LogError:      "panic",
				}).Error(err)
				w.WriteHeader(code)
			case nil:
			case net.Error:
				log := logFromCtx(r.Context())
				log.WithFields(logrus.Fields{
					LogHTTPStatus: code,
					LogError:      "recovered",
				}).Error(err)
				w.WriteHeader(code)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

type loggerKey uint8

const logKey loggerKey = 1

func logger(basePath string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logrus.New().WithFields(logrus.Fields{
				LogRemote:     r.RemoteAddr,
				LogHTTPStatus: "",
				LogHTTPMethod: r.Method,
				LogFunc:       strings.TrimPrefix(r.URL.Path, basePath),
			})

			ctx := context.WithValue(r.Context(), logKey, log)
			r = r.WithContext(ctx)

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

			log := logFromCtx(r.Context())
			if m.Code < 500 {
				log.WithField(LogHTTPStatus, m.Code).Info("handled")
			} else {
				log.WithField(LogHTTPStatus, m.Code).Warn("failed to handle")
			}
		})
	}
}
