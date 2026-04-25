package envfilter_test

import (
	"testing"

	"github.com/yourorg/envoy-conf/internal/envfilter"
)

var sampleEnv = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_URL":       "postgres://localhost/db",
	"DB_PASSWORD":  "secret",
	"LOG_LEVEL":    "info",
	"FEATURE_FLAG": "true",
}

func TestFilter_NoOptions(t *testing.T) {
	out, err := envfilter.Filter(sampleEnv, envfilter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(sampleEnv) {
		t.Errorf("expected %d entries, got %d", len(sampleEnv), len(out))
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	out, err := envfilter.Filter(sampleEnv, envfilter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestFilter_ByPattern(t *testing.T) {
	out, err := envfilter.Filter(sampleEnv, envfilter.Options{Pattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestFilter_ExcludeKeys(t *testing.T) {
	out, err := envfilter.Filter(sampleEnv, envfilter.Options{ExcludeKeys: []string{"DB_PASSWORD", "LOG_LEVEL"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, found := out["DB_PASSWORD"]; found {
		t.Error("DB_PASSWORD should have been excluded")
	}
	if len(out) != len(sampleEnv)-2 {
		t.Errorf("expected %d entries, got %d", len(sampleEnv)-2, len(out))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := envfilter.Filter(sampleEnv, envfilter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestFilter_CombinedPrefixAndExclude(t *testing.T) {
	out, err := envfilter.Filter(sampleEnv, envfilter.Options{
		Prefix:      "DB_",
		ExcludeKeys: []string{"DB_PASSWORD"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
	if _, ok := out["DB_URL"]; !ok {
		t.Error("expected DB_URL in result")
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	out, err := envfilter.Filter(map[string]string{}, envfilter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected 0 entries for empty input, got %d", len(out))
	}
}
