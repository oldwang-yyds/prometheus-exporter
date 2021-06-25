package main

import (
	"net/http"

	"github.com/friendly-u/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	m := metrics.Metrics
	reg := metrics.StartHandMetrics(m)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8080", nil)
}
