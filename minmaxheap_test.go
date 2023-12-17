package hnsw

import (
	"math/rand"
	"slices"
	"testing"
)

func TestMinMaxHeap(t *testing.T) {
	// Test data
	data := []int{8, 71, 41, 31, 10, 11, 16, 46, 51, 31, 21}
	rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })

	t.Logf("data: %v", data)
	// Create a new minMaxQueue
	queue := NewMinMaxHeap[int](func(a, b int) bool {
		return a < b
	}, data)

	t.Logf("queue: %v", queue.Data)

	// Test IsMinMaxQueue2() function
	if !queue.isMinMaxHeap() {
		t.Errorf("Expected isMinMaxQueue() to return true, but got false")
		t.FailNow()
	}

	// Test Len() method
	if queue.Len() != uint(len(data)) {
		t.Errorf("Expected Len() to return %d, but got %d", len(data), queue.Len())
	}

	// Test Push() method
	queue.Push(13)
	if queue.Len() != uint(len(data)+1) {
		t.Errorf("Expected Len() to return %d, but got %d", len(data)+1, queue.Len())
	}

	// Test PopMin() method
	min := queue.PopMin()
	if min != 8 {
		t.Errorf("Expected PopMin() to return %d, but got %d", 8, min)
	}

	// Test PopMax() method
	max := queue.PopMax()
	if max != 71 {
		t.Errorf("Expected PopMax() to return %d, but got %d", 71, max)
	}

	// Test PeekMin() method
	peekMin := queue.PeekMin()
	if peekMin != 10 {
		t.Errorf("Expected PeekMin() to return %d, but got %d", 10, peekMin)
	}

	// Test PeekMax() method
	peekMax := queue.PeekMax()
	if peekMax != 51 {
		t.Errorf("Expected PeekMax() to return %d, but got %d", 51, peekMax)
	}
}

func FuzzMinMaxHeap(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1})
	f.Add([]byte{1, 2})
	f.Add([]byte{2, 1})

	f.Fuzz(func(t *testing.T, bytes []byte) {
		data := loadAsUint8(bytes)

		heap := NewMinMaxHeap[int8](func(a, b int8) bool {
			return a < b
		}, data)

		// Test IsMinMaxHeap() function
		if !heap.isMinMaxHeap() {
			t.Errorf("Expected isMinMaxQueue() to return true, but got false")
			t.FailNow()
		}

		if len(data) == 0 {
			return
		}

		realMin := slices.Min(data)
		realMax := slices.Max(data)

		// Test PeakMin() method
		peekMin := heap.PeekMin()
		if peekMin != realMin {
			t.Errorf("Expected PeekMin() to return %d, but got %d", realMin, peekMin)
		}

		// Test PeakMax() method
		peekMax := heap.PeekMax()
		if peekMax != realMax {
			t.Errorf("Expected PeekMax() to return %d, but got %d", realMax, peekMax)
		}

		// Test Len() method
		if heap.Len() != uint(len(data)) {
			t.Errorf("Expected Len() to return %d, but got %d", len(data), heap.Len())
		}

		// Test PopMin() method
		min := heap.PopMin()
		if min != realMin {
			t.Errorf("Expected PopMin() to return %d, but got %d", realMin, min)
		}

		if len(data) == 1 {
			return
		}

		// Test PopMax() method
		max := heap.PopMax()
		if max != realMax {
			t.Errorf("Expected PopMax() to return %d, but got %d", realMax, max)
		}

		// Test Len() method
		if heap.Len() != uint(len(data)-2) {
			t.Errorf("Expected Len() to return %d, but got %d", len(data)-2, heap.Len())
		}

	})
}

func loadAsUint8(data []byte) []int8 {
	var out []int8
	for _, b := range data {
		out = append(out, int8(b))
	}
	return out
}
