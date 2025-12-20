package model

type JobState int

// Define job states
const (
	Pending JobState = iota
	Running
)

type Job struct {
	ID       int64     // Unique identifier for the job
	Required Resource  // Resources required by the job
	State    JobState  // Current state of the job
}