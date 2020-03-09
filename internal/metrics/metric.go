// Package metrics contains the common metrics needed for different packages.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// PanicsTotal contains metrics for rates of panic.
	PanicsTotal struct{ prometheus.Counter } //nolint:gochecknoglobals,gocritic
)

// InitMetrics must be called once before using this package.
// It registers and initializes metrics used by this package.
func InitMetrics() {
	PanicsTotal.Counter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "panics_total",
			Help: "Amount of recovered panics.",
		},
	)
}
