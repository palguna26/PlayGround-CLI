package orchestrator

import (
	"testing"

	"github.com/yourusername/playground/internal/model"
)

func TestNew(t *testing.T) {
	orch, err := New()
	if err != nil {
		// This might fail if system doesn't have enough RAM
		t.Logf("New() failed (expected on low-RAM systems): %v", err)
		return
	}

	if orch == nil {
		t.Fatal("Expected orchestrator, got nil")
	}

	if orch.systemInfo == nil {
		t.Error("System info should not be nil")
	}

	if orch.planner == nil {
		t.Error("Planner should not be nil")
	}

	if orch.executor == nil {
		t.Error("Executor should not be nil")
	}

	t.Logf("Orchestrator created:\n%s", orch.String())
}

func TestNewWithModels(t *testing.T) {
	planner := &model.ModelSpec{
		Name:   "Test Planner",
		Role:   model.RolePlanner,
		MinRAM: uint64(2) * model.GB,
		SizeGB: 0.5,
	}

	executor := &model.ModelSpec{
		Name:   "Test Executor",
		Role:   model.RoleExecutor,
		MinRAM: uint64(4) * model.GB,
		SizeGB: 1.0,
	}

	orch, err := NewWithModels(planner, executor)
	if err != nil {
		t.Fatalf("NewWithModels() failed: %v", err)
	}

	if orch.planner.Name != "Test Planner" {
		t.Errorf("Expected planner 'Test Planner', got %s", orch.planner.Name)
	}

	if orch.executor.Name != "Test Executor" {
		t.Errorf("Expected executor 'Test Executor', got %s", orch.executor.Name)
	}
}

func TestGetters(t *testing.T) {
	planner := &model.ModelSpec{
		Name:   "Test Planner",
		Role:   model.RolePlanner,
		MinRAM: uint64(2) * model.GB,
		SizeGB: 0.5,
	}

	executor := &model.ModelSpec{
		Name:   "Test Executor",
		Role:   model.RoleExecutor,
		MinRAM: uint64(4) * model.GB,
		SizeGB: 1.0,
	}

	orch, err := NewWithModels(planner, executor)
	if err != nil {
		t.Fatalf("NewWithModels() failed: %v", err)
	}

	// Test GetSystemInfo
	sysInfo := orch.GetSystemInfo()
	if sysInfo == nil {
		t.Error("GetSystemInfo() should not return nil")
	}

	// Test GetPlanner
	p := orch.GetPlanner()
	if p == nil || p.Name != "Test Planner" {
		t.Error("GetPlanner() returned incorrect planner")
	}

	// Test GetExecutor
	e := orch.GetExecutor()
	if e == nil || e.Name != "Test Executor" {
		t.Error("GetExecutor() returned incorrect executor")
	}
}

func TestValidateModels(t *testing.T) {
	tests := []struct {
		name           string
		plannerSizeGB  float64
		executorSizeGB float64
		shouldFail     bool
		skipOnLowRAM   bool
	}{
		{
			name:           "Tiny models should fit",
			plannerSizeGB:  0.1,
			executorSizeGB: 0.2,
			shouldFail:     false,
			skipOnLowRAM:   false,
		},
		{
			name:           "Huge planner should fail",
			plannerSizeGB:  100.0,
			executorSizeGB: 0.2,
			shouldFail:     true,
			skipOnLowRAM:   false,
		},
		{
			name:           "Huge executor should fail",
			plannerSizeGB:  0.1,
			executorSizeGB: 100.0,
			shouldFail:     true,
			skipOnLowRAM:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planner := &model.ModelSpec{
				Name:   "Test Planner",
				Role:   model.RolePlanner,
				MinRAM: uint64(2) * model.GB,
				SizeGB: tt.plannerSizeGB,
			}

			executor := &model.ModelSpec{
				Name:   "Test Executor",
				Role:   model.RoleExecutor,
				MinRAM: uint64(4) * model.GB,
				SizeGB: tt.executorSizeGB,
			}

			orch, err := NewWithModels(planner, executor)
			if err != nil {
				t.Fatalf("NewWithModels() failed: %v", err)
			}

			// Skip validation test on low-RAM systems for non-huge models
			if tt.skipOnLowRAM && orch.systemInfo.AvailableRAM < uint64(2)*model.GB {
				t.Skip("Skipping on low-RAM system")
			}

			err = orch.ValidateModels()
			if tt.shouldFail {
				if err == nil {
					t.Error("Expected validation to fail, but it passed")
				} else {
					t.Logf("Expected error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected validation to pass, but got error: %v", err)
				}
			}
		})
	}
}

