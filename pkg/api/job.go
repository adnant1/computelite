package api

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

type Job struct {
	ID       int64
	Requires Resource // Resources required by the job
	State    JobState          // Current state of the job
}