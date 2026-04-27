package envpromote

import (
	"strings"
	"testing"
)

func TestRender_Ready(t *testing.T) {
	r := PromoteResult{
		From:  "staging",
		To:    "production",
		Ready: true,
	}
	var sb strings.Builder
	Render(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "READY") {
		t.Errorf("expected READY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "staging") || !strings.Contains(out, "production") {
		t.Errorf("expected stage names in output")
	}
}

func TestRender_NotReady_ShowsMissing(t *testing.T) {
	r := PromoteResult{
		From:        "staging",
		To:          "production",
		Ready:       false,
		MissingKeys: []string{"DB_URL", "API_KEY"},
	}
	var sb strings.Builder
	Render(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "NOT READY") {
		t.Errorf("expected NOT READY in output")
	}
	if !strings.Contains(out, "DB_URL") || !strings.Contains(out, "API_KEY") {
		t.Errorf("expected missing keys in output")
	}
}

func TestRender_ShowsChangedAndNew(t *testing.T) {
	r := PromoteResult{
		From:        "staging",
		To:          "production",
		Ready:       true,
		ChangedKeys: []string{"TIMEOUT"},
		NewKeys:     []string{"FEATURE_FLAG"},
	}
	var sb strings.Builder
	Render(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "TIMEOUT") {
		t.Errorf("expected TIMEOUT in changed keys section")
	}
	if !strings.Contains(out, "FEATURE_FLAG") {
		t.Errorf("expected FEATURE_FLAG in new keys section")
	}
}
