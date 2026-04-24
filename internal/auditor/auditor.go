package auditor

import (
	"fmt"
	"time"
)

// EventType represents the kind of audit event.
type EventType string

const (
	EventDiff     EventType = "diff"
	EventValidate EventType = "validate"
	EventResolve  EventType = "resolve"
	EventSnapshot EventType = "snapshot"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Event     EventType         `json:"event"`
	Targets   []string          `json:"targets,omitempty"`
	Meta      map[string]string `json:"meta,omitempty"`
	Success   bool              `json:"success"`
	Message   string            `json:"message,omitempty"`
}

// Log holds a list of audit entries.
type Log struct {
	Entries []Entry `json:"entries"`
}

// Auditor records audit entries for CLI operations.
type Auditor struct {
	log Log
}

// New creates a new Auditor instance.
func New() *Auditor {
	return &Auditor{}
}

// Record appends an audit entry.
func (a *Auditor) Record(event EventType, targets []string, meta map[string]string, success bool, message string) {
	a.log.Entries = append(a.log.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Targets:   targets,
		Meta:      meta,
		Success:   success,
		Message:   message,
	})
}

// Entries returns all recorded entries.
func (a *Auditor) Entries() []Entry {
	return a.log.Entries
}

// Summary returns a human-readable summary of recorded events.
func (a *Auditor) Summary() string {
	total := len(a.log.Entries)
	if total == 0 {
		return "No audit events recorded."
	}
	successes := 0
	for _, e := range a.log.Entries {
		if e.Success {
			successes++
		}
	}
	return fmt.Sprintf("%d event(s) recorded: %d succeeded, %d failed.", total, successes, total-successes)
}
