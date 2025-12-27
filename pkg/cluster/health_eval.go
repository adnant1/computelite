package cluster

import (
	"time"

	"github.com/adnant1/computelite/pkg/api"
)

const HeartbeatTimeout = 3 * time.Second

// EvaluateNodeHealth determines the health status of a node based on its last heartbeat timestamp
func EvaluateNodeHealth(lastHeartbeat time.Time, now time.Time) api.NodeHealth {
	if lastHeartbeat.IsZero() {
		return api.Unknown
	}

	if now.Sub(lastHeartbeat) > HeartbeatTimeout {
		return api.Unhealthy
	}

	return api.Healthy
}