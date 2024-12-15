package grid

import (
	"testing"
)

func TestCalculateGridLevels(t *testing.T) {
	tests := []struct {
		name           string
		lowerPrice     float64
		upperPrice     float64
		gridNum        int
		expectedLevels []float64
	}{
		{
			name:           "Basic grid with 5 levels",
			lowerPrice:     100.0,
			upperPrice:     200.0,
			gridNum:        5,
			expectedLevels: []float64{100.0, 125.0, 150.0, 175.0, 200.0},
		},
		{
			name:           "Grid with 3 levels",
			lowerPrice:     1000.0,
			upperPrice:     1300.0,
			gridNum:        3,
			expectedLevels: []float64{1000.0, 1150.0, 1300.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			levels := CalculateGridLevels(tt.lowerPrice, tt.upperPrice, tt.gridNum)

			if len(levels) != len(tt.expectedLevels) {
				t.Errorf("Expected %d levels, got %d", len(tt.expectedLevels), len(levels))
				return
			}

			for i := range levels {
				if levels[i] != tt.expectedLevels[i] {
					t.Errorf("Level %d: expected %.2f, got %.2f", i, tt.expectedLevels[i], levels[i])
				}
			}
		})
	}
}

func TestValidateGridParams(t *testing.T) {
	tests := []struct {
		name       string
		lowerPrice float64
		upperPrice float64
		gridNum    int
		wantErr    bool
	}{
		{
			name:       "Valid parameters",
			lowerPrice: 100.0,
			upperPrice: 200.0,
			gridNum:    5,
			wantErr:    false,
		},
		{
			name:       "Invalid - lower price greater than upper",
			lowerPrice: 200.0,
			upperPrice: 100.0,
			gridNum:    5,
			wantErr:    true,
		},
		{
			name:       "Invalid - zero grid number",
			lowerPrice: 100.0,
			upperPrice: 200.0,
			gridNum:    0,
			wantErr:    true,
		},
		{
			name:       "Invalid - negative prices",
			lowerPrice: -100.0,
			upperPrice: 200.0,
			gridNum:    5,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGridParams(tt.lowerPrice, tt.upperPrice, tt.gridNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGridParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
