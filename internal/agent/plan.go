package agent

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Plan represents the structured output from the planner model
type Plan struct {
	Goal           string   `json:"goal"`
	Constraints    []string `json:"constraints"`
	Steps          []string `json:"steps"`
	FilesToInspect []string `json:"files_to_inspect"`
}

// ValidatePlannerOutput extracts and validates JSON plan from planner output
func ValidatePlannerOutput(output string) (*Plan, error) {
	// Clean the output
	output = strings.TrimSpace(output)

	// Try to extract JSON from the output
	jsonStr := extractJSON(output)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in planner output")
	}

	// Parse JSON
	var plan Plan
	if err := json.Unmarshal([]byte(jsonStr), &plan); err != nil {
		return nil, fmt.Errorf("failed to parse planner JSON: %w\nOutput: %s", err, jsonStr)
	}

	// Validate required fields
	if plan.Goal == "" {
		return nil, fmt.Errorf("plan missing required field: goal")
	}

	if len(plan.Steps) == 0 {
		return nil, fmt.Errorf("plan missing required field: steps")
	}

	return &plan, nil
}

// extractJSON attempts to extract JSON from text that may contain other content
func extractJSON(text string) string {
	// Try to find JSON object boundaries
	start := strings.Index(text, "{")
	if start == -1 {
		return ""
	}

	// Find the matching closing brace
	depth := 0
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return text[start : i+1]
			}
		}
	}

	return ""
}

// String returns a human-readable representation of the plan
func (p *Plan) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Goal: %s\n", p.Goal))

	if len(p.Constraints) > 0 {
		sb.WriteString("\nConstraints:\n")
		for _, c := range p.Constraints {
			sb.WriteString(fmt.Sprintf("  - %s\n", c))
		}
	}

	if len(p.Steps) > 0 {
		sb.WriteString("\nSteps:\n")
		for i, s := range p.Steps {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, s))
		}
	}

	if len(p.FilesToInspect) > 0 {
		sb.WriteString("\nFiles to inspect:\n")
		for _, f := range p.FilesToInspect {
			sb.WriteString(fmt.Sprintf("  - %s\n", f))
		}
	}

	return sb.String()
}

// ToJSON converts the plan back to JSON string
func (p *Plan) ToJSON() (string, error) {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
