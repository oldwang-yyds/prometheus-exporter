package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ManagerTotal struct {
	PingTotal *prometheus.Desc
}

type Metrics struct {
	Help  string
	Type  string
	Zone  string
	Label map[string](string)
}

var m Metrics

// Simulate prepare the data
func (c *ManagerTotal) SetYourVariable() (oomCountByHost map[string]int) {
	// Just example fake data.
	oomCountByHost = map[string]int{
		"www.baidu.com":   0,
		"bar.example.org": 2001,
	}
	return oomCountByHost
}

// Describe simply sends the two Descs in the struct to the channel.
func (c *ManagerTotal) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.PingTotal
}

func (c *ManagerTotal) Collect(ch chan<- prometheus.Metric) {
	oomCountByHost := c.SetYourVariable()
	for host, oomCount := range oomCountByHost {
		ch <- prometheus.MustNewConstMetric(
			c.PingTotal,
			prometheus.CounterValue,
			float64(oomCount),
			host,
		)
	}
}

func NewManagerTotal(zone string) *ManagerTotal {
	return &ManagerTotal{
		PingTotal: prometheus.NewDesc(
			"friendly_u_ping_total",
			"how many times about ping.",
			[]string{"host"},
			prometheus.Labels{"zone": zone},
		),
	}
}

func StartHandMetrics(mm Metrics) *prometheus.Registry {
	m := mm
	total := NewManagerTotal(m.Zone)

	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(total)

	return reg
}
