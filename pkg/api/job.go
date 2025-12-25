package api

import "time"

type JobState int

// Define job states
const (
	Pending JobState = iota
	Assigned
	Running
	Succeeded
	Failed
	Evicted
)

const RunDuration = time.Minute * 1

type Job struct {
	ID              int64
	Requires        Resource   // Resources required by the job
	State           JobState   // Current state of the job
	AssignedNodeID  string     // ID of the node to which the job is assigned

	StartedAt       time.Time  // Timestamp when the job started running
}