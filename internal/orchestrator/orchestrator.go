package orchestrator

import (
	"fmt"
	"os/exec"

	"github.com/yourusername/playground/internal/model"
	"github.com/yourusername/playground/internal/system"
)

// Orchestrator manages the two-stage pipeline (planner → executor)
type Orchestrator struct {
	systemInfo *system.SystemInfo
	planner    *model.ModelSpec
	executor   *model.ModelSpec
}

// New creates a new orchestrator with automatic model selection
func New() (*Orchestrator, error) {
	// Detect system capabilities
	sysInfo, err := system.DetectSystem()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system: %w", err)
	}

	// Select best models for this system
	planner, executor, err := model.SelectModels(sysInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to select models: %w", err)
	}

	return &Orchestrator{
		systemInfo: sysInfo,
		planner:    planner,
		executor:   executor,
	}, nil
}

// NewWithModels creates an orchestrator with specific models (for testing)
func NewWithModels(planner, executor *model.ModelSpec) (*Orchestrator, error) {
	sysInfo, err := system.DetectSystem()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system: %w", err)
	}

	return &Orchestrator{
		systemInfo: sysInfo,
		planner:    planner,
		executor:   executor,
	}, nil
}

// GetSystemInfo returns the detected system information
func (o *Orchestrator) GetSystemInfo() *system.SystemInfo {
	return o.systemInfo
}

// GetPlanner returns the selected planner model
func (o *Orchestrator) GetPlanner() *model.ModelSpec {
	return o.planner
}

// GetExecutor returns the selected executor model
func (o *Orchestrator) GetExecutor() *model.ModelSpec {
	return o.executor
}

// ValidateModels checks if selected models can fit in available RAM
func (o *Orchestrator) ValidateModels() error {
	usableRAM := o.systemInfo.GetUsableRAM()

	// Check planner
	plannerSizeBytes := uint64(o.planner.SizeGB * 1024 * 1024 * 1024)
	if !o.systemInfo.CanFitModel(plannerSizeBytes) {
		return fmt.Errorf(
			"planner model %s (%.1f GB) exceeds usable RAM (%.1f GB)",
			o.planner.Name,
			o.planner.SizeGB,
			float64(usableRAM)/(1024*1024*1024),
		)
	}

	// Check executor
	executorSizeBytes := uint64(o.executor.SizeGB * 1024 * 1024 * 1024)
	if !o.systemInfo.CanFitModel(executorSizeBytes) {
		return fmt.Errorf(
			"executor model %s (%.1f GB) exceeds usable RAM (%.1f GB)",
			o.executor.Name,
			o.executor.SizeGB,
			float64(usableRAM)/(1024*1024*1024),
		)
	}

	return nil
}

// CheckLlamaCppInstalled verifies that llama-cli is available
func (o *Orchestrator) CheckLlamaCppInstalled() error {
	_, err := exec.LookPath("llama-cli")
	if err != nil {
		return fmt.Errorf("llama-cli not found in PATH: %w", err)
	}
	return nil
}

// String returns a human-readable summary of the orchestrator configuration
func (o *Orchestrator) String() string {
	return fmt.Sprintf(
		"Orchestrator:\n  System: %s\n  Planner: %s\n  Executor: %s",
		o.systemInfo.String(),
		o.planner.String(),
		o.executor.String(),
	)
}

// GetModelPaths returns the local paths for planner and executor models
func (o *Orchestrator) GetModelPaths() (plannerPath, executorPath string) {
	return o.planner.LocalPath, o.executor.LocalPath
}

// SetModelPaths sets the local paths for planner and executor models
func (o *Orchestrator) SetModelPaths(plannerPath, executorPath string) {
	o.planner.LocalPath = plannerPath
	o.executor.LocalPath = executorPath
}

// NeedsDownload checks if any models need to be downloaded
func (o *Orchestrator) NeedsDownload() (plannerNeeded, executorNeeded bool) {
	plannerNeeded = o.planner.LocalPath == ""
	executorNeeded = o.executor.LocalPath == ""
	return
}
