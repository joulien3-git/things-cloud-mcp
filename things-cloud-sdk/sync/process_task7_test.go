package sync

import (
	"encoding/json"
	"path/filepath"
	"testing"

	things "github.com/arthursoares/things-cloud-sdk"
)

func TestProcessTask7(t *testing.T) {
	t.Parallel()
	syncer, err := Open(filepath.Join(t.TempDir(), "test.db"), nil)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer syncer.Close()

	changes, err := syncer.processItems([]things.Item{
		{UUID: "YBFp6Fbjz72UTK2cUnJ5FX", Kind: things.ItemKindTask7, Action: things.ItemActionCreated, P: json.RawMessage(`{"tt":"Repeat me","st":1}`)},
	}, 0)
	if err != nil {
		t.Fatalf("processItems failed: %v", err)
	}
	for _, change := range changes {
		if _, unknown := change.(UnknownChange); unknown {
			t.Fatalf("Task7 produced an UnknownChange: %#v", change)
		}
	}
	task, err := syncer.getTask("YBFp6Fbjz72UTK2cUnJ5FX")
	if err != nil {
		t.Fatalf("getTask: %v", err)
	}
	if task == nil || task.Title != "Repeat me" {
		t.Fatalf("Task7 not stored as a task: %#v", task)
	}
}

func TestIsTaskKind(t *testing.T) {
	t.Parallel()
	for _, kind := range []things.ItemKind{"Task7", "Task6", "Task4", "Task3", "Task2", "Task"} {
		if !things.IsTaskKind(kind) {
			t.Errorf("IsTaskKind(%q) = false", kind)
		}
	}
	for _, kind := range []things.ItemKind{"Task8", "Tag3", "Command", ""} {
		if things.IsTaskKind(kind) {
			t.Errorf("IsTaskKind(%q) = true", kind)
		}
	}
}
