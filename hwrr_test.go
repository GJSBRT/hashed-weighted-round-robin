package hwrr

import (
	"fmt"
	"testing"
	"math/rand/v2"
)

// TestNewHWRR tests the NewHWRR function.
func TestNewHWRR(t *testing.T) {
	backends := []*Backend{
		{"backend1", 1},
		{"backend2", 2},
		{"backend3", 3},
	}

	hwrr := NewHWRR(backends)

	if len(hwrr.backends) != 6 {
		t.Errorf("Expected 6 backends, got %d", len(hwrr.backends))
	}
}

// TestGetNextBackend tests if the GetNextBackend function returns the correct backend based on the hash given.
func TestGetNextBackend(t *testing.T) {
	backends := []*Backend{
		{"backend1", 1},
		{"backend2", 2},
		{"backend3", 3},
	}

	hwrr := NewHWRR(backends)

	backend := hwrr.GetNextBackend(0)

	if backend.Name != "backend1" {
		t.Errorf("Expected backend1, got %s", backend.Name)
	}

	backend = hwrr.GetNextBackend(1)

	if backend.Name != "backend2" {
		t.Errorf("Expected backend2, got %s", backend.Name)
	}

	backend = hwrr.GetNextBackend(2)

	if backend.Name != "backend2" {
		t.Errorf("Expected backend2, got %s", backend.Name)
	}

	backend = hwrr.GetNextBackend(3)

	if backend.Name != "backend3" {
		t.Errorf("Expected backend3, got %s", backend.Name)
	}

	backend = hwrr.GetNextBackend(4)

	if backend.Name != "backend3" {
		t.Errorf("Expected backend3, got %s", backend.Name)
	}

	backend = hwrr.GetNextBackend(5)

	if backend.Name != "backend3" {
		t.Errorf("Expected backend3, got %s", backend.Name)
	}
}

// TestGetNextBackendWithSingleBackend tests if the GetNextBackend function returns the correct backend when there is only one backend.
func TestGetNextBackendWithSingleBackend(t *testing.T) {
	backends := []*Backend{
		{"backend1", 1},
	}

	hwrr := NewHWRR(backends)

	backend := hwrr.GetNextBackend(0)

	if backend.Name != "backend1" {
		t.Errorf("Expected backend1, got %s", backend.Name)
	}
}

// TestGetNextBackendWithZeroWeight tests if the GetNextBackend function returns nil when there is a backend with zero weight.
func TestGetNextBackendWithZeroWeight(t *testing.T) {
	backends := []*Backend{
		{"backend1", 0},
	}

	hwrr := NewHWRR(backends)

	backend := hwrr.GetNextBackend(0)

	if backend != nil {
		t.Errorf("Expected nil, got %s", backend.Name)
	}
}

// BenchmarkGetNextBackend benchmarks the GetNextBackend function.
func BenchmarkGetNextBackend(b *testing.B) {
	backends := []*Backend{
		{"backend1", 1},
		{"backend2", 2},
		{"backend3", 3},
	}

	hwrr := NewHWRR(backends)

	for i := 0; i < b.N; i++ {
		hwrr.GetNextBackend(i)
	}
}

// BenchmarkGetNextBackendWithSingleBackend benchmarks the GetNextBackend function with a single backend.
func BenchmarkGetNextBackendWithSingleBackend(b *testing.B) {
	backends := []*Backend{
		{"backend1", 1},
	}

	hwrr := NewHWRR(backends)

	for i := 0; i < b.N; i++ {
		hwrr.GetNextBackend(i)
	}
}

// BenchmarkGetNextBackendWithZeroWeight benchmarks the GetNextBackend function with a backend with zero weight.
func BenchmarkGetNextBackendWithZeroWeight(b *testing.B) {
	backends := []*Backend{
		{"backend1", 0},
	}

	hwrr := NewHWRR(backends)

	for i := 0; i < b.N; i++ {
		hwrr.GetNextBackend(i)
	}
}

// BenchmarkAllocationTestEvenWeight "benchmarks" the distribution of backends with even weight.
func BenchmarkAllocationTestEvenWeight(b *testing.B) {
	backends := []*Backend{
		{"backend1", 2},
		{"backend2", 2},
		{"backend3", 2},
	}

	hwrr := NewHWRR(backends)

	backendDistribution := make(map[string]int)

	for i := 0; i < b.N; i++ {
		sourcePort := rand.IntN(65535-1) + 1 
		destPort := rand.IntN(65535-1) + 1

		backend := hwrr.GetNextBackend(hash(fmt.Sprintf("192.168.1.1:%d:10.10.0.1:%d", sourcePort, destPort)))
		if backend == nil {
			b.Fatal("Expected backend, got nil")
		}

		backendDistribution[backend.Name]++
	}

	backendDistributionPercentage := make(map[string]float64)
	for name, count := range backendDistribution {
		backendDistributionPercentage[name] = float64(count) / float64(b.N) * 100
	}

	// Check if distribution is within 3% of each other
	if b.N > 100 {
		for _, count := range backendDistribution {
			if float64(count) < float64(b.N)/3*0.97 || float64(count) > float64(b.N)/3*1.03 {
				b.Fatalf("Expected 33.33%% distribution, got %v", backendDistributionPercentage)
			}
		}
	}

	for name, percentage := range backendDistributionPercentage {
		b.ReportMetric(percentage, name)
	}
}

// BenchmarkAllocationTestTopHeavyWeight "benchmarks" the distribution of backends with a top-heavy weighted backend.
func BenchmarkAllocationTestTopHeavyWeight(b *testing.B) {
	backends := []*Backend{
		{"backend1", 3},
		{"backend2", 1},
		{"backend3", 1},
	}

	hwrr := NewHWRR(backends)

	backendDistribution := make(map[string]int)

	for i := 0; i < b.N; i++ {
		sourcePort := rand.IntN(65535-1) + 1 
		destPort := rand.IntN(65535-1) + 1

		backend := hwrr.GetNextBackend(hash(fmt.Sprintf("192.168.1.1:%d:10.10.0.1:%d", sourcePort, destPort)))
		if backend == nil {
			b.Fatal("Expected backend, got nil")
		}

		backendDistribution[backend.Name]++
	}

	backendDistributionPercentage := make(map[string]float64)
	for name, count := range backendDistribution {
		backendDistributionPercentage[name] = float64(count) / float64(b.N) * 100
	}

	// Check if backend1 has 60% distribution
	if b.N > 100 {
		if backendDistributionPercentage["backend1"] < 60*0.97 || backendDistributionPercentage["backend1"] > 60*1.03 {
			b.Fatalf("Expected 60%% distribution, got %v", backendDistributionPercentage)
		}
	}

	for name, percentage := range backendDistributionPercentage {
		b.ReportMetric(percentage, name)
	}
}
