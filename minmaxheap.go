package hnsw

// minMaxHeap is a min-max heap that supports the following operations:
// - Push(x T) - add x to the queue
// - PopMin() T - remove and return the smallest element in the queue
// - PopMax() T - remove and return the largest element in the queue
// - PeekMin() T - return the smallest element in the queue
// - PeekMax() T - return the largest element in the queue
// - Len() uint - return the number of elements in the queue
// This implementation is based off of the repository:
// https://github.com/dogmatiq/kyu/blob/main/README.md
type minMaxHeap[T any] struct {
	Data      []T
	lessFn    compare[T]
	greaterFn compare[T]
}

type compare[T any] func([]T, int, int) bool

func NewMinMaxHeap[T any](less func(T, T) bool, data []T) *minMaxHeap[T] {
	n := len(data)
	q := &minMaxHeap[T]{
		lessFn: func(data []T, i int, j int) bool {
			return less(data[i], data[j])
		},
		greaterFn: func(data []T, i int, j int) bool {
			return less(data[j], data[i])
		},
		Data: data,
	}
	if n == 0 {
		return q
	}
	for i := n/2 - 1; i >= 0; i-- {
		q.down(i, n)
	}
	return q
}

func NewMinMaxHeap2[T any](less func(T, T) bool, data []T) *minMaxHeap[T] {
	n := len(data)
	q := &minMaxHeap[T]{
		lessFn: func(data []T, i int, j int) bool {
			return less(data[i], data[j])
		},
		greaterFn: func(data []T, i int, j int) bool {
			return less(data[j], data[i])
		},
		Data: data,
	}
	for i := n/2 - 1; i >= 0; i-- {
		q.down(i, n)
	}
	return q
}

func (heap *minMaxHeap[T]) Len() int {
	return len(heap.Data)
}

func (heap *minMaxHeap[T]) Push(x T) {
	i := int(heap.Len())
	heap.Data = append(heap.Data, x)
	heap.up(i, i+1)
}

func (heap *minMaxHeap[T]) PopMin() T {
	n := int(heap.Len() - 1)

	heap.swap(0, n)
	heap.down(0, n)
	item := heap.Data[n]
	heap.Data = heap.Data[:n]

	return item
}

func (heap *minMaxHeap[T]) PopMax() T {
	n := int(heap.Len())
	i := heap.indexOfMax()
	j := n - 1

	heap.swap(i, j)
	heap.down(i, j)

	// i := 0
	// l := 1
	// if l < n && heap.lessFn(heap.Data, l, i) {
	// 	i = l
	// }
	// r := 2
	// if r < n && heap.lessFn(heap.Data, r, i) {
	// 	i = r
	// }

	// heap.swap(i, n-1)
	// heap.down(i, n-1)

	item := heap.Data[j]
	heap.Data = heap.Data[:j]

	return item
}

// indexOfMax returns index of the maximum element in h.
func (heap *minMaxHeap[T]) indexOfMax() int {
	n := int(heap.Len())
	if n <= 2 {
		return n - 1
	}

	if heap.lessFn(heap.Data, 2, 1) {
		return 1
	}

	return 2
}

func (heap *minMaxHeap[T]) PeekMin() T {
	return heap.Data[0]
}

func (heap *minMaxHeap[T]) PeekMax() T {
	return heap.Data[heap.indexOfMax()]
}

func (heap *minMaxHeap[T]) swap(i int, j int) {
	heap.Data[i], heap.Data[j] = heap.Data[j], heap.Data[i]
}

// up moves the element at i upwards within the heap until it occupies an
// appropriate node.
func (heap *minMaxHeap[T]) up(i, n int) {
	parent := (i - 1) / 2

	if isMinLevelIndex(uint(i)) {
		if i > 0 && heap.swapIfGreater(i, parent) {
			heap.upX(heap.greaterFn, parent, n)
		} else {
			heap.upX(heap.lessFn, i, n)
		}
	} else {
		if i > 0 && heap.swapIfLess(i, parent) {
			heap.upX(heap.lessFn, parent, n)
		} else {
			heap.upX(heap.greaterFn, i, n)
		}
	}
}

func (heap *minMaxHeap[T]) swapIf(cmp compare[T], i int, j int) bool {
	if cmp(heap.Data, i, j) {
		heap.swap(i, j)
		return true
	}
	return false
}

func (heap *minMaxHeap[T]) upX(cmp compare[T], i, n int) {
	for i > 2 {
		grandparent := (((i - 1) / 2) - 1) / 2

		if !heap.swapIf(cmp, i, grandparent) {
			return
		}

		i = grandparent
	}
}

func (heap *minMaxHeap[T]) down(i int, n int) bool {
	if isMinLevelIndex(uint(i)) {
		return heap.downMin(i, n)
	}
	return heap.downMax(i, n)
}

func (heap *minMaxHeap[T]) downMin(i int, n int) bool {
	recursed := false

	for {
		m := heap.minDescendent(heap.lessFn, i, n)
		if m == -1 {
			// i has no children.
			return recursed
		}

		parent := (m - 1) / 2

		if i == parent {
			// m is a direct child of i.
			heap.swapIfLess(m, i)
			return recursed
		}

		// m is a grandchild of i.
		if !heap.swapIfLess(m, i) {
			return recursed
		}

		heap.swapIfLess(parent, m)

		i = m
		recursed = true
	}
}

