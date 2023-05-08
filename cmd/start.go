package cmd

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/collector"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Logstash exporter",
	Long:  "Start the Prometheus Logstash exporter with the specified configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		startExporter(constants.LogstashURL, constants.ListenAddress)
	},
}

func startExporter(logstashURL, listenAddress string) {
	logstashCollector, err := collector.NewLogstashCollector(logstashURL)
	if err != nil {
		log.Fatalf("Cannot register a new collector: %v", err)
	}
	prometheus.MustRegister(logstashCollector)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/-/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Printf("Starting Logstash exporter on %s...", listenAddress)

	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		log.Fatalf("Error starting the HTTP server: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
