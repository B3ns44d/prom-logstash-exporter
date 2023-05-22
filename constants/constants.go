package constants

import "time"

const Version = "v1.0.0"

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
