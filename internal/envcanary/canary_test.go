package envcanary

import (
	"testing"
)

func TestCheck_RequiredKeys_AllPresent(t *testing.T) {
	baseline := map[string]string{"APP_ENV": "staging", "PORT": "8080"}
	current := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	opts := Options{RequiredKeys: []string{"APP_ENV", "PORT"}}

	results := Check(baseline, current, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != StatusOK {
			t.Errorf("key %q: expected OK, got %v", r.Key, r.Status)
		}
	}
}

func TestCheck_RequiredKey_Missing_Critical(t *testing.T) {
	baseline := map[string]string{"APP_ENV": "staging"}
	current := map[string]string{}
	opts := Options{RequiredKeys: []string{"APP_ENV"}}

	results := Check(baseline, current, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != StatusCritical {
		t.Errorf("expected Critical, got %v", results[0].Status)
	}
	if !HasCritical(results) {
		t.Error("HasCritical should return true")
	}
}

func TestCheck_RequiredKey_Missing_AllowMissing(t *testing.T) {
	baseline := map[string]string{"APP_ENV": "staging"}
	current := map[string]string{}
	opts := Options{RequiredKeys: []string{"APP_ENV"}, AllowMissing: true}

	results := Check(baseline, current, opts)
	if results[0].Status != StatusWarning {
		t.Errorf("expected Warning, got %v", results[0].Status)
	}
	if HasCritical(results) {
		t.Error("HasCritical should return false when AllowMissing")
	}
}

func TestCheck_WatchKey_Changed(t *testing.T) {
	baseline := map[string]string{"DB_HOST": "db-staging"}
	current := map[string]string{"DB_HOST": "db-prod"}
	opts := Options{WatchKeys: []string{"DB_HOST"}}

	results := Check(baseline, current, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != StatusWarning {
		t.Errorf("expected Warning, got %v", results[0].Status)
	}
}

func TestCheck_WatchKey_Unchanged(t *testing.T) {
	baseline := map[string]string{"DB_HOST": "db-prod"}
	current := map[string]string{"DB_HOST": "db-prod"}
	opts := Options{WatchKeys: []string{"DB_HOST"}}

	results := Check(baseline, current, opts)
	if results[0].Status != StatusOK {
		t.Errorf("expected OK, got %v", results[0].Status)
	}
}

func TestCheck_WatchKey_AbsentInCurrent(t *testing.T) {
	baseline := map[string]string{"FEATURE_FLAG": "true"}
	current := map[string]string{}
	opts := Options{WatchKeys: []string{"FEATURE_FLAG"}}

	results := Check(baseline, current, opts)
	if results[0].Status != StatusWarning {
		t.Errorf("expected Warning, got %v", results[0].Status)
	}
}

func TestHasCritical_False(t *testing.T) {
	results := []Result{
		{Key: "A", Status: StatusOK},
		{Key: "B", Status: StatusWarning},
	}
	if HasCritical(results) {
		t.Error("expected HasCritical to be false")
	}
}
