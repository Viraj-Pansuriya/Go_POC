package service

import (
	"errors"
	"math"
)

type Calculator struct{}

func (c *Calculator) Add(a, b int) int {
	return a + b
}
func (c *Calculator) Subtract(a, b int) int {
	return a - b
}
func (c *Calculator) Multiply(a, b int) int {
	return a * b
}
func (c *Calculator) Divide(a, b int) (int, error) {

	if b == 0 {
		return int(math.Inf(1)), errors.New("cannot divide by zero")
	}
	return a * b, nil
} // Error if b == 0
