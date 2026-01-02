package beads

import (
	"testing"

	"github.com/intent-solutions-io/gastown-viewer-intent/internal/model"
)

func TestParseIssueList(t *testing.T) {
	input := []byte(`[
		{
			"id": "test-1",
			"title": "Test Issue",
			"description": "Test description\n\nDone when:\n- First item\n- Second item",
			"status": "open",
			"priority": 1,
			"issue_type": "task",
			"created_at": "2026-01-01T10:00:00Z",
			"updated_at": "2026-01-01T12:00:00Z"
		}
	]`)

	issues, err := ParseIssueList(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}

	issue := issues[0]
	if issue.ID != "test-1" {
		t.Errorf("expected ID 'test-1', got '%s'", issue.ID)
	}
	if issue.Title != "Test Issue" {
		t.Errorf("expected title 'Test Issue', got '%s'", issue.Title)
	}
	if issue.Priority != 1 {
		t.Errorf("expected priority 1, got %d", issue.Priority)
	}
}

func TestMapStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected model.Status
	}{
		{"open", model.StatusPending},
		{"pending", model.StatusPending},
		{"in_progress", model.StatusInProgress},
		{"in-progress", model.StatusInProgress},
		{"closed", model.StatusDone},
		{"done", model.StatusDone},
		{"blocked", model.StatusBlocked},
		{"unknown", model.StatusPending},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapStatus(tt.input)
			if result != tt.expected {
				t.Errorf("mapStatus(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMapPriority(t *testing.T) {
	tests := []struct {
		input    int
		expected model.Priority
	}{
		{1, model.PriorityHigh},
		{2, model.PriorityMedium},
		{3, model.PriorityLow},
		{0, model.PriorityMedium},
		{99, model.PriorityMedium},
	}

	for _, tt := range tests {
		result := mapPriority(tt.input)
		if result != tt.expected {
			t.Errorf("mapPriority(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestParseDoneWhen(t *testing.T) {
	description := `Implement the feature.

Done when:
- First criterion
- Second criterion
- Third criterion

Additional notes here.`

	items := parseDoneWhen(description)

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	expected := []string{"First criterion", "Second criterion", "Third criterion"}
	for i, item := range items {
		if item != expected[i] {
			t.Errorf("item %d: expected %q, got %q", i, expected[i], item)
		}
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"bd version 0.29.0 (dev)\n", "0.29.0"},
		{"bd version 1.0.0\n", "1.0.0"},
		{"0.42.0", "0.42.0"},
	}

	for _, tt := range tests {
		result := ParseVersion([]byte(tt.input))
		if result != tt.expected {
			t.Errorf("ParseVersion(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestBDIssueToModelIssue(t *testing.T) {
	bdIssue := BDIssue{
		ID:          "test-1",
		Title:       "Test Issue",
		Description: "Done when:\n- Item 1\n- Item 2",
		Status:      "in_progress",
		Priority:    1,
		Dependencies: []BDIssue{
			{ID: "dep-1", Title: "Dependency", Status: "closed", Priority: 2, DepType: "blocks"},
		},
		Dependents: []BDIssue{
			{ID: "dep-2", Title: "Dependent", Status: "open", Priority: 2, DepType: "blocks"},
		},
	}

	issue := bdIssue.ToModelIssue()

	if issue.ID != "test-1" {
		t.Errorf("expected ID 'test-1', got '%s'", issue.ID)
	}
	if issue.Status != model.StatusInProgress {
		t.Errorf("expected status in_progress, got %s", issue.Status)
	}
	if issue.Priority != model.PriorityHigh {
		t.Errorf("expected priority high, got %s", issue.Priority)
	}
	if len(issue.DoneWhen) != 2 {
		t.Errorf("expected 2 done_when items, got %d", len(issue.DoneWhen))
	}
	if len(issue.BlockedBy) != 1 {
		t.Errorf("expected 1 blocked_by, got %d", len(issue.BlockedBy))
	}
	if len(issue.Blocks) != 1 {
		t.Errorf("expected 1 blocks, got %d", len(issue.Blocks))
	}
}

func TestParseEmptyList(t *testing.T) {
	issues, err := ParseIssueList([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected empty list, got %d issues", len(issues))
	}
}
