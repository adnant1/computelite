package cluster

import "github.com/adnant1/computelite/pkg/scheduler"

// JobQueue represents a simple FIFO queue for jobs
type JobQueue struct {
	jobs []*scheduler.Job
}

// Enqueue adds a job to the end of the queue
func (q *JobQueue) Enqueue(job *scheduler.Job) {
	q.jobs = append(q.jobs, job)
}

// Dequeue removes and returns the job at the front of the queue
func (q *JobQueue) Dequeue() (*scheduler.Job, bool) {
	if len(q.jobs) == 0 {
		return nil, false
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]

	return job, true
}

// Peek returns the job at the front of the queue without removing it
func (q *JobQueue) Peek() (*scheduler.Job, bool) {
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