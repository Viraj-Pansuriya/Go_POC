package test

import (
	"testing"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/05-testing/service"
)

// calculator/calculator_test.go
func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"zero", 0, 0, 0},
		{"mixed", -5, 10, 5},
	}

	calc := &service.Calculator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 6},
		{"negative numbers", -2, -3, 6},
		{"zero", 0, 0, 0},
		{"mixed", -5, 10, -50},
	}

	calc := &service.Calculator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	calc := &service.Calculator{}

	// Test division by zero
	_, err := calc.Divide(10, 0)
	if err == nil {
		t.Error("Expected error for division by zero")
	}

	// Test normal division
	result, err := calc.Divide(10, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}
