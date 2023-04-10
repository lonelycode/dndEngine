package game

import (
	"testing"
)

func TestParseDiceRoll(t *testing.T) {
	testCases := []struct {
		input        string
		expectedDice *Dice
		expectErr    bool
	}{
		{"2d6", &Dice{NumDice: 2, Dice: D6, Mod: 0}, false},
		{"1d20+3", &Dice{NumDice: 1, Dice: D20, Mod: 3}, false},
		{"d8", nil, true},
		{"1d6+", nil, true},
		{"1d6+x", nil, true},
		{"2d5", nil, true},
	}

	for _, tc := range testCases {
		actualDice, actualErr := ParseDiceRoll(tc.input)
		if actualErr != nil && !tc.expectErr {
			t.Errorf("Unexpected error for input %s: %v", tc.input, actualErr)
		}
		if actualErr == nil && tc.expectErr {
			t.Errorf("Expected an error, but none occurred for input %s", tc.input)
		}
		if actualDice != nil && !actualDice.Equals(tc.expectedDice) {
			t.Errorf("Unexpected dice for input %s; got %+v, want %+v", tc.input, actualDice, tc.expectedDice)
		}
	}
}
