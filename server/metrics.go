package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var dataMessagesAddedSum = promauto.NewCounter(prometheus.CounterOpts{
	Name: "data_messages_added_sum",
	Help: "Lifetime sum of message additions",
})

var dataMessagesDeletedSum = promauto.NewCounter(prometheus.CounterOpts{
	Name: "data_messages_deleted_sum",
	Help: "Lifetime sum of message deletions",
})

var httpRequestDurationSecondsSum = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "http_request_duration_seconds_sum",
	Help: "Cumulative duration in seconds of requests",
}, []string{"method", "path"})

var httpRequestDurationSecondsCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "http_request_duration_seconds_count",
	Help: "Cumulative count of requests",
}, []string{"method", "path"})
