package calculator

import (
	"testing"
)

// BenchmarkAdd benchmarks the Add method
func BenchmarkAdd(b *testing.B) {
	calc := NewCalculator()
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	
	// Run the benchmark b.N times
	for i := 0; i < b.N; i++ {
		calc.Add(2.0, 3.0)
	}
}

// BenchmarkSubtract benchmarks the Subtract method
func BenchmarkSubtract(b *testing.B) {
	calc := NewCalculator()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		calc.Subtract(5.0, 3.0)
	}
}

// BenchmarkMultiply benchmarks the Multiply method
func BenchmarkMultiply(b *testing.B) {
	calc := NewCalculator()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		calc.Multiply(2.0, 3.0)
	}
}

// BenchmarkDivide benchmarks the Divide method
func BenchmarkDivide(b *testing.B) {
	calc := NewCalculator()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, _ = calc.Divide(6.0, 3.0)
	}
}

// BenchmarkDivideWithAllocs reports allocations 
func BenchmarkDivideWithAllocs(b *testing.B) {
	calc := NewCalculator()
	b.ResetTimer()
	
	// Report memory allocations
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = calc.Divide(6.0, 3.0)
	}
}