package system

import (
	"runtime"
	"testing"
)

func TestDetectSystem(t *testing.T) {
	info, err := DetectSystem()
	if err != nil {
		t.Fatalf("DetectSystem() failed: %v", err)
	}

	// Validate RAM detection
	if info.TotalRAM == 0 {
		t.Error("TotalRAM should not be zero")
	}
	if info.AvailableRAM == 0 {
		t.Error("AvailableRAM should not be zero")
	}
	if info.AvailableRAM > info.TotalRAM {
		t.Errorf("AvailableRAM (%d) should not exceed TotalRAM (%d)", info.AvailableRAM, info.TotalRAM)
	}

	// Validate CPU detection
	if info.CPUCores == 0 {
		t.Error("CPUCores should not be zero")
	}
	if info.CPUCores != runtime.NumCPU() {
		t.Errorf("CPUCores (%d) should match runtime.NumCPU() (%d)", info.CPUCores, runtime.NumCPU())
	}

	// Validate OS detection
	if info.OS == "" {
		t.Error("OS should not be empty")
	}
	if info.OS != runtime.GOOS {
		t.Errorf("OS (%s) should match runtime.GOOS (%s)", info.OS, runtime.GOOS)
	}

	t.Logf("System Info: %s", info.String())
}

func TestGetUsableRAM(t *testing.T) {
	info := &SystemInfo{
		TotalRAM:     uint64(16) * 1024 * 1024 * 1024, // 16 GB
		AvailableRAM: uint64(8) * 1024 * 1024 * 1024,  // 8 GB
	}

	usableRAM := info.GetUsableRAM()
	// 70% of 8GB using integer arithmetic to avoid float conversion issues
	expectedRAM := (uint64(8) * 1024 * 1024 * 1024 * 7) / 10

	if usableRAM != expectedRAM {
		t.Errorf("GetUsableRAM() = %d, want %d", usableRAM, expectedRAM)
	}

	// Should be approximately 5.6 GB
	usableGB := float64(usableRAM) / (1024 * 1024 * 1024)
	if usableGB < 5.5 || usableGB > 5.7 {
		t.Errorf("Usable RAM should be ~5.6 GB, got %.2f GB", usableGB)
	}
}

func TestCanFitModel(t *testing.T) {
	info := &SystemInfo{
		TotalRAM:     uint64(16) * 1024 * 1024 * 1024, // 16 GB
		AvailableRAM: uint64(8) * 1024 * 1024 * 1024,  // 8 GB
	}

	tests := []struct {
		name        string
		modelSizeGB float64
		shouldFit   bool
	}{
		{"Small model (1GB)", 1.0, true},
		{"Medium model (3GB)", 3.0, true},
		{"Large model (5GB)", 5.0, true},
		{"Too large (6GB)", 6.0, false},
		{"Way too large (10GB)", 10.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelSizeBytes := uint64(tt.modelSizeGB * 1024 * 1024 * 1024)
			canFit := info.CanFitModel(modelSizeBytes)
			if canFit != tt.shouldFit {
				t.Errorf("CanFitModel(%s) = %v, want %v", tt.name, canFit, tt.shouldFit)
			}
		})
	}
}

func TestSystemInfoString(t *testing.T) {
	info := &SystemInfo{
		TotalRAM:     uint64(16) * 1024 * 1024 * 1024, // 16 GB
		AvailableRAM: uint64(8) * 1024 * 1024 * 1024,  // 8 GB
		CPUCores:     8,
		HasGPU:       true,
		OS:           "linux",
	}

	str := info.String()
	if str == "" {
		t.Error("String() should not be empty")
	}

	// Should contain key information
	if !contains(str, "linux") {
		t.Error("String() should contain OS")
	}
	if !contains(str, "8 cores") {
		t.Error("String() should contain CPU cores")
	}
	if !contains(str, "true") {
		t.Error("String() should contain GPU status")
	}

	t.Logf("System Info String: %s", str)
}

func TestForceGC(t *testing.T) {
	// This test just ensures ForceGC doesn't panic
	// We can't easily verify memory was actually freed
	ForceGC()
	t.Log("ForceGC() completed without panic")
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
