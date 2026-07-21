// Package scheduler defines the scheduling and runtime-authority ports of the
// DevOS kernel (SDD §11). It is the contract only; the stateful DAG executor
// and scheduler implementation are delivered in later milestones (ADR-002).
//
// The KernelAuthority is the sole runtime authority tenants call to schedule
// work, provision workspaces, and register plugins (governance/04-engineering
// -standards.md §11). Implementations must enforce Constitutional tenets
// (budget T5, HITL T3, isolation T4) at this layer.
package scheduler

import (
	"context"
	"time"

	"github.com/Hafiz-Muhammad-Umar12/ForgeOS/core/registry"
)

// TaskState is the lifecycle state of a scheduled task.
type TaskState string

const (
	// TaskStatePending means the task is queued but not yet scheduled.
	TaskStatePending TaskState = "pending"
	// TaskStateScheduled means the task has been assigned to an agent.
	TaskStateScheduled TaskState = "scheduled"
	// TaskStateRunning means the agent is executing the task.
	TaskStateRunning TaskState = "running"
	// TaskStateDone means the task completed successfully.
	TaskStateDone TaskState = "done"
	// TaskStateFailed means the task failed.
	TaskStateFailed TaskState = "failed"
)

// Valid reports whether s is a known task state.
func (s TaskState) Valid() bool {
	switch s {
	case TaskStatePending, TaskStateScheduled, TaskStateRunning, TaskStateDone, TaskStateFailed:
		return true
	default:
		return false
	}
}

// Task is a unit of scheduled work assigned to an agent.
type Task struct {
	ID        string
	OrgID     string
	ProjectID string
	AgentID   string
	IntentID  string
	State     TaskState
	Payload   any
	CreatedAt time.Time
}

// ScheduleRequest asks the kernel to schedule work.
type ScheduleRequest struct {
	OrgID     string
	ProjectID string
	AgentID   string
	IntentID  string
	Payload   any
}

// Scheduler assigns and tracks tasks.
type Scheduler interface {
	RequestSchedule(ctx context.Context, req ScheduleRequest) (Task, error)
	Cancel(ctx context.Context, taskID string) error
	Status(ctx context.Context, taskID string) (Task, error)
}

// KernelAuthority is the tenant-facing runtime authority. Tenant services
// request; the kernel authorizes and executes (SDD §11, governance §11).
type KernelAuthority interface {
	RequestSchedule(ctx context.Context, req ScheduleRequest) (Task, error)
	RequestProvision(ctx context.Context, workspaceID string) error
	RegisterPlugin(ctx context.Context, info registry.ServiceInfo) error
	Health(ctx context.Context) (string, error)
}
