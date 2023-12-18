package hnsw

// import (
// 	"container/heap"
// 	"fmt"
// 	"github.com/stretchr/testify/assert"
// 	"math/rand"
// 	"slices"
// 	"sort"
// 	"testing"
// 	"time"
// )

// func TestMinDistancePriorityQueue(t *testing.T) {
// 	dpq := DistancePriorityQueue[[]float32, float32]{
// 		PriorityQueue: []*NodeAndVector[[]float32, float32]{},
// 		Origin:        []float32{0, 0, 0},
// 		Dist:          L2Distance32,
// 		LessFn:        Less[float32],
// 	}

// 	oneNode := &NodeAndVector[[]float32, float32]{
// 		Node:     &HNSWNode{},
// 		Vector:   []float32{1, 0, 0},
// 		distance: 0,
// 	}
// 	twoNode := &NodeAndVector[[]float32, float32]{
// 		Node:     &HNSWNode{},
// 		Vector:   []float32{0, 2, 0},
// 		distance: 0,
// 	}
// 	threeNode := &NodeAndVector[[]float32, float32]{
// 		Node:     &HNSWNode{},
// 		Vector:   []float32{0, 0, 3},
// 		distance: 0,
// 	}

// 	heap.Push(&dpq, twoNode)

// 	assert.Equal(t, 1, dpq.Len())
// 	assert.Equal(t, twoNode, dpq.Peek())

// 	heap.Push(&dpq, oneNode)

// 	assert.Equal(t, 2, dpq.Len())
// 	assert.Equal(t, oneNode, dpq.Peek())

// 	heap.Push(&dpq, threeNode)

// 	assert.Equal(t, 3, dpq.Len())
// 	assert.Equal(t, oneNode, dpq.Peek())

// 	// Drain and check order
// 	assert.Equal(t, oneNode, heap.Pop(&dpq))
// 	assert.Equal(t, twoNode, heap.Pop(&dpq))
// 	assert.Equal(t, threeNode, heap.Pop(&dpq))
// }

// func TestMaxDistancePriorityQueue(t *testing.T) {
// 	dpq := DistancePriorityQueue[[]float32, float32]{}
// 	dpq.PriorityQueue = []*NodeAndVector[[]float32, float32]{}
// 	dpq.Origin = []float32{0, 0, 0}
// 	dpq.Dist = L2Distance32
// 	dpq.LessFn = GreaterEqual[float32]

// 	oneNode := &NodeAndVector[[]float32, float32]{
// 		Node:   &HNSWNode{},
// 		Vector: []float32{1, 0, 0},
// 	}

// 	twoNode := &NodeAndVector[[]float32, float32]{
// 		Node:   &HNSWNode{},
// 		Vector: []float32{0, 2, 0},
// 	}

// 	threeNode := &NodeAndVector[[]float32, float32]{
// 		Node:   &HNSWNode{},
// 		Vector: []float32{0, 0, 3},
// 	}

// 	heap.Push(&dpq, twoNode)

// 	assert.Equal(t, 1, dpq.Len())
// 	assert.Equal(t, twoNode, dpq.Peek())

// 	heap.Push(&dpq, oneNode)

// 	assert.Equal(t, 2, dpq.Len())
// 	assert.Equal(t, twoNode, dpq.Peek())

// 	heap.Push(&dpq, threeNode)

// 	assert.Equal(t, 3, dpq.Len())
// 	assert.Equal(t, threeNode, dpq.Peek())

// 	// Drain and check order
// 	assert.Equal(t, threeNode, heap.Pop(&dpq))
// 	assert.Equal(t, twoNode, heap.Pop(&dpq))
// 	assert.Equal(t, oneNode, heap.Pop(&dpq))
// }

// // TestDistancePriorityQueue - Generalized test for both Min and Max DistancePriorityQueue
// func TestDistancePriorityQueue(t *testing.T) {
// 	testRuns := 1 // Number of iterations for the test

// 	dimensions := rand.Intn(2000) + 1 // Randomly generate the number of dimensions for the vectors

// 	for i := 0; i < testRuns; i++ {
// 		numNodes := rand.Intn(1000000) // Randomly generate the number of nodes to test with
// 		// Test for MinDistancePriorityQueue
// 		dpqMin, nodesMin, sortedMin := setupTestQueue("min", numNodes, dimensions)
// 		performTest(t, dpqMin, nodesMin, sortedMin, "min")

