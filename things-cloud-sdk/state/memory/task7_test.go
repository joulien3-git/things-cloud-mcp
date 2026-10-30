package memory

import (
	"encoding/json"
	"testing"

	things "github.com/arthursoares/things-cloud-sdk"
)

func TestStateUpdateHandlesTask7(t *testing.T) {
	t.Parallel()
	const id = "YBFp6Fbjz72UTK2cUnJ5FX"
	s := NewState()
	if err := s.Update(
		things.Item{UUID: id, Kind: things.ItemKindTask7, Action: things.ItemActionCreated, P: json.RawMessage(`{"tt":"Repeat me","st":1}`)},
	); err != nil {
		t.Fatalf("Task7 create blocked state update: %v", err)
	}
	if s.Tasks[id] == nil || s.Tasks[id].Title != "Repeat me" {
		t.Fatalf("Task7 create not applied under its own ID: %#v", s.Tasks)
	}

	// Task6 and Task7 events for the same object apply to one task.
	if err := s.Update(
		things.Item{UUID: id, Kind: things.ItemKindTask, Action: things.ItemActionModified, P: json.RawMessage(`{"tt":"Renamed"}`)},
	); err != nil {
		t.Fatal(err)
	}
	if len(s.Tasks) != 1 || s.Tasks[id].Title != "Renamed" {
		t.Fatalf("Task6 modification of a Task7 task: %#v", s.Tasks)
	}

	if err := s.Update(things.Item{UUID: id, Kind: things.ItemKindTask7, Action: things.ItemActionDeleted, P: json.RawMessage(`{}`)}); err != nil {
		t.Fatal(err)
	}
	if len(s.Tasks) != 0 {
		t.Fatalf("Task7 deletion not applied: %#v", s.Tasks)
	}
}

// TestStateUpdateMixedKindsFromIssueReports replays one batch containing every
// kind reported to break the rebuild (Task7, Task2, Tag2, Command, Command3).
func TestStateUpdateMixedKindsFromIssueReports(t *testing.T) {
	t.Parallel()
	s := NewState()
	err := s.Update(
		things.Item{UUID: "YBFp6Fbjz72UTK2cUnJ5FX", Kind: things.ItemKindTask7, Action: things.ItemActionCreated, P: json.RawMessage(`{"tt":"Task7"}`)},
		things.Item{UUID: "66848BEF-2FE7-4E88-99BD-FA933FF68997", Kind: things.ItemKindTask2, Action: things.ItemActionCreated, P: json.RawMessage(`{"tt":"Task2"}`)},
		things.Item{UUID: "13AED47B-1D02-4B11-8475-D54EB69DE64A", Kind: things.ItemKindTask2, Action: things.ItemActionModified},
		things.Item{UUID: "CC-Things-Tag-Low", Kind: things.ItemKindTag2, Action: things.ItemActionCreated, P: json.RawMessage(`{"tt":"Low"}`)},
		things.Item{UUID: "8F63FEC2-1202-43D8-80D7-55B4AC4FCA37", Kind: things.ItemKindCommand, Action: things.ItemActionCreated, P: json.RawMessage(`{"tp":0,"if":{"tt":"Fwd"}}`)},
		things.Item{UUID: "457jbco4TLEAXyjgnrgE6X", Kind: things.ItemKindCommand3, Action: things.ItemActionDeleted},
	)
	if err != nil {
		t.Fatalf("mixed batch blocked state update: %v", err)
	}
	if len(s.Tasks) != 2 || len(s.Tags) != 1 {
		t.Fatalf("tasks=%d tags=%d, want 2 tasks and 1 tag", len(s.Tasks), len(s.Tags))
	}
}
