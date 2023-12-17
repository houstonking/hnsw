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
	len  uint
	Data []T
	Less func(T, T) bool
}

type compare[T any] func([]T, int, int) bool

func NewMinMaxHeap[T any](less func(T, T) bool, data []T) *minMaxHeap[T] {
	n := len(data)
	q := &minMaxHeap[T]{Less: less, Data: data, len: uint(n)}
	for i := n/2 - 1; i >= 0; i-- {
		q.down(i, n)
	}
	return q
}

func (heap *minMaxHeap[T]) Len() uint {
	return heap.len
}

func (heap *minMaxHeap[T]) Push(x T) {
	i := int(heap.Len())
	heap.Data = append(heap.Data, x)
	heap.up(i, i+1)
	heap.len++
}

func (heap *minMaxHeap[T]) PopMin() T {
	i := int(heap.Len() - 1)

	heap.swap(0, i)
	heap.down(0, i)
	heap.len--
	item := heap.Data[i]

	return item
}

func (heap *minMaxHeap[T]) PopMax() T {
	n := int(heap.Len())
	i := heap.indexOfMax(n)
	j := n - 1

	item := heap.Data[i]

	heap.swap(i, j)
	heap.down(i, j)
	heap.len--

	return item
}

// indexOfMax returns index of the maximum element in h.
func (heap *minMaxHeap[T]) indexOfMax(n int) int {
	if n <= 2 {
		return n - 1
	}

	if heap.LessFn()(heap.Data, 2, 1) {
		return 1
	}

	return 2
}

func (heap *minMaxHeap[T]) PeekMin() T {
	return heap.Data[0]
}

func (heap *minMaxHeap[T]) PeekMax() T {
	return heap.Data[heap.indexOfMax(int(heap.Len()))]
}

func (heap *minMaxHeap[T]) swap(i int, j int) {
	heap.Data[i], heap.Data[j] = heap.Data[j], heap.Data[i]
}

func (heap *minMaxHeap[T]) LessFn() compare[T] {
	return func(data []T, i int, j int) bool {
		return heap.Less(data[i], data[j])
	}
}

func (heap *minMaxHeap[T]) GreaterFn() compare[T] {
	return func(data []T, i int, j int) bool {
		return heap.Less(data[j], data[i])
	}
}

// up moves the element at i upwards within the heap until it occupies an
// appropriate node.
func (heap *minMaxHeap[T]) up(i, n int) {
	parent := (i - 1) / 2

	if isMinLevelIndex(uint(i)) {
		if i > 0 && heap.swapIf(heap.GreaterFn(), i, parent) {
			heap.upX(heap.GreaterFn(), parent, n)
		} else {
			heap.upX(heap.LessFn(), i, n)
		}
	} else {
		if i > 0 && heap.swapIf(heap.LessFn(), i, parent) {
			heap.upX(heap.LessFn(), parent, n)
		} else {
			heap.upX(heap.GreaterFn(), i, n)
		}
	}
}

func (heap *minMaxHeap[T]) upX(less compare[T], i, n int) {
	for i > 2 {
		grandparent := (((i - 1) / 2) - 1) / 2

		if !heap.swapIf(less, i, grandparent) {
			return
		}

		i = grandparent
	}
}

func (heap *minMaxHeap[T]) down(i int, n int) bool {
	if isMinLevelIndex(uint(i)) {
		return heap.downX(heap.LessFn(), i, n)
	}
	return heap.downX(heap.GreaterFn(), i, n)
}

func (heap *minMaxHeap[T]) downX(cmp compare[T], i int, n int) bool {
	recursed := false

	for {
		m := heap.minDescendent(cmp, i, n)
		if m == -1 {
			// i has no children.
			return recursed
		}

		parent := (m - 1) / 2

		if i == parent {
			// m is a direct child of i.
			heap.swapIf(cmp, m, i)
			return recursed
		}

		// m is a grandchild of i.
		if !heap.swapIf(cmp, m, i) {
			return recursed
		}

		heap.swapIf(cmp, parent, m)

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

// least returns the index of the smaller element of those elements at i and j.
//
// If j overruns the heap, done is true.
func (heap *minMaxHeap[T]) least(cmp compare[T], i, j, n int) (_ int, done bool) {
	if j >= n {
		return i, true
	}

	if cmp(heap.Data, i, j) {
		return i, false
	}

	return j, false
}

func (heap *minMaxHeap[T]) swapIf(cmp compare[T], i int, j int) bool {
	if cmp(heap.Data, i, j) {
		heap.swap(i, j)
		return true
	}
	return false
}

// isMinMaxHeap returns true if the heap is a min-max heap, false otherwise.
func (heap *minMaxHeap[T]) isMinMaxHeap() bool {
	length := heap.len
	testIndex := func(index uint, compareIndex func(uint) bool) bool {
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
	for i := uint(0); i < length; i++ {
		if isMinLevelIndex(i) {
			// fmt.Println("value: ", heap.Data[i], "index: ", i, "is min")
			compareOne := func(child uint) bool {
				ret := child >= length || !heap.Less(heap.Data[child], heap.Data[i])
				// fmt.Println("child: ", child, "value: ", heap.Data[child], "index: ", child, "is less than: ", heap.Data[i], "index: ", i, "result: ", ret)
				return ret
			}
			if !testIndex(i, compareOne) {
				return false
			}
		} else {
			// fmt.Println("value: ", heap.Data[i], "index: ", i, "is max")
			compareOne := func(child uint) bool {

				ret := child >= length || !heap.Less(heap.Data[i], heap.Data[child])
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
