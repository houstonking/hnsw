package hnsw

import (
	"container/heap"

	"golang.org/x/exp/constraints"
)

type NodeAndVector[S ~[]F, F constraints.Float] struct {
	Node     *HNSWNode
	Vector   S
	distance F
}

type PriorityQueue[S ~[]F, F constraints.Float] []*NodeAndVector[S, F]

func (pq PriorityQueue[S, F]) Len() int { return len(pq) }

func (pq PriorityQueue[S, F]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue[S, F]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

// A DistancePriorityQueue implements heap.Interface, and holds NodeAndVector.
// The priority of a node is its distance from the origin. If LessFn is less,
// the priority is the smaller distance, otherwise it is the larger distance.
type DistancePriorityQueue[S ~[]F, F constraints.Float] struct {
	PriorityQueue[S, F]
	Origin S
	Dist   DistanceFunc[S, F]
	LessFn func(a, b F) bool
}

func (dpq *DistancePriorityQueue[S, F]) Push(x any) {
	item := x.(*NodeAndVector[S, F])
	item.distance = dpq.Dist(item.Vector, dpq.Origin)
	dpq.PriorityQueue = append(dpq.PriorityQueue, item)
}

func (dpq DistancePriorityQueue[S, F]) Less(i, j int) bool {
	return dpq.LessFn(dpq.PriorityQueue[i].distance, dpq.PriorityQueue[j].distance)
}

func (dpq *DistancePriorityQueue[S, F]) Peek() *NodeAndVector[S, F] {
	if dpq.Len() == 0 {
		return nil
	}
	return dpq.PriorityQueue[0]
}

func MinDistQueue[S ~[]F, F constraints.Float](dist DistanceFunc[S, F], nodes []*HNSWNode, vectors []S, origin S) *DistancePriorityQueue[S, F] {
	nv := make([]*NodeAndVector[S, F], len(nodes))
	for i, node := range nodes {
		nv[i] = &NodeAndVector[S, F]{
			Node:   node,
			Vector: vectors[i],
		}
	}
	dpq := &DistancePriorityQueue[S, F]{
		PriorityQueue: nv,
		Origin:        origin,
		Dist:          dist,
		LessFn:        Less[F],
	}
	heap.Init(dpq)
	return dpq
}

func MaxDistQueue[S ~[]F, F constraints.Float](dist DistanceFunc[S, F], nodes []*HNSWNode, vectors []S, origin S) *DistancePriorityQueue[S, F] {
	dpq := &DistancePriorityQueue[S, F]{
		PriorityQueue: []*NodeAndVector[S, F]{},
		Origin:        origin,
		Dist:          dist,
		LessFn:        GreaterEqual[F],
	}
	heap.Init(dpq)
	for i, node := range nodes {
		heap.Push(dpq, &NodeAndVector[S, F]{
			Node:   node,
			Vector: vectors[i],
		})
	}
	return dpq
}

type DistanceMinMaxHeap[S ~[]F, F constraints.Float] struct {
	heap   *minMaxHeap[NodeAndVector[S, F]]
	origin S
	dist   DistanceFunc[S, F]
}

func NewDistanceMinMaxHeap[S ~[]F, F constraints.Float](dist DistanceFunc[S, F], origin S) *DistanceMinMaxHeap[S, F] {
	return &DistanceMinMaxHeap[S, F]{
		heap: NewMinMaxHeap[NodeAndVector[S, F]](func(nav1, nav2 NodeAndVector[S, F]) bool {
			return Less[F](dist(nav1.Vector, origin), dist(nav2.Vector, origin))
		}, []NodeAndVector[S, F]{}),
		origin: origin,
		dist:   dist,
	}
}
