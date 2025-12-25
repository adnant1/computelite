package scheduler

import "github.com/adnant1/computelite/pkg/api"

// JobQueue represents a simple FIFO queue for jobs
type JobQueue struct {
	jobs []*api.Job
}

// Jobs returns a snapshot of pending jobs
// this is read-only and intended for observability/debugging
func (q *JobQueue) Jobs() []*api.Job {
	return q.jobs
}

// Enqueue adds a job to the end of the queue
func (q *JobQueue) Enqueue(job *api.Job) {
	q.jobs = append(q.jobs, job)
}

// Dequeue removes and returns the job at the front of the queue
func (q *JobQueue) Dequeue() (*api.Job, bool) {
	if len(q.jobs) == 0 {
		return nil, false
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]

	return job, true
}

// Peek returns the job at the front of the queue without removing it
func (q *JobQueue) Peek() (*api.Job, bool) {
	if len(q.jobs) == 0 {
		return nil, false
	}

	job := q.jobs[0]
	return job, true
}

// Size returns the number of jobs in the queue
func (q *JobQueue) Size() int {
	return len(q.jobs)
}