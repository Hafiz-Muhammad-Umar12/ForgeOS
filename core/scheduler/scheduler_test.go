package scheduler

import (
	"context"
	"testing"

	"github.com/Hafiz-Muhammad-Umar12/ForgeOS/core/registry"
)

func TestTaskStateValid(t *testing.T) {
	for _, s := range []TaskState{TaskStatePending, TaskStateScheduled, TaskStateRunning, TaskStateDone, TaskStateFailed} {
		if !s.Valid() {
			t.Errorf("%s invalid", s)
		}
	}
	if TaskState("bogus").Valid() {
		t.Error("bogus should be invalid")
	}
}

type fakeScheduler struct{}

func (fakeScheduler) RequestSchedule(context.Context, ScheduleRequest) (Task, error) {
	return Task{ID: "t1", State: TaskStateScheduled}, nil
}
func (fakeScheduler) Cancel(context.Context, string) error { return nil }
func (fakeScheduler) Status(context.Context, string) (Task, error) {
	return Task{ID: "t1"}, nil
}

type fakeAuthority struct{}

func (fakeAuthority) RequestSchedule(context.Context, ScheduleRequest) (Task, error) {
	return Task{ID: "t1"}, nil
}
func (fakeAuthority) RequestProvision(context.Context, string) error { return nil }
func (fakeAuthority) RegisterPlugin(context.Context, registry.ServiceInfo) error {
	return nil
}
func (fakeAuthority) Health(context.Context) (string, error) { return "ok", nil }

var (
	_ Scheduler       = (*fakeScheduler)(nil)
	_ KernelAuthority = (*fakeAuthority)(nil)
)

func TestSchedulerFake(t *testing.T) {
	f := fakeScheduler{}
	task, err := f.RequestSchedule(context.Background(), ScheduleRequest{AgentID: "a1"})
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != "t1" || task.State != TaskStateScheduled {
		t.Fatalf("task=%+v", task)
	}
}
