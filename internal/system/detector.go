package system

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/shirou/gopsutil/v3/mem"
)

// SystemInfo contains hardware and OS information
type SystemInfo struct {
	TotalRAM     uint64 // Total system RAM in bytes
	AvailableRAM uint64 // Available RAM in bytes
	CPUCores     int    // Number of CPU cores
	HasGPU       bool   // GPU detected (read-only, not used for selection)
	OS           string // Operating system (windows, linux, darwin)
}

// DetectSystem detects current system capabilities
func DetectSystem() (*SystemInfo, error) {
	info := &SystemInfo{
		OS: runtime.GOOS,
	}

	// Detect RAM
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to detect RAM: %w", err)
	}
	info.TotalRAM = vmStat.Total
	info.AvailableRAM = vmStat.Available

	// Detect CPU cores
	info.CPUCores = runtime.NumCPU()

	// Detect GPU (optional, read-only)
	info.HasGPU = detectGPU()

	return info, nil
}

// detectGPU attempts to detect GPU presence (read-only, not used for model selection)
func detectGPU() bool {
	// This is a simple heuristic - check for common GPU indicators
	// We don't use this for model selection, just for informational purposes

	switch runtime.GOOS {
	case "windows":
		// On Windows, check for CUDA or DirectML
		return checkWindowsGPU()
	case "darwin":
		// On macOS, Metal is always available on modern Macs
		return checkMacGPU()
	case "linux":
		// On Linux, check for CUDA or ROCm
		return checkLinuxGPU()
	default:
		return false
	}
}

// checkWindowsGPU checks for GPU on Windows
func checkWindowsGPU() bool {
	// Simple check: if we have NVIDIA or AMD drivers, assume GPU exists
	// This is read-only and doesn't affect model selection
	// TODO: Could use wmi or registry checks for more accuracy
	return false // Conservative default
}

// checkMacGPU checks for GPU on macOS
func checkMacGPU() bool {
	// Metal is available on all modern Macs (2012+)
	// We could check sysctl for more details, but this is sufficient
	return true // Most Macs have Metal support
}

// checkLinuxGPU checks for GPU on Linux
func checkLinuxGPU() bool {
	// Check for NVIDIA (CUDA) or AMD (ROCm)
	// This is a conservative check - we don't use it for model selection
	// TODO: Could check /proc/driver/nvidia/version or lspci
	return false // Conservative default
}

// GetUsableRAM returns the amount of RAM we can safely use (70% of available)
func (s *SystemInfo) GetUsableRAM() uint64 {
	return uint64(float64(s.AvailableRAM) * 0.7)
}

// CanFitModel checks if a model of given size can fit in available RAM
func (s *SystemInfo) CanFitModel(modelSizeBytes uint64) bool {
	usableRAM := s.GetUsableRAM()
	return modelSizeBytes <= usableRAM
}

// String returns a human-readable representation of system info
func (s *SystemInfo) String() string {
	return fmt.Sprintf(
		"OS: %s | RAM: %.1f GB total, %.1f GB available | CPU: %d cores | GPU: %v",
		s.OS,
		float64(s.TotalRAM)/(1024*1024*1024),
		float64(s.AvailableRAM)/(1024*1024*1024),
		s.CPUCores,
		s.HasGPU,
	)
}

// ForceGC forces garbage collection to free up memory before model loading
func ForceGC() {
	runtime.GC()
	debug.FreeOSMemory()
}
