package envdrift

import (
	"testing"
)

func TestDetect_NoChanges(t *testing.T) {
	baseline := map[string]string{"A": "1", "B": "2"}
	live := map[string]string{"A": "1", "B": "2"}
	report := Detect("prod", baseline, live)
	if report.HasDrift() {
		t.Fatal("expected no drift")
	}
	a, r, c := report.Counts()
	if a != 0 || r != 0 || c != 0 {
		t.Fatalf("unexpected counts: added=%d removed=%d changed=%d", a, r, c)
	}
}

func TestDetect_Added(t *testing.T) {
	baseline := map[string]string{"A": "1"}
	live := map[string]string{"A": "1", "B": "new"}
	report := Detect("prod", baseline, live)
	if !report.HasDrift() {
		t.Fatal("expected drift")
	}
	a, _, _ := report.Counts()
	if a != 1 {
		t.Fatalf("expected 1 added, got %d", a)
	}
}

func TestDetect_Removed(t *testing.T) {
	baseline := map[string]string{"A": "1", "GONE": "bye"}
	live := map[string]string{"A": "1"}
	report := Detect("staging", baseline, live)
	_, r, _ := report.Counts()
	if r != 1 {
		t.Fatalf("expected 1 removed, got %d", r)
	}
}

func TestDetect_Changed(t *testing.T) {
	baseline := map[string]string{"HOST": "localhost"}
	live := map[string]string{"HOST": "prod.example.com"}
	report := Detect("prod", baseline, live)
	_, _, c := report.Counts()
	if c != 1 {
		t.Fatalf("expected 1 changed, got %d", c)
	}
	entry := report.Entries[0]
	if entry.Baseline != "localhost" || entry.Live != "prod.example.com" {
		t.Fatalf("unexpected entry values: %+v", entry)
	}
}

func TestDetect_SortedKeys(t *testing.T) {
	baseline := map[string]string{"Z": "1", "A": "2", "M": "3"}
	live := map[string]string{"Z": "1", "A": "2", "M": "3"}
	report := Detect("prod", baseline, live)
	for i := 1; i < len(report.Entries); i++ {
		if report.Entries[i].Key < report.Entries[i-1].Key {
			t.Fatal("entries are not sorted")
		}
	}
}

func TestReport_Target(t *testing.T) {
	report := Detect("canary", map[string]string{}, map[string]string{})
	if report.Target != "canary" {
		t.Fatalf("expected target canary, got %s", report.Target)
	}
}
