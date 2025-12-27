# ComputeLite

**ComputeLite** is a Kubernetes-inspired compute scheduler and control-plane simulator written in Go.
It models nodes, jobs, heartbeats, and controller-based reconciliation to demonstrate how modern compute platforms schedule, monitor, and reschedule workloads under failure and churn.

The system runs as a **long-lived daemon**, converging cluster state over time rather than executing a single scheduling pass.

---

## Motivation

Most scheduling demos stop at “assign jobs to nodes.”
ComputeLite goes further by modeling:

- continuous reconciliation instead of one-shot scheduling
- node health via heartbeats
- failure detection based on time
- job eviction and rescheduling
- clean startup and graceful shutdown

This project is designed to explore **compute infrastructure fundamentals**, similar to the internals of Kubernetes or Borg without any external dependencies.

---

## High-Level Architecture

ComputeLite is structured as a **control plane + node agents** system.

### Core Components

**ClusterState**

- Single source of truth
- Thread-safe (RWMutex)
- Owns all nodes, jobs, and resource accounting
- Enforces valid job state transitions

**Controllers (Reconciliation Loops)**

- Run continuously in goroutines
- Observe cluster state via snapshots
- Mutate state only through ClusterState APIs

Controllers include:

- **SchedulerController** – assigns pending jobs to healthy nodes using a pluggable policy
- **JobController** – advances jobs through Assigned → Running → Succeeded
- **HealthController** – evaluates node health based on heartbeat freshness
- **ReschedulerController** – evicts and requeues jobs from unhealthy nodes

**Node Agents (Simulated)**

- Periodically emit heartbeats
- Represent node-local agents (similar to kubelet)
- Failure is modeled by stopping heartbeats

---

## Execution Model

ComputeLite runs as a **daemon**, not a script.

1. ClusterState is initialized
2. Controllers start and reconcile continuously
3. Node agents emit heartbeats
4. A scenario injects nodes, jobs, and failures
5. State converges over time
6. The system shuts down gracefully on SIGINT/SIGTERM

There is no “end” condition—only convergence.

---

## Project Structure

```text
computelite/
├── cmd/
│   └── computelite/
│       ├── main.go          # minimal entrypoint
│       └── app/             # binary-specific runtime logic
│           ├── run.go
│           ├── controllers.go
│           ├── reporting.go
│           ├── scenario.go
│           └── node_agent.go
├── pkg/
│   ├── api/                 # core types (Job, Node, Resource, states)
│   ├── cluster/             # ClusterState + snapshot reporting
│   ├── controller/          # control-plane controllers
│   └── scheduler/           # scheduling policies (e.g. BestFit)
```

This layout mirrors real infrastructure projects.

---

## Scheduling Policy

ComputeLite supports **pluggable scheduling policies** via a simple interface.

Current implementation:

- **FirstFit** - selects the first healthy node that can satisfy the job’s resource requirements
- **BestFit** – selects the node that minimizes leftover resources after placement
- **RoundRobin** - distributes jobs evenly across healthy nodes in cyclic order

Policies are injected into the scheduler controller and can be swapped without changing core logic.

---

## Failure & Rescheduling

- Nodes emit heartbeats at fixed intervals
- HealthController marks nodes unhealthy when heartbeats exceed a timeout
- ReschedulerController evicts jobs from unhealthy nodes
- Evicted jobs are requeued and rescheduled on healthy nodes

Failures emerge naturally from time—not explicit conditionals.

---

## Observability

The system periodically prints a **cluster snapshot**:

- job counts by state
- node resource utilization
- node health and heartbeat status

This provides a live view of system convergence while running.

---

## Running the Project

From the repository root:

```bash
go run ./cmd/computelite
```

The system will:

- start controllers
- start node agents
- run the default scenario
- print periodic cluster snapshots

To stop:

```text
Ctrl + C
```

Shutdown is graceful and deterministic.

---

## Design Principles

- Controller-based reconciliation
- Clear ownership of shared state
- No direct shared-memory access outside ClusterState
- Time-based failure detection
- Idempotent, race-safe controllers
- Minimal `main.go`, logic lives in `app/`

---

## Future Extensions (Optional)

ComputeLite is intentionally extensible. Possible next steps include:

- graceful node draining
- job retry backoff
- event-driven scheduler wakeups
- multiple scheduling policies
- metrics instead of logs
- CLI flags for scenario selection

---

## Summary

ComputeLite demonstrates how a real compute control plane behaves:

- always running
- always converging
- resilient to failure
- cleanly structured
