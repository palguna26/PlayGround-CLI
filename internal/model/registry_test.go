package model

import (
	"testing"

	"github.com/yourusername/playground/internal/system"
)

func TestSelectPlanner(t *testing.T) {
	tests := []struct {
		name          string
		availableRAM  uint64
		expectedModel string
		shouldBeNil   bool
	}{
		{
			name:          "High RAM (16GB) - should select Qwen2.5-0.5B",
			availableRAM:  uint64(16) * GB,
			expectedModel: "Qwen2.5-0.5B-Instruct",
			shouldBeNil:   false,
		},
		{
			name:          "Medium RAM (4GB) - should select Qwen2.5-0.5B",
			availableRAM:  uint64(4) * GB,
			expectedModel: "Qwen2.5-0.5B-Instruct",
			shouldBeNil:   false,
		},
		{
			name:          "Low RAM (2.5GB) - should select TinyLLaMA",
			availableRAM:  uint64(2500) * 1024 * 1024,
			expectedModel: "TinyLLaMA-1.1B",
			shouldBeNil:   false,
		},
		{
			name:         "Very Low RAM (1GB) - should return nil",
			availableRAM: uint64(1) * GB,
			shouldBeNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planner := SelectPlanner(tt.availableRAM)
			if tt.shouldBeNil {
				if planner != nil {
					t.Errorf("Expected nil, got %s", planner.Name)
				}
			} else {
				if planner == nil {
					t.Fatalf("Expected model, got nil")
				}
				if planner.Name != tt.expectedModel {
					t.Errorf("Expected %s, got %s", tt.expectedModel, planner.Name)
				}
				t.Logf("Selected: %s", planner.String())
			}
		})
	}
}

func TestSelectExecutor(t *testing.T) {
	tests := []struct {
		name          string
		availableRAM  uint64
		expectedModel string
		shouldBeNil   bool
	}{
		{
			name:          "High RAM (16GB) - should select DeepSeek-Coder-6.7B",
			availableRAM:  uint64(16) * GB,
			expectedModel: "DeepSeek-Coder-6.7B-Instruct",
			shouldBeNil:   false,
		},
		{
			name:          "Medium RAM (8GB) - should select Qwen2.5-Coder-3B",
			availableRAM:  uint64(8) * GB,
			expectedModel: "Qwen2.5-Coder-3B-Instruct",
			shouldBeNil:   false,
		},
		{
			name:          "Low RAM (5GB) - should select Qwen2.5-Coder-1.5B",
			availableRAM:  uint64(5) * GB,
			expectedModel: "Qwen2.5-Coder-1.5B-Instruct",
			shouldBeNil:   false,
		},
		{
			name:          "Very Low RAM (3.5GB) - should select StarCoder-Base-1B",
			availableRAM:  uint64(3500) * 1024 * 1024,
			expectedModel: "StarCoder-Base-1B",
			shouldBeNil:   false,
		},
		{
			name:         "Insufficient RAM (2GB) - should return nil",
			availableRAM: uint64(2) * GB,
			shouldBeNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := SelectExecutor(tt.availableRAM)
			if tt.shouldBeNil {
				if executor != nil {
					t.Errorf("Expected nil, got %s", executor.Name)
				}
			} else {
				if executor == nil {
					t.Fatalf("Expected model, got nil")
				}
				if executor.Name != tt.expectedModel {
					t.Errorf("Expected %s, got %s", tt.expectedModel, executor.Name)
				}
				t.Logf("Selected: %s", executor.String())
			}
		})
	}
}

