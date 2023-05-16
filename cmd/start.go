package cmd

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/collector"
	"time"
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
		logrus.Fatalf("Cannot register a new collector: %v", err)
	}
	prometheus.MustRegister(logstashCollector)
	prometheus.MustRegister(version.NewCollector("prom_logstash_exporter"))

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/-/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/-/health", func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now()
		response := constants.HealthResponse{
			Status: "OK",
			Time:   currentTime,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("Health check write response error: %v", err)
		}
	})

	logrus.Info("prom-logstash-exporter version: ", constants.Version)
	logrus.Printf("Logstash exporter is running on %s...", listenAddress)

	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		logrus.Fatalf("Error starting the HTTP server: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
}
