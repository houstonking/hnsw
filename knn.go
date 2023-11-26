package hnsw

import "golang.org/x/exp/constraints"

func KNN[S ~[]F, F constraints.Float](candidates []S, D DistanceFunc[S, F], q S, k int) []int {
	distances := make([]F, len(candidates))
	for i, candidate := range candidates {
		distances[i] = D(candidate, q)
	}
	return TopK(distances, k)
}

func TopK[F constraints.Float](distances []F, k int) []int {
	ids := make([]int, len(distances))
	for i := range ids {
		ids[i] = i
	}
	quickSelect(ids, distances, 0, len(ids)-1, k)
	return quickSort(ids[:k], distances[:k])
}

func quickSort[F constraints.Float](ids []int, distances []F) []int {
	if len(ids) < 2 {
		return ids
	}
	pivotIndex := len(ids) / 2
	pivotIndex = partition(ids, distances, 0, len(ids)-1, pivotIndex)
	quickSort(ids[:pivotIndex], distances[:pivotIndex])
	quickSort(ids[pivotIndex+1:], distances[pivotIndex+1:])
	return ids
}

func quickSelect[F constraints.Float](ids []int, distances []F, left, right, k int) {
	for {
		if left == right {
			return
		}
		pivotIndex := left + (right-left)/2
		pivotIndex = partition(ids, distances, left, right, pivotIndex)
		if k == pivotIndex {
			return
		} else if k < pivotIndex {
			right = pivotIndex - 1
		} else if k > pivotIndex {
			left = pivotIndex + 1
		}
	}
}

func partition[F constraints.Float](ids []int, distances []F, left, right, pivotIndex int) int {
	pivotValue := distances[pivotIndex]
	swap(ids, distances, pivotIndex, right)
	storeIndex := left
	for i := left; i < right; i++ {
		if distances[i] < pivotValue {
			swap(ids, distances, i, storeIndex)
			storeIndex++
		}
	}
	swap(ids, distances, storeIndex, right)
	return storeIndex
}

func swap[F constraints.Float](ids []int, distances []F, index int, right int) {
	ids[index], ids[right] = ids[right], ids[index]
	distances[index], distances[right] = distances[right], distances[index]
}
