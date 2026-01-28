package orchestrator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/yourusername/playground/internal/system"
)

// PlannerOutput represents the structured output from the planner model
type PlannerOutput struct {
	Goal           string   `json:"goal"`
	Constraints    []string `json:"constraints"`
	Steps          []string `json:"steps"`
	FilesToInspect []string `json:"files_to_inspect"`
}

// ExecutorOutput represents the output from the executor model
type ExecutorOutput struct {
	Response  string
	ToolCalls []interface{} // Will be defined by agent package
}

// Execute runs the two-stage pipeline: planner → executor
func (o *Orchestrator) Execute(userInput string) (string, error) {
	// Validate models can fit in RAM
	if err := o.ValidateModels(); err != nil {
		return "", err
	}

	// Check llama-cli is available
	if err := o.CheckLlamaCppInstalled(); err != nil {
		return "", err
	}

	// Stage 1: Planning
	// Force GC before loading planner
	system.ForceGC()

	plan, err := o.runPlanner(userInput)
	if err != nil {
		return "", fmt.Errorf("planner failed: %w", err)
	}

	// Unload planner (GC will clean up)
	system.ForceGC()

	// Stage 2: Execution
	result, err := o.runExecutor(plan, userInput)
	if err != nil {
		return "", fmt.Errorf("executor failed: %w", err)
	}

	// Unload executor
	system.ForceGC()

	return result, nil
}

// runPlanner executes the planner model
func (o *Orchestrator) runPlanner(userInput string) (string, error) {
	if o.planner.LocalPath == "" {
		return "", fmt.Errorf("planner model not downloaded: %s", o.planner.Name)
	}

	// Check if model file exists
	if _, err := os.Stat(o.planner.LocalPath); err != nil {
		return "", fmt.Errorf("planner model file not found: %s", o.planner.LocalPath)
	}

	// Build planner prompt
	prompt := buildPlannerPrompt(userInput)

	// Run llama-cli with planner model
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "llama-cli",
		"--model", o.planner.LocalPath,
		"--prompt", prompt,
		"--ctx-size", "2048",
		"--n-predict", "512",
		"--temp", "0.1",
		"--top-k", "40",
		"--top-p", "0.9",
		"--threads", fmt.Sprintf("%d", o.systemInfo.CPUCores),
		"--no-display-prompt",
		"--log-disable",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("llama-cli failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// runExecutor executes the executor model
func (o *Orchestrator) runExecutor(plan, userInput string) (string, error) {
	if o.executor.LocalPath == "" {
		return "", fmt.Errorf("executor model not downloaded: %s", o.executor.Name)
	}

	// Check if model file exists
	if _, err := os.Stat(o.executor.LocalPath); err != nil {
		return "", fmt.Errorf("executor model file not found: %s", o.executor.LocalPath)
	}

	// Build executor prompt
	prompt := buildExecutorPrompt(plan, userInput)

	// Run llama-cli with executor model
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "llama-cli",
		"--model", o.executor.LocalPath,
		"--prompt", prompt,
		"--ctx-size", "4096",
		"--n-predict", "2048",
		"--temp", "0.1",
		"--top-k", "40",
		"--top-p", "0.9",
		"--threads", fmt.Sprintf("%d", o.systemInfo.CPUCores),
		"--no-display-prompt",
		"--log-disable",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("llama-cli failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// buildPlannerPrompt creates the prompt for the planner model
func buildPlannerPrompt(userInput string) string {
	return fmt.Sprintf(`You are a task planner. Convert user intent into structured JSON.

Output ONLY valid JSON in this format:
{
  "goal": "string",
  "constraints": ["string"],
  "steps": ["string"],
  "files_to_inspect": ["string"]
}

Rules:
- No markdown
- No prose
- No code
- JSON only

User request: %s

JSON output:`, userInput)
}

// buildExecutorPrompt creates the prompt for the executor model
func buildExecutorPrompt(plan, userInput string) string {
	return fmt.Sprintf(`You are a coding assistant. Generate code changes based on the plan.

Plan:
%s

Original request:
%s

Generate the necessary code changes. Use JSON tool calls for file operations.`, plan, userInput)
}
