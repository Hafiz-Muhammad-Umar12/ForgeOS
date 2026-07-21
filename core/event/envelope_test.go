package event

import (
	"testing"
	"time"
)

func TestNewEnvelope(t *testing.T) {
	now := time.Date(2026, 7, 21, 12, 0, 0, 0, time.UTC)
	e := New(TypeIntentCreated, "ingress", map[string]string{"x": "y"},
		WithTraceID("trace-1"), WithOrgID("org-1"), WithProjectID("proj-1"), WithTime(now))
	if e.ID == "" {
		t.Fatal("ID empty")
	}
	if e.Type != TypeIntentCreated {
		t.Errorf("type=%s", e.Type)
	}
	if e.SchemaVersion != CurrentSchemaVersion {
		t.Errorf("schema=%d", e.SchemaVersion)
	}
	if e.TraceID != "trace-1" {
		t.Errorf("trace=%s", e.TraceID)
	}
	if e.OrgID != "org-1" {
		t.Errorf("org=%s", e.OrgID)
	}
	if e.ProjectID != "proj-1" {
		t.Errorf("proj=%s", e.ProjectID)
	}
	if e.ProducedBy != "ingress" {
		t.Errorf("by=%s", e.ProducedBy)
	}
	if !e.ProducedAt.Equal(now) {
		t.Errorf("time=%v", e.ProducedAt)
	}
	if e.Payload["x"] != "y" {
		t.Errorf("payload=%v", e.Payload)
	}
}

func TestSubjectRoundTrip(t *testing.T) {
	for _, typ := range []EventType{
		TypeIntentCreated, TypePlanProposed, TypeTaskAssigned, TypeAgentToken,
		TypeDeployCompleted, TypeBudgetExceeded, TypeWorkspaceReady, TypeAgentRegistered,
	} {
		s := typ.Subject()
		if s.String() != "devos."+string(typ) {
			t.Errorf("%s subject=%s", typ, s.String())
		}
		if s.EventType() != typ {
			t.Errorf("%s roundtrip=%s", typ, s.EventType())
		}
	}
}

func TestUniqueIDs(t *testing.T) {
	a := New(TypeIntentCreated, "x", "")
	b := New(TypeIntentCreated, "x", "")
	if a.ID == b.ID {
		t.Fatal("expected unique IDs")
	}
}

func TestWithTimeDeterminism(t *testing.T) {
	now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	e := New(TypePlanApproved, "x", "", WithTime(now))
	if !e.ProducedAt.Equal(now) {
		t.Fatalf("time=%v", e.ProducedAt)
	}
}