func TestModelPaths(t *testing.T) {
	planner := &model.ModelSpec{
		Name:   "Test Planner",
		Role:   model.RolePlanner,
		MinRAM: uint64(2) * model.GB,
		SizeGB: 0.5,
	}

	executor := &model.ModelSpec{
		Name:   "Test Executor",
		Role:   model.RoleExecutor,
		MinRAM: uint64(4) * model.GB,
		SizeGB: 1.0,
	}

	orch, err := NewWithModels(planner, executor)
	if err != nil {
		t.Fatalf("NewWithModels() failed: %v", err)
	}

	// Test initial state
	plannerPath, executorPath := orch.GetModelPaths()
	if plannerPath != "" || executorPath != "" {
		t.Error("Initial paths should be empty")
	}

	// Test NeedsDownload
	plannerNeeded, executorNeeded := orch.NeedsDownload()
	if !plannerNeeded || !executorNeeded {
		t.Error("Both models should need download initially")
	}

	// Test SetModelPaths
	orch.SetModelPaths("/path/to/planner.gguf", "/path/to/executor.gguf")
	plannerPath, executorPath = orch.GetModelPaths()
	if plannerPath != "/path/to/planner.gguf" {
		t.Errorf("Expected planner path '/path/to/planner.gguf', got %s", plannerPath)
	}
	if executorPath != "/path/to/executor.gguf" {
		t.Errorf("Expected executor path '/path/to/executor.gguf', got %s", executorPath)
	}

	// Test NeedsDownload after setting paths
	plannerNeeded, executorNeeded = orch.NeedsDownload()
	if plannerNeeded || executorNeeded {
		t.Error("No models should need download after setting paths")
	}
}

func TestString(t *testing.T) {
	planner := &model.ModelSpec{
		Name:         "Test Planner",
		Role:         model.RolePlanner,
		MinRAM:       uint64(2) * model.GB,
		SizeGB:       0.5,
		Quantization: "Q4_K_M",
	}

	executor := &model.ModelSpec{
		Name:         "Test Executor",
		Role:         model.RoleExecutor,
		MinRAM:       uint64(4) * model.GB,
		SizeGB:       1.0,
		Quantization: "Q4_K_M",
	}

	orch, err := NewWithModels(planner, executor)
	if err != nil {
		t.Fatalf("NewWithModels() failed: %v", err)
	}

	str := orch.String()
	if str == "" {
		t.Error("String() should not be empty")
	}

	// Should contain key information
	if !contains(str, "Test Planner") {
		t.Error("String() should contain planner name")
	}
	if !contains(str, "Test Executor") {
		t.Error("String() should contain executor name")
	}

	t.Logf("Orchestrator string:\n%s", str)
}

func TestCheckLlamaCppInstalled(t *testing.T) {
	planner := &model.ModelSpec{
		Name:   "Test Planner",
		Role:   model.RolePlanner,
		MinRAM: uint64(2) * model.GB,
		SizeGB: 0.5,
	}

	executor := &model.ModelSpec{
		Name:   "Test Executor",
		Role:   model.RoleExecutor,
		MinRAM: uint64(4) * model.GB,
		SizeGB: 1.0,
	}

	orch, err := NewWithModels(planner, executor)
	if err != nil {
		t.Fatalf("NewWithModels() failed: %v", err)
	}

	err = orch.CheckLlamaCppInstalled()
	if err != nil {
		t.Logf("llama-cli not installed (expected in CI): %v", err)
	} else {
		t.Log("llama-cli is installed")
	}
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
