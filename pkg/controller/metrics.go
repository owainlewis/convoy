package controller

import "github.com/prometheus/client_golang/prometheus"

var (
	eventsProcessed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "convoy",
		Subsystem: "controller",
		Name:      "events_processed",
		Help:      "Total number of events processed",
	})
	eventsQueued = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "convoy",
		Subsystem: "controller",
		Name:      "events_queued",
		Help:      "Total number of events queued for processing",
	})
)

func init() {
	prometheus.MustRegister(eventsProcessed)
}
