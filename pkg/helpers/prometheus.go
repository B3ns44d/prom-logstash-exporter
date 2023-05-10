package helpers

import (
	"github.com/prometheus/client_golang/prometheus"
)

func NewDescFQ(namespace, subsystem string) func(name, help string, labels ...string) *prometheus.Desc {
	return func(name, help string, labels ...string) *prometheus.Desc {
		return prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, name), help, labels, nil)
	}
}
