package envwatch

import (
	"errors"
	"testing"
	"time"
)

func mapPoll(maps []map[string]string) (PollFunc, *int) {
	call := 0
	return func() (map[string]string, error) {
		idx := call
		call++
		if idx >= len(maps) {
			return maps[len(maps)-1], nil
		}
		return maps[idx], nil
	}, &call
}

func TestDiff_NoChanges(t *testing.T) {
	result := diff(
		map[string]string{"A": "1"},
		map[string]string{"A": "1"},
	)
	if len(result) != 0 {
		t.Fatalf("expected no changes, got %d", len(result))
	}
}

func TestDiff_Added(t *testing.T) {
	result := diff(
		map[string]string{},
		map[string]string{"NEW": "val"},
	)
	if len(result) != 1 || result[0].Type != Added {
		t.Fatalf("expected one Added change, got %+v", result)
	}
}

func TestDiff_Removed(t *testing.T) {
	result := diff(
		map[string]string{"OLD": "val"},
		map[string]string{},
	)
	if len(result) != 1 || result[0].Type != Removed {
		t.Fatalf("expected one Removed change, got %+v", result)
	}
}

func TestDiff_Changed(t *testing.T) {
	result := diff(
		map[string]string{"K": "old"},
		map[string]string{"K": "new"},
	)
	if len(result) != 1 || result[0].Type != Changed {
		t.Fatalf("expected one Changed, got %+v", result)
	}
	if result[0].OldVal != "old" || result[0].NewVal != "new" {
		t.Fatalf("unexpected values: %+v", result[0])
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	poll, _ := mapPoll([]map[string]string{
		{"A": "1"},
		{"A": "2"},
	})
	opts := Options{Interval: 10 * time.Millisecond, MaxPolls: 1}
	done := make(chan struct{})
	defer close(done)

	ch, errs := Watch(poll, opts, done)
	select {
	case changes := <-ch:
		if len(changes) != 1 || changes[0].Type != Changed {
			t.Fatalf("unexpected changes: %+v", changes)
		}
	case err := <-errs:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change")
	}
}

func TestWatch_InitialPollError(t *testing.T) {
	poll := func() (map[string]string, error) {
		return nil, errors.New("boom")
	}
	opts := Options{Interval: 10 * time.Millisecond}
	done := make(chan struct{})
	defer close(done)

	_, errs := Watch(poll, opts, done)
	select {
	case err := <-errs:
		if err == nil {
			t.Fatal("expected error")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for error")
	}
}