func (heap *minMaxHeap[T]) downMax(i int, n int) bool {
	recursed := false

	for {
		m := heap.minDescendent(heap.greaterFn, i, n)
		if m == -1 {
			// i has no children.
			return recursed
		}

		parent := (m - 1) / 2

		if i == parent {
			// m is a direct child of i.
			heap.swapIfGreater(m, i)
			return recursed
		}

		// m is a grandchild of i.
		if !heap.swapIfGreater(m, i) {
			return recursed
		}

		heap.swapIfGreater(parent, m)

		i = m
		recursed = true
	}
}

// minDescendent returns the index of the smallest child or grandchild of i.
//
// It returns -1 if i is a leaf node.
func (heap *minMaxHeap[T]) minDescendent(cmp compare[T], i int, n int) int {
	left := i*2 + 1
	if left >= n {
		// i is a leaf node.
		return -1
	}

	// check i's right-hand child.
	right := left + 1
	min, done := heap.least(cmp, left, right, n)
	if done {
		return min
	}

	// check i's left-hand child's own left-hand child.
	min, done = heap.least(cmp, min, left*2+1, n)
	if done {
		return min
	}

	// check i's left-hand child's right-hand child.
	min, done = heap.least(cmp, min, left*2+2, n)
	if done {
		return min
	}

	// check i's right-hand child's right-hand child.
	min, done = heap.least(cmp, min, right*2+1, n)
	if done {
		return min
	}

	// check i's right-hand child's own right-hand child.
	min, _ = heap.least(cmp, min, right*2+2, n)

	return min

}

// firstChild := 2*i + 1
// if firstChild >= n {
// 	return i
// }
// secondChild := firstChild + 1
// if secondChild >= n {
// 	return firstChild
// }

// firstGrandchild := 2*firstChild + 1
// if firstGrandchild >= n {
// 	if cmp(heap.Data, secondChild, firstChild) {
// 		return secondChild
// 	}
// 	return firstChild
// }

// secondGrandchild := firstGrandchild + 1
// if secondGrandchild >= n {
// 	if cmp(heap.Data, firstGrandchild, secondChild) {
// 		return firstGrandchild
// 	}
// 	return secondChild
// }

// minGrandchild := firstGrandchild
// if cmp(heap.Data, secondGrandchild, firstGrandchild) {
// 	minGrandchild = secondGrandchild
// }

// thirdGrandchild := secondGrandchild + 1
// if thirdGrandchild >= n {
// 	if cmp(heap.Data, minGrandchild, secondChild) {
// 		return minGrandchild
// 	}
// 	return secondChild
// }

// if cmp(heap.Data, thirdGrandchild, minGrandchild) {
// 	return thirdGrandchild
// }
// return minGrandchild

// func (heap *minMaxHeap[T]) swapIf(cmp compare[T], i int, j int) bool {
// 	if cmp(heap.Data, i, j) {
// 		heap.swap(i, j)
// 		return true
// 	}
// 	return false
// }

func (heap *minMaxHeap[T]) least(cmp compare[T], i, j, n int) (_ int, done bool) {
	if j >= n {
		return i, true
	}

	if cmp(heap.Data, i, j) {
		return i, false
	}

	return j, false
}

func (heap *minMaxHeap[T]) swapIfLess(i int, j int) bool {
	if heap.lessFn(heap.Data, i, j) {
		heap.swap(i, j)
		return true
	}
	return false
}

func (heap *minMaxHeap[T]) swapIfGreater(i int, j int) bool {
	if heap.greaterFn(heap.Data, i, j) {
		heap.swap(i, j)
		return true
	}
	return false
}

// isMinMaxHeap returns true if the heap is a min-max heap, false otherwise.
func (heap *minMaxHeap[T]) isMinMaxHeap() bool {
	length := heap.Len()
	testIndex := func(index int, compareIndex func(int) bool) bool {
		firstChild := firstChildIndex(index)
		secondChild := firstChild + 1
		firstGrandchild := firstChildIndex(firstChild)
		secondGrandchild := firstGrandchild + 1
		thirdGrandchild := firstChildIndex(secondChild)
		fourthGrandchild := thirdGrandchild + 1
		return compareIndex(firstChild) && compareIndex(secondChild) &&
			compareIndex(firstGrandchild) && compareIndex(secondGrandchild) &&
			compareIndex(thirdGrandchild) && compareIndex(fourthGrandchild)
	}
	for i := 0; i < length; i++ {
		if isMinLevelIndex(uint(i)) {
			// fmt.Println("value: ", heap.Data[i], "index: ", i, "is min")
			compareOne := func(child int) bool {
				ret := child >= length || !heap.lessFn(heap.Data, child, i)
				// fmt.Println("child: ", child, "value: ", heap.Data[child], "index: ", child, "is less than: ", heap.Data[i], "index: ", i, "result: ", ret)
				return ret
			}
			if !testIndex(i, compareOne) {
				return false
			}
		} else {
			// fmt.Println("value: ", heap.Data[i], "index: ", i, "is max")
			compareOne := func(child int) bool {

				ret := child >= length || !heap.lessFn(heap.Data, i, child)
				// fmt.Println("child: ", child, "value: ", heap.Data[child], "index: ", child, "is less than: ", heap.Data[i], "index: ", i, "result: ", ret)

				return ret
			}
			if !testIndex(i, compareOne) {
				// fmt.Println("failed")
				return false
			}
		}
	}
	return true
}
