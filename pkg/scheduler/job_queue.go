package scheduler

// JobQueue represents a simple FIFO queue for jobs
type JobQueue struct {
	jobs []*Job
}

// Jobs returns a snapshot of pending jobs
// this is read-only and intended for observability/debugging
func (q *JobQueue) Jobs() []*Job {
	return q.jobs
}

// Enqueue adds a job to the end of the queue
func (q *JobQueue) Enqueue(job *Job) {
	q.jobs = append(q.jobs, job)
}

// Dequeue removes and returns the job at the front of the queue
func (q *JobQueue) Dequeue() (*Job, bool) {
	if len(q.jobs) == 0 {
		return nil, false
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]

	return job, true
}

// Peek returns the job at the front of the queue without removing it
func (q *JobQueue) Peek() (*Job, bool) {
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