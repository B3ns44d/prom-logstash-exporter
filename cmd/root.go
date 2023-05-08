package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"prom-logstash-exporter/constants"
)

var rootCmd = &cobra.Command{
	Use:   "prom-logstash-exporter",
	Short: "A Prometheus exporter for Logstash metrics",
	Long: `prom-logstash-exporter is a Prometheus exporter that collects
Logstash metrics using the Logstash monitoring API and exposes them
for consumption by Prometheus.`,
	Version: constants.Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prom-logstash-exporter version", constants.Version)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&constants.LogstashURL, "logstash-url", "http://localhost:9600", "URL of the Logstash instance to monitor")
	startCmd.PersistentFlags().StringVar(&constants.ListenAddress, "listen-address", ":2112", "The address to listen on for Prometheus metrics")
}
