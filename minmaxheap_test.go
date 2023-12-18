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
	}

	// Test Len() method
	if queue.Len() != len(data) {
		t.Errorf("Expected Len() to return %d, but got %d", len(data), queue.Len())
	}

	// Test Push() method
	queue.Push(13)
	if queue.Len() != len(data)+1 {
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

func FuzzMinMaxHeap2(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1})
	f.Add([]byte{1, 2})
	f.Add([]byte{2, 1})

	f.Fuzz(func(t *testing.T, bytes []byte) {
		data := loadAsUint8(bytes)

		heap := NewMinMaxHeap[int8](func(a, b int8) bool {
			return a < b
		}, data)

		if heap.Len() == 0 {
			return
		}
		extracted := make([]int8, 0, len(data))
		for heap.Len() > 0 {
			extracted = append(extracted, heap.PopMin())
		}

		if !slices.IsSorted(extracted) {
			t.Errorf("Expected extracted data to be sorted, but got %v", extracted)
		}

		heap = NewMinMaxHeap[int8](func(a, b int8) bool {
			return a < b
		}, data)

		extracted = make([]int8, 0, len(data))
		for heap.Len() > 0 {
			extracted = append(extracted, heap.PopMax())
		}

		if !slices.IsSortedFunc(extracted, func(a, b int8) int { return int(b) - int(a) }) {
			t.Errorf("Expected extracted data to be sorted, but got %v", extracted)
		}
	})
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
		if heap.Len() != len(data) {
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
		if heap.Len() != len(data)-2 {
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

func BenchmarkMinMaxHeap(b *testing.B) {
	b.Run("Min", func(b *testing.B) {
		benchmarkMinMaxHeap(b, func(a, b int) bool {
			return a < b
		})
	})
	b.Run("Max", func(b *testing.B) {
		benchmarkMinMaxHeap(b, func(a, b int) bool {
			return a > b
		})
	})
	b.Run("Push", func(b *testing.B) {
		benchmarkMinMaxHeapPush(b)
	})
	b.Run("PopMin", func(b *testing.B) {
		benchmarkMinMaxHeapPopMin(b)
	})
	b.Run("PopMax", func(b *testing.B) {
		benchmarkMinMaxHeapPopMax(b)
	})
	b.Run("VeryLarge", func(b *testing.B) {
		origin := make([]float32, 1000)
		benchmarkLargeMinMaxHeap(b, func(a, b []float32) bool {
			return L2Distance32(a, origin) < L2Distance32(b, origin)
		})
	})
}

func benchmarkLargeMinMaxHeap(b *testing.B, cmp func(a, b []float32) bool) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make([][]float32, 1_000_000)
		for i := range data {
			data[i] = make([]float32, 1000)
			for j := range data[i] {
				data[i][j] = rand.Float32()
			}
		}
		b.StartTimer()
		NewMinMaxHeap[[]float32](cmp, data)
	}
}

func benchmarkMinMaxHeap(b *testing.B, cmp func(a, b int) bool) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make([]int, 1_000_000)
		for i := range data {
			data[i] = rand.Int()
		}
		b.StartTimer()
		NewMinMaxHeap[int](func(a, b int) bool {
			return a < b
		}, data)
	}
}

func benchmarkMinMaxHeapPush(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		heap := NewMinMaxHeap[int](func(a, b int) bool {
			return a < b
		}, []int{})
		b.StartTimer()
		for i := 0; i < 1_000_000; i++ {
			heap.Push(rand.Int())
		}
	}
}

func benchmarkMinMaxHeapPopMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make([]int, 1_000_000)
		for i := range data {
			data[i] = rand.Int()
		}
		heap := NewMinMaxHeap[int](func(a, b int) bool {
			return a < b
		}, data)
		b.StartTimer()
		for i := 0; i < 1_000_000; i++ {
			heap.PopMin()
		}
	}
}

func benchmarkMinMaxHeapPopMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make([]int, 1_000_000)
		for i := range data {
			data[i] = rand.Int()
		}
		heap := NewMinMaxHeap[int](func(a, b int) bool {
			return a < b
		}, data)
		b.StartTimer()
		for i := 0; i < 1_000_000; i++ {
			heap.PopMax()
		}
	}
}

func TestMinMaxHeap_Unit(t *testing.T) {

	heap := NewMinMaxHeap[int](func(a, b int) bool {
		return a < b
	}, []int{})

	for i := 0; i < 100; i++ {
		heap.Push(i)
		if !heap.isMinMaxHeap() {
			t.Errorf("Expected isMinMaxQueue() to return true, but got false")
			t.FailNow()
		}
	}

	max := heap.PopMax()

	if max != 99 {
		t.Errorf("Expected max to be 99, but got %d", max)
	}

}
