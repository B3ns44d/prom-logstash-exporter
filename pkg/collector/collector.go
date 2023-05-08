package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"prom-logstash-exporter/constants"
	"sync"
	"time"
)

type Collector interface {
	Collect(ch chan<- prometheus.Metric) error
}

var scrapeDurations = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace: constants.Namespace,
		Subsystem: "exporter",
		Name:      "scrape_duration_seconds",
		Help:      "prom_logstash_exporter: Duration of a scrape job.",
	},
	[]string{"collector", "result"},
)

type LogstashCollector struct {
	collectors map[string]Collector
}

func NewLogstashCollector(logstashEndpoint string) (*LogstashCollector, error) {
	nodeInfoCollector, err := NewNodeInfoCollector(logstashEndpoint)
	if err != nil {
		return nil, fmt.Errorf("cannot register a new collector: %v", err)
	}

	collectors := make(map[string]Collector)
	collectors["info"] = nodeInfoCollector

	return &LogstashCollector{collectors: collectors}, nil
}

func (coll LogstashCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

func (coll LogstashCollector) Collect(ch chan<- prometheus.Metric) {
	wg := &sync.WaitGroup{}
	wg.Add(len(coll.collectors))

	for name, collector := range coll.collectors {
		go coll.collectWith(name, collector, ch, wg)
	}

	wg.Wait()
	scrapeDurations.Collect(ch)
}

func (coll LogstashCollector) collectWith(name string, c Collector, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()
	err := c.Collect(ch)
	duration := time.Since(start)

	if err != nil {
		logrus.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		scrapeDurations.WithLabelValues(name, "error").Observe(duration.Seconds())
	} else {
		logrus.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		scrapeDurations.WithLabelValues(name, "success").Observe(duration.Seconds())
	}
}
