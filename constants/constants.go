package constants

import "time"

const Version = "v0.0.1"

var (
	LogstashURL   string
	ListenAddress string
)

const (
	Namespace = "logstash"
	StatsPath = "/_node/stats"
)

type HealthResponse struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}
