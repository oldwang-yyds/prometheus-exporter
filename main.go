package main

import (
	"metrics/metrics"
	"metrics/ping"
)

func main() {
	go ping.Start()

	metrics.Start()
}
