package agent

import (
	"testing"
)

func TestValidatePlannerOutput(t *testing.T) {
	tests := []struct {
		name         string
		output       string
		shouldError  bool
		expectedGoal string
	}{
		{
			name: "Valid JSON plan",
			output: `{
  "goal": "Add error handling",
  "constraints": ["Must not break tests"],
  "steps": ["Read file", "Add checks"],
  "files_to_inspect": ["main.go"]
}`,
			shouldError:  false,
			expectedGoal: "Add error handling",
		},
		{
			name: "JSON with surrounding text",
			output: `Here's the plan:
{
  "goal": "Refactor database",
  "constraints": [],
  "steps": ["Read code", "Refactor"],
  "files_to_inspect": ["db.go"]
}
That's my plan.`,
			shouldError:  false,
			expectedGoal: "Refactor database",
		},
		{
			name:        "Missing goal",
			output:      `{"constraints": [], "steps": ["step1"], "files_to_inspect": []}`,
			shouldError: true,
		},
		{
			name:        "Missing steps",
			output:      `{"goal": "test", "constraints": [], "files_to_inspect": []}`,
			shouldError: true,
		},
		{
			name:        "Invalid JSON",
			output:      `{invalid json}`,
			shouldError: true,
		},
		{
			name:        "No JSON at all",
			output:      `This is just plain text with no JSON`,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := ValidatePlannerOutput(tt.output)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else {
					t.Logf("Expected error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if plan == nil {
					t.Fatal("Expected plan, got nil")
				}
				if plan.Goal != tt.expectedGoal {
					t.Errorf("Expected goal %q, got %q", tt.expectedGoal, plan.Goal)
				}
			}
		})
	}
}

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple JSON object",
			input:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
		{
			name:     "JSON with surrounding text",
			input:    `Some text {"key": "value"} more text`,
			expected: `{"key": "value"}`,
		},
		{
			name:     "Nested JSON",
			input:    `{"outer": {"inner": "value"}}`,
			expected: `{"outer": {"inner": "value"}}`,
		},
		{
			name:     "No JSON",
			input:    `Just plain text`,
			expected: ``,
		},
		{
			name:     "Incomplete JSON",
			input:    `{"key": "value"`,
			expected: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractJSON(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPlanString(t *testing.T) {
	plan := &Plan{
		Goal:           "Test goal",
		Constraints:    []string{"constraint1", "constraint2"},
		Steps:          []string{"step1", "step2", "step3"},
		FilesToInspect: []string{"file1.go", "file2.go"},
	}

	str := plan.String()
	if str == "" {
		t.Error("String() should not be empty")
	}

	// Should contain all components
	if !contains(str, "Test goal") {
		t.Error("String() should contain goal")
	}
	if !contains(str, "constraint1") {
		t.Error("String() should contain constraints")
	}
	if !contains(str, "step1") {
		t.Error("String() should contain steps")
	}
	if !contains(str, "file1.go") {
		t.Error("String() should contain files")
	}

	t.Logf("Plan string:\n%s", str)
}

func TestPlanToJSON(t *testing.T) {
	plan := &Plan{
		Goal:           "Test goal",
		Constraints:    []string{"constraint1"},
		Steps:          []string{"step1", "step2"},
		FilesToInspect: []string{"file.go"},
	}

	jsonStr, err := plan.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}

	if jsonStr == "" {
		t.Error("ToJSON() should not return empty string")
	}

	// Should be valid JSON
	_, err = ValidatePlannerOutput(jsonStr)
	if err != nil {
		t.Errorf("ToJSON() produced invalid JSON: %v", err)
	}

	t.Logf("Plan JSON:\n%s", jsonStr)
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