// 		// Test for MaxDistancePriorityQueue
// 		dpqMax, nodesMax, sortedMax := setupTestQueue("max", numNodes, dimensions)
// 		performTest(t, dpqMax, nodesMax, sortedMax, "max")
// 	}
// }

// // TODO: this is not a good benchmark!
// func BenchmarkDistancePriorityQueue(b *testing.B) {
// 	// numNodes := []int{64}
// 	numNodes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048}

// 	dimensions := []int{1536}
// 	//dimensions := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048}
// 	for _, n := range numNodes {
// 		for _, dim := range dimensions {
// 			dpqMin, nodesMin, sortedMin := setupTestQueue("min", n, dim)
// 			dpqMax, nodesMax, sortedMax := setupTestQueue("max", n, dim)

// 			b.Run(fmt.Sprintf("MinDistancePriorityQueue (n=%d, dim=%d)", n, dim), func(b *testing.B) {
// 				for i := 0; i < b.N; i++ {
// 					performTest(b, dpqMin, nodesMin, sortedMin, "min")
// 				}
// 			})

// 			b.Run(fmt.Sprintf("MaxDistancePriorityQueue (n=%d, dim=%d)", n, dim), func(b *testing.B) {
// 				for i := 0; i < b.N; i++ {
// 					performTest(b, dpqMax, nodesMax, sortedMax, "max")
// 				}
// 			})
// 		}
// 	}
// }

// func setupTestQueue(queueType string, numNodes int, dimensions int) (DistancePriorityQueue[[]float32, float32], []*NodeAndVector[[]float32, float32], []*NodeAndVector[[]float32, float32]) {
// 	var dpq DistancePriorityQueue[[]float32, float32]
// 	dpq.PriorityQueue = []*NodeAndVector[[]float32, float32]{}
// 	dpq.Origin = randomVector(dimensions) // Assuming 3D vectors

// 	if queueType == "min" {
// 		dpq.Dist = L2Distance32
// 		dpq.LessFn = Less[float32]
// 	} else {
// 		dpq.Dist = L2Distance32
// 		dpq.LessFn = GreaterEqual[float32]
// 	}

// 	nodes := make([]*NodeAndVector[[]float32, float32], numNodes)
// 	for i := 0; i < numNodes; i++ {
// 		nodes[i] = &NodeAndVector[[]float32, float32]{
// 			Node:   &HNSWNode{},
// 			Vector: randomVector(dimensions),
// 		}
// 	}

// 	sortedNodes := sortNodes(nodes, dpq.Dist, dpq.Origin, queueType)

// 	return dpq, nodes, sortedNodes
// }

// func randomVector(size int) []float32 {
// 	rand.New(rand.NewSource(time.Now().UnixNano()))
// 	vector := make([]float32, size)
// 	for i := range vector {
// 		vector[i] = rand.Float32()
// 	}
// 	return vector
// }

// func performTest(t testing.TB, dpq DistancePriorityQueue[[]float32, float32], nodes []*NodeAndVector[[]float32, float32], sortedNodes []*NodeAndVector[[]float32, float32], queueType string) {
// 	for _, node := range nodes {
// 		heap.Push(&dpq, node)
// 		// Additional checks for queue length and top element can be added here
// 	}
// 	// Drain and check order
// 	for _, expectedNode := range sortedNodes {
// 		actualNode := heap.Pop(&dpq).(*NodeAndVector[[]float32, float32])
// 		// In the event of a tie, the order of the nodes is not guaranteed, so we only
// 		// fail the test if the distance is not the same
// 		assert.Equal(t, expectedNode.distance, actualNode.distance)
// 	}
// }

// func sortNodes(nodes []*NodeAndVector[[]float32, float32], distFn func([]float32, []float32) float32, origin []float32, queueType string) []*NodeAndVector[[]float32, float32] {
// 	// Calculate the distance for each node and store it in a map
// 	distances := make(map[*NodeAndVector[[]float32, float32]]float32)
// 	for _, node := range nodes {
// 		distances[node] = distFn(node.Vector, origin)
// 	}

// 	sorted := slices.Clone(nodes)

// 	// Sort the nodes based on the calculated distances
// 	sort.Slice(sorted, func(i, j int) bool {
// 		if queueType == "min" {
// 			return distances[sorted[i]] < distances[sorted[j]]
// 		} else { // for "max"
// 			return distances[sorted[i]] > distances[sorted[j]]
// 		}
// 	})

// 	return sorted
// }
