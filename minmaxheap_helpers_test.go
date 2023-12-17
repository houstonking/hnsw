package hnsw

import (
	"testing"
)

func TestHighestSetBit(t *testing.T) {
	// Test data
	testData := []struct {
		input    uint
		expected uint
	}{
		{0b00000000, 0},
		{0b00000001, 1},
		{0b00000010, 2},
		{0b00000011, 2},
		{0b00000100, 3},
		{0b00000111, 3},
		{0b00001000, 4},
		{0b00010000, 5},
		{0b00100000, 6},
		{0b01000000, 7},
		{0b10000000, 8},
		{0b10000001, 8},
		{0b10000010, 8},
		{0b10000100, 8},
		{0b10001000, 8},
		{0b10010000, 8},
		{0b10100000, 8},
		{0b11000000, 8},
		{0b11111111, 8},
	}

	// Test highestSetBit() function
	for _, test := range testData {
		result := highestSetBit(test.input)
		if result != test.expected {
			t.Errorf("Expected highestSetBit(%d) to return %d, but got %d", test.input, test.expected, result)
		}
	}
}

// func BenchmarkHighestSetBit(b *testing.B) {
// 	b.Run("highestSetBit", func(b *testing.B) {
// 		var res uint
// 		for i := 0; i < b.N; i++ {
// 			res = highestSetBit(uint(i))
// 		}
// 		b.Log("res: ", res)
// 	})

// 	b.Run("highestSetBit2", func(b *testing.B) {
// 		var res uint
// 		for i := 0; i < b.N; i++ {
// 			res = highestSetBit2(uint(i))
// 		}
// 		b.Log("res: ", res)
// 	})
// }

func Test_isNewItemMin(t *testing.T) {
	// Test data
	testData := []struct {
		input    uint
		expected bool
	}{
		{0b00000000, true},
		{0b00000001, false},
		{0b00000010, true},
		{0b00000011, true},
		{0b00000100, false},
		{0b00000111, false},
		{0b00001000, true},
		{0b00001011, true},
		{0b00001100, true},
		{0b00001111, true},
		{0b00010000, false},
	}

	// Test isNewItemMin() function
	for _, test := range testData {
		result := isNewItemMin(test.input)
		if result != test.expected {
			t.Errorf("Expected isNewItemMin(%d) to return %t, but got %t", test.input, test.expected, result)
		}
	}
}

func Test_isMinItem(t *testing.T) {
	// Test data
	testData := []struct {
		input    uint
		expected bool
	}{
		{0, true},
		{1, false},
		{2, false},
		{3, true},
		{4, true},
		{5, true},
		{6, true},
		{7, false},
		{8, false},
	}

	// Test isMinItem() function
	for _, test := range testData {
		result := isMinLevelIndex(test.input)
		if result != test.expected {
			t.Errorf("Expected isMinItem(%d) to return %t, but got %t", test.input, test.expected, result)
		}
	}
}
