package hnsw

import (
	"math/bits"
)

// highestSetBit returns the index of the highest set bit in x. This is a nice shortcut
// for determining what layer a new item should be inserted into. For example, if we
// have 7 items in the queue, the highest set bit is 3, which is odd, so we know that
// the new item should be inserted into the min layer 3.
func highestSetBit(x uint) uint {
	return uint(bits.Len64(uint64(x)))
}

func isNewItemMin(len uint) bool {
	return (highestSetBit(len) & 1) == 0
}

func isMinLevelIndex(index uint) bool {
	return (highestSetBit(index+1) & 1) == 1
}

func grandparentIndex(i uint) uint {
	return (i - 3) / 4
}

func parentIndex(i uint) uint {
	return (i - 1) / 2
}

func firstChildIndex(i uint) uint {
	return i*2 + 1
}

func lastGrandchildIndex(i uint) uint {
	return i*4 + 6
}
