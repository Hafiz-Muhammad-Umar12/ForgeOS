package domain

import "testing"

type fakeAggregate struct{ id, typ string }

func (f fakeAggregate) AggregateID() string   { return f.id }
func (f fakeAggregate) AggregateType() string { return f.typ }

var _ Aggregate = fakeAggregate{}

type idCase struct {
	id interface {
		IsValid() bool
		String() string
	}
	valid bool
	val   string
}

func TestIdentifiers(t *testing.T) {
	cases := []idCase{
		{OrgID("org-1"), true, "org-1"},
		{OrgID(""), false, ""},
		{ProjectID("p"), true, "p"},
		{IntentID(""), false, ""},
		{PlanID("pl"), true, "pl"},
		{TaskID("t"), true, "t"},
		{AgentID("ag"), true, "ag"},
		{WorkspaceID("w"), true, "w"},
	}
	for _, c := range cases {
		if c.id.IsValid() != c.valid {
			t.Errorf("%v valid=%v want %v", c.id, c.id.IsValid(), c.valid)
		}
		if c.id.String() != c.val {
			t.Errorf("%v string=%q want %q", c.id, c.id.String(), c.val)
		}
	}
}