func TestSelectModels(t *testing.T) {
	tests := []struct {
		name             string
		totalRAM         uint64
		availableRAM     uint64
		expectError      bool
		expectedPlanner  string
		expectedExecutor string
	}{
		{
			name:             "High-end system (16GB total, 15GB available)",
			totalRAM:         uint64(16) * GB,
			availableRAM:     uint64(15) * GB,
			expectError:      false,
			expectedPlanner:  "Qwen2.5-0.5B-Instruct",
			expectedExecutor: "DeepSeek-Coder-6.7B-Instruct", // 70% of 15GB = 10.5GB, fits DeepSeek (10GB)
		},
		{
			name:             "Mid-range system (8GB total, 6GB available)",
			totalRAM:         uint64(8) * GB,
			availableRAM:     uint64(6) * GB,
			expectError:      false,
			expectedPlanner:  "Qwen2.5-0.5B-Instruct",
			expectedExecutor: "Qwen2.5-Coder-1.5B-Instruct", // 70% of 6GB = 4.2GB, fits 1.5B (4GB)
		},
		{
			name:             "Low-end system (8GB total, 5GB available)",
			totalRAM:         uint64(8) * GB,
			availableRAM:     uint64(5) * GB,
			expectError:      false,
			expectedPlanner:  "Qwen2.5-0.5B-Instruct", // 70% of 5GB = 3.5GB, fits Qwen (3GB min)
			expectedExecutor: "StarCoder-Base-1B",     // 70% of 5GB = 3.5GB, fits StarCoder (3GB min)
		},
		{
			name:         "Insufficient RAM (2GB total, 1GB available)",
			totalRAM:     uint64(2) * GB,
			availableRAM: uint64(1) * GB,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sysInfo := &system.SystemInfo{
				TotalRAM:     tt.totalRAM,
				AvailableRAM: tt.availableRAM,
				CPUCores:     4,
				OS:           "linux",
			}

			planner, executor, err := SelectModels(sysInfo)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				t.Logf("Expected error: %v", err)
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if planner.Name != tt.expectedPlanner {
					t.Errorf("Expected planner %s, got %s", tt.expectedPlanner, planner.Name)
				}
				if executor.Name != tt.expectedExecutor {
					t.Errorf("Expected executor %s, got %s", tt.expectedExecutor, executor.Name)
				}
				t.Logf("Selected planner: %s", planner.String())
				t.Logf("Selected executor: %s", executor.String())
			}
		})
	}
}

func TestGetModelByName(t *testing.T) {
	tests := []struct {
		name        string
		modelName   string
		shouldExist bool
	}{
		{"Existing planner", "Qwen2.5-0.5B-Instruct", true},
		{"Existing executor", "DeepSeek-Coder-6.7B-Instruct", true},
		{"Non-existent model", "GPT-4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := GetModelByName(tt.modelName)
			if tt.shouldExist {
				if model == nil {
					t.Errorf("Expected to find model %s, got nil", tt.modelName)
				} else if model.Name != tt.modelName {
					t.Errorf("Expected %s, got %s", tt.modelName, model.Name)
				}
			} else {
				if model != nil {
					t.Errorf("Expected nil for %s, got %s", tt.modelName, model.Name)
				}
			}
		})
	}
}

func TestGetModelsByRole(t *testing.T) {
	planners := GetModelsByRole(RolePlanner)
	if len(planners) == 0 {
		t.Error("Expected at least one planner model")
	}
	for _, model := range planners {
		if model.Role != RolePlanner {
			t.Errorf("Expected planner role, got %s", model.Role)
		}
	}
	t.Logf("Found %d planner models", len(planners))

	executors := GetModelsByRole(RoleExecutor)
	if len(executors) == 0 {
		t.Error("Expected at least one executor model")
	}
	for _, model := range executors {
		if model.Role != RoleExecutor {
			t.Errorf("Expected executor role, got %s", model.Role)
		}
	}
	t.Logf("Found %d executor models", len(executors))
}

func TestModelSpecString(t *testing.T) {
	model := &ModelSpec{
		Name:         "Test Model",
		Role:         RolePlanner,
		MinRAM:       uint64(4) * GB,
		SizeGB:       2.5,
		Quantization: "Q4_K_M",
	}

	str := model.String()
	if str == "" {
		t.Error("String() should not be empty")
	}
	t.Logf("Model string: %s", str)
}

func TestGetDownloadURL(t *testing.T) {
	model := &ModelSpec{
		HuggingFaceID: "Qwen/Qwen2.5-0.5B-Instruct-GGUF",
		FileName:      "qwen2.5-0.5b-instruct-q4_k_m.gguf",
	}

	url := model.GetDownloadURL()
	expectedURL := "https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct-GGUF/resolve/main/qwen2.5-0.5b-instruct-q4_k_m.gguf"
	if url != expectedURL {
		t.Errorf("Expected URL:\n%s\nGot:\n%s", expectedURL, url)
	}
	t.Logf("Download URL: %s", url)
}
