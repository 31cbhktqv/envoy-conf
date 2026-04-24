package auditor

import (
	"testing"
)

func TestRecord_SingleEntry(t *testing.T) {
	a := New()
	a.Record(EventDiff, []string{"staging", "prod"}, nil, true, "")

	entries := a.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Event != EventDiff {
		t.Errorf("expected event %q, got %q", EventDiff, entries[0].Event)
	}
	if !entries[0].Success {
		t.Errorf("expected success=true")
	}
	if len(entries[0].Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(entries[0].Targets))
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	a := New()
	a.Record(EventValidate, []string{"dev"}, nil, true, "all rules passed")
	a.Record(EventResolve, []string{"prod"}, nil, false, "missing source")

	if len(a.Entries()) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(a.Entries()))
	}
}

func TestRecord_WithMeta(t *testing.T) {
	a := New()
	meta := map[string]string{"file": "rules.yaml", "keys": "5"}
	a.Record(EventValidate, nil, meta, true, "")

	entry := a.Entries()[0]
	if entry.Meta["file"] != "rules.yaml" {
		t.Errorf("expected meta file=rules.yaml, got %q", entry.Meta["file"])
	}
}

func TestSummary_NoEvents(t *testing.T) {
	a := New()
	got := a.Summary()
	if got != "No audit events recorded." {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummary_MixedResults(t *testing.T) {
	a := New()
	a.Record(EventDiff, nil, nil, true, "")
	a.Record(EventSnapshot, nil, nil, false, "write error")
	a.Record(EventResolve, nil, nil, true, "")

	got := a.Summary()
	expected := "3 event(s) recorded: 2 succeeded, 1 failed."
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRecord_TimestampSet(t *testing.T) {
	a := New()
	a.Record(EventDiff, nil, nil, true, "")
	if a.Entries()[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
