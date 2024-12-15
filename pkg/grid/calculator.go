package grid

import (
	"fmt"
)

// CalculateGridLevels calculates the price levels for a grid trading strategy
// lowerPrice: the lowest price in the grid
// upperPrice: the highest price in the grid
// gridNum: number of grid levels (must be >= 2)
func CalculateGridLevels(lowerPrice, upperPrice float64, gridNum int) []float64 {
	if err := ValidateGridParams(lowerPrice, upperPrice, gridNum); err != nil {
		return nil
	}

	levels := make([]float64, gridNum)
	// Calculate the price difference between each grid level
	interval := (upperPrice - lowerPrice) / float64(gridNum-1)

	// Generate grid levels
	for i := 0; i < gridNum; i++ {
		levels[i] = lowerPrice + (float64(i) * interval)
	}

	return levels
}

// ValidateGridParams validates the input parameters for grid calculation
func ValidateGridParams(lowerPrice, upperPrice float64, gridNum int) error {
	if lowerPrice <= 0 || upperPrice <= 0 {
		return fmt.Errorf("prices must be positive")
	}

	if lowerPrice >= upperPrice {
		return fmt.Errorf("upper price must be greater than lower price")
	}

	if gridNum < 2 {
		return fmt.Errorf("grid number must be at least 2")
	}

	return nil
}
