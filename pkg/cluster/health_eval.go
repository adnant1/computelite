package cluster

import "time"

const HeartbeatTimeout = 15 * time.Second

// EvaluateNodeHealth determines the health status of a node based on its last heartbeat timestamp
func EvaluateNodeHealth(lastHeartbeat time.Time, now time.Time) NodeHealth {
	if lastHeartbeat.IsZero() {
		return Unknown
	}

	if now.Sub(lastHeartbeat) > HeartbeatTimeout {
		return Unhealthy
	}

	return Healthy
}