package model

import (
	"fmt"
	"sort"

	"github.com/yourusername/playground/internal/system"
)

// ModelRole defines the role of a model in the orchestration pipeline
type ModelRole string

const (
	RolePlanner  ModelRole = "planner"
	RoleExecutor ModelRole = "executor"
)

// ModelSpec defines a model's specifications and requirements
type ModelSpec struct {
	Name          string    // Human-readable name
	Role          ModelRole // planner or executor
	MinRAM        uint64    // Minimum RAM required in bytes
	SizeGB        float64   // Approximate model size in GB
	Quantization  string    // Quantization format (e.g., Q4_K_M)
	HuggingFaceID string    // HuggingFace repository ID
	FileName      string    // GGUF filename
	LocalPath     string    // Local path (set after download)
}

// GB constant for readability
const GB = 1024 * 1024 * 1024

// ModelRegistry contains all available models
var ModelRegistry = []ModelSpec{
	// ===== PLANNERS =====
	// Tiny, fast models for task planning
	{
		Name:          "Qwen2.5-0.5B-Instruct",
		Role:          RolePlanner,
		MinRAM:        uint64(3) * GB,
		SizeGB:        0.5,
		Quantization:  "Q4_K_M",
		HuggingFaceID: "Qwen/Qwen2.5-0.5B-Instruct-GGUF",
		FileName:      "qwen2.5-0.5b-instruct-q4_k_m.gguf",
	},
	{
		Name:          "TinyLLaMA-1.1B",
		Role:          RolePlanner,
		MinRAM:        uint64(2) * GB,
		SizeGB:        0.7,
		Quantization:  "Q4_K_M",
		HuggingFaceID: "TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF",
		FileName:      "tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf",
	},

	// ===== EXECUTORS =====
	// Larger, more capable models for code generation
	{
		Name:          "DeepSeek-Coder-6.7B-Instruct",
		Role:          RoleExecutor,
		MinRAM:        uint64(10) * GB,
		SizeGB:        4.1,
		Quantization:  "Q4_K_M",
		HuggingFaceID: "QuantFactory/deepseek-coder-6.7b-instruct-GGUF",
		FileName:      "deepseek-coder-6.7b-instruct.Q4_K_M.gguf",
	},
	{
		Name:          "Qwen2.5-Coder-3B-Instruct",
		Role:          RoleExecutor,
		MinRAM:        uint64(7) * GB,
		SizeGB:        2.0,
		Quantization:  "Q4_K_M",
		HuggingFaceID: "Qwen/Qwen2.5-Coder-3B-Instruct-GGUF",
		FileName:      "qwen2.5-coder-3b-instruct-q4_k_m.gguf",
	},
	{
		Name:          "Qwen2.5-Coder-1.5B-Instruct",
		Role:          RoleExecutor,
		MinRAM:        uint64(4) * GB,
		SizeGB:        1.1,
		Quantization:  "Q4_K_M",
		HuggingFaceID: "Qwen/Qwen2.5-Coder-1.5B-Instruct-GGUF",
		FileName:      "qwen2.5-coder-1.5b-instruct-q4_k_m.gguf",
	},
	{
		Name:          "StarCoder-Base-1B",
		Role:          RoleExecutor,
		MinRAM:        uint64(3) * GB,
		SizeGB:        0.7,
		Quantization:  "Q4_K_M",
		HuggingFaceID: "bigcode/starcoderbase-1b",
		FileName:      "starcoderbase-1b.Q4_K_M.gguf",
	},
}

// SelectPlanner selects the best planner model based on available RAM
func SelectPlanner(availableRAM uint64) *ModelSpec {
	// Filter planners that fit in available RAM
	var candidates []ModelSpec
	for _, model := range ModelRegistry {
		if model.Role == RolePlanner && model.MinRAM <= availableRAM {
			candidates = append(candidates, model)
		}
	}

	if len(candidates) == 0 {
		return nil // No planner fits
	}

	// Sort by MinRAM (descending) - prefer higher-quality models if they fit
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].MinRAM > candidates[j].MinRAM
	})

	return &candidates[0]
}

// SelectExecutor selects the best executor model based on available RAM
func SelectExecutor(availableRAM uint64) *ModelSpec {
	// Filter executors that fit in available RAM
	var candidates []ModelSpec
	for _, model := range ModelRegistry {
		if model.Role == RoleExecutor && model.MinRAM <= availableRAM {
			candidates = append(candidates, model)
		}
	}

	if len(candidates) == 0 {
		return nil // No executor fits
	}

	// Sort by MinRAM (descending) - prefer higher-quality models if they fit
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].MinRAM > candidates[j].MinRAM
	})

	return &candidates[0]
}

// SelectModels selects both planner and executor based on system info
func SelectModels(sysInfo *system.SystemInfo) (*ModelSpec, *ModelSpec, error) {
	usableRAM := sysInfo.GetUsableRAM()

	planner := SelectPlanner(usableRAM)
	if planner == nil {
		return nil, nil, fmt.Errorf(
			"no planner model fits in available RAM (%.1f GB available, minimum %.1f GB required)",
			float64(usableRAM)/GB,
			float64(ModelRegistry[0].MinRAM)/GB,
		)
	}

	executor := SelectExecutor(usableRAM)
	if executor == nil {
		return nil, nil, fmt.Errorf(
			"no executor model fits in available RAM (%.1f GB available, minimum %.1f GB required)",
			float64(usableRAM)/GB,
			float64(ModelRegistry[len(ModelRegistry)-1].MinRAM)/GB,
		)
	}

	return planner, executor, nil
}

// GetModelByName retrieves a model spec by name
func GetModelByName(name string) *ModelSpec {
	for _, model := range ModelRegistry {
		if model.Name == name {
			return &model
		}
	}
	return nil
}

// GetModelsByRole retrieves all models with a specific role
func GetModelsByRole(role ModelRole) []ModelSpec {
	var models []ModelSpec
	for _, model := range ModelRegistry {
		if model.Role == role {
			models = append(models, model)
		}
	}
	return models
}

// String returns a human-readable representation of the model spec
func (m *ModelSpec) String() string {
	return fmt.Sprintf(
		"%s (%s, %.1f GB, %s, min RAM: %.1f GB)",
		m.Name,
		m.Role,
		m.SizeGB,
		m.Quantization,
		float64(m.MinRAM)/GB,
	)
}

// GetDownloadURL returns the HuggingFace download URL for this model
func (m *ModelSpec) GetDownloadURL() string {
	return fmt.Sprintf(
		"https://huggingface.co/%s/resolve/main/%s",
		m.HuggingFaceID,
		m.FileName,
	)
}
