package hnsw

import (
	"math/rand"

	"github.com/bits-and-blooms/bitset"
	"github.com/chewxy/math32"
	"golang.org/x/exp/constraints"
)

// NodeID defines the type of the node ID. It should be unique for each node in the graph.
type NodeID uint64

type HNSWGraph[S ~[]F, F constraints.Float] struct {
	// D defines the dimensionality of the data.
	D int
	// Dist defines the distance function.
	Dist DistanceFunc[S, F]
	// EfConstruction defines the size of the dynamic list for the nearest neighbors search during the construction of the graph.
	EfConstruction int
	// M defines the maximum number of outgoing connections in the graph.
	M int
	// LevelMult defines the multiplier for the number of outgoing connections per level.
	LevelMult float32
	// LevelProbabilities defines the probabilities of the levels.
	LevelProbabilities []float32
	// LevelNeighbors defines the cumulative number of neighbors per level.
	LevelNeighbors []int
	// RngSource defines the random number generator source.
	Rng *rand.Rand
	// NodeStore defines the nodes of the graph.
	NodeStore NodeStore
	// VectorStore defines the vectors of the graph.
	VectorStore VectorStore[S, F]
	// ep defines the entry point of the graph.
	ep *NodeID
}

func (f *HNSWGraph[S, F]) EntryPoint() NodeID {
	return *f.ep
}

func (f *HNSWGraph[S, F]) Store(node *HNSWNode, vector S) error {
	//fmt.Println("HNSWGraph.Store", node)
	err := f.NodeStore.Add(node.ID, node)
	if err != nil {
		return err
	}
	err = f.VectorStore.Add(node.ID, vector)
	if err != nil {
		return err
	}
	return nil
}

func (f *HNSWGraph[S, F]) EmptyNeighbors() [][]NodeID {
	neighbors := make([][]NodeID, len(f.LevelNeighbors))
	for i := range neighbors {
		neighbors[i] = make([]NodeID, 0)
	}
	return neighbors
}

func (f *HNSWGraph[S, F]) EmptyNode() *HNSWNode {
	return &HNSWNode{
		ID:               f.NodeStore.NextID(),
		Layer:            0,
		NeighborsByLevel: f.EmptyNeighbors(),
	}

}

func (f *HNSWGraph[S, F]) UpdateNeighbor(neighborID NodeID, newNode *HNSWNode, newNodeVector S, layer int) (*HNSWNode, error) {
	mLayer := f.M
	if layer == 0 {
		mLayer += f.M
	}
	neighborNode, err := f.NodeStore.Get(neighborID)
	if err != nil {
		return nil, err
	}
	neighborVector, err := f.VectorStore.Get(neighborID)
	if err != nil {
		return nil, err
	}

	// nearest := MinDistQueue(f.Dist, []*HNSWNode{newNode}, []S{newNodeVector}, neighborVector)
	nearest := NewMinMaxHeap(func(a, b *NodeAndVector[S, F]) bool {
		return f.Dist(a.Vector, neighborVector) < f.Dist(b.Vector, neighborVector)
	}, []*NodeAndVector[S, F]{
		{
			Node:   newNode,
			Vector: newNodeVector,
		},
	})

	currentNeighborNodes, err := f.NodeStore.BatchGet(neighborNode.GetNeighbors(layer))
	if err != nil {
		return nil, err
	}
	currentNeighborVectors, err := f.VectorStore.BatchGet(neighborNode.GetNeighbors(layer))
	if err != nil {
		return nil, err
	}

	for i := range currentNeighborNodes {
		nearest.Push(&NodeAndVector[S, F]{Node: currentNeighborNodes[i], Vector: currentNeighborVectors[i]})
	}

	SelectNeighborsSimple(nearest, f.M, layer)

	updatedNeighbors := make([]NodeID, 0)

	for _, nn := range nearest.Data {
		updatedNeighbors = append(updatedNeighbors, nn.Node.ID)
	}

	neighborNode.NeighborsByLevel[layer] = updatedNeighbors
	return neighborNode, nil
}

func NewHNSWGraph[S ~[]F, F constraints.Float](D int, M int, distanceFunc DistanceFunc[S, F], EfConstruction int) *HNSWGraph[S, F] {
	d := D
	m := M
	efConstruction := EfConstruction
	levelMult := 1 / math32.Log(float32(m))
	probabilities, neighborsPerLevel := calculateProbabilitiesAndNeighborsPerLevel(m, levelMult)
	rng := rand.New(rand.NewSource(0))
	return &HNSWGraph[S, F]{
		D:                  d,
		EfConstruction:     efConstruction,
		Dist:               distanceFunc,
		M:                  m,
		LevelMult:          levelMult,
		LevelProbabilities: probabilities,
		LevelNeighbors:     neighborsPerLevel,
		Rng:                rng,
		NodeStore:          NewSliceNodeStore(),
		VectorStore:        NewSliceVectorStore[S, F](),
	}
}

func calculateProbabilitiesAndNeighborsPerLevel(m int, levelMult float32) (probabilities []float32, neighborsPerLevel []int) {
	neighbors := 0
	neighborsPerLevel = make([]int, 0)
	probabilities = make([]float32, 0)
	for level := 0; ; level++ {
		levelProb := math32.Exp(-float32(level)/levelMult) * (1 - math32.Exp(-1/levelMult))
		if levelProb < 1e-9 {
			break
		}
		neighbors += m
		if level == 0 {
			neighbors += m
		}
		neighborsPerLevel = append(neighborsPerLevel, neighbors)
		probabilities = append(probabilities, levelProb)
	}
	return probabilities, neighborsPerLevel
}

func selectLevel(probabilities []float32, rng *rand.Rand) int {
	f := rng.Float32()
	for level, probability := range probabilities {
		if f < probability {
			return level
		}
		f -= probability
	}
	return len(probabilities) - 1
}

type VistedSet interface {
	Contains(id NodeID) bool
	Add(id NodeID)
}

type MapVistedSet struct {
	Nodes map[NodeID]bool
}

func NewMapVistedSet(nodes []NodeID) *MapVistedSet {
	m := make(map[NodeID]bool)
	for _, node := range nodes {
		m[node] = true
	}
	return &MapVistedSet{
		Nodes: m,
	}
}

func (m *MapVistedSet) Contains(id NodeID) bool {
	_, ok := m.Nodes[id]
	return ok
}

func (m *MapVistedSet) Add(id NodeID) {
	m.Nodes[id] = true
}

type BitsetVistedSet struct {
	bitset.BitSet
}

func NewBitsetVistedSet(nodes []NodeID) *BitsetVistedSet {
	b := bitset.New(uint(len(nodes)))
	for _, node := range nodes {
		b.Set(uint(node))
	}
	return &BitsetVistedSet{
		BitSet: *b,
	}
}

func (b *BitsetVistedSet) Contains(id NodeID) bool {
	return b.Test(uint(id))
}

func (b *BitsetVistedSet) Add(id NodeID) {
	b.Set(uint(id))
}

func Zip[S ~[]F, F constraints.Float](nodes []*HNSWNode, vectors []S) []*NodeAndVector[S, F] {
	nv := make([]*NodeAndVector[S, F], len(nodes))
	for i, node := range nodes {
		nv[i] = &NodeAndVector[S, F]{
			Node:   node,
			Vector: vectors[i],
		}
	}
	return nv
}

// SearchLayer searches layer `layerNumber` of the HNSWGraph for the `ef` nearest neighbors of `query` starting from
// the provided set of `entryPoints`. Returns the indexes of the nearest neighbors.
func SearchLayer[S ~[]F, F constraints.Float](
	hnsw *HNSWGraph[S, F],
	query S,
	entryPoints []NodeID,
	ef int,
	layerNumber int,
) (*minMaxHeap[*NodeAndVector[S, F]], error) {
	// fmt.Println("SearchLayer", entryPoints, layerNumber)

	visted := NewBitsetVistedSet(entryPoints)
	vectors, err := hnsw.VectorStore.BatchGet(entryPoints)
	if err != nil {
		return nil, err
	}
	nodes, err := hnsw.NodeStore.BatchGet(entryPoints)
	if err != nil {
		return nil, err
	}
	data := Zip[S, F](nodes, vectors)
	// candidates := MinDistQueue(hnsw.Dist, nodes, vectors, query)
	candidates := NewMinMaxHeap(func(a, b *NodeAndVector[S, F]) bool {
		return hnsw.Dist(a.Vector, query) < hnsw.Dist(b.Vector, query)
	}, data)
	// nearestNeighbors := MaxDistQueue(hnsw.Dist, nodes, vectors, query)
	nearestNeighbors := NewMinMaxHeap(func(a, b *NodeAndVector[S, F]) bool {
		return hnsw.Dist(a.Vector, query) < hnsw.Dist(b.Vector, query)
	}, data)

	for candidates.Len() > 0 {
		nearestCandidate := candidates.PopMin()
		furthestFound := nearestNeighbors.PeekMax()
		// fmt.Println("nearestCandidate", nearestCandidate.Node.ID, "furthestFound", furthestFound.Node.ID)

		if hnsw.Dist(nearestCandidate.Vector, query) > hnsw.Dist(furthestFound.Vector, query) {
			return nearestNeighbors, nil
		}

		for _, neighborID := range nearestCandidate.Node.GetNeighbors(layerNumber) {
			if !visted.Contains(neighborID) {
				visted.Add(neighborID)
				furthestFound = nearestNeighbors.PeekMax()
				// fmt.Println("set furthestFound", furthestFound.Node.ID)
				neighborVector, err := hnsw.VectorStore.Get(neighborID)
				if err != nil {
					return nil, err
				}
				neighborNode, err := hnsw.NodeStore.Get(neighborID)
				if err != nil {
					return nil, err
				}
				if nearestNeighbors.Len() < ef || hnsw.Dist(neighborVector, query) < hnsw.Dist(furthestFound.Vector, query) {
					nv := NodeAndVector[S, F]{Node: neighborNode, Vector: neighborVector}
					nearestNeighbors.Push(&nv)
					candidates.Push(&nv)
					if nearestNeighbors.Len() > ef {
						// fmt.Println("throwing away from ef", nearestNeighbors.PeekMax().Node.ID)
						nearestNeighbors.PopMax() // todo: maybe pop min?
					}
				}
			}
		}
	}
	nearestIDs := make([]NodeID, 0)
	for _, nodeAndVector := range nearestNeighbors.Data {
		nearestIDs = append(nearestIDs, nodeAndVector.Node.ID)
	}
	// fmt.Println("returning nearestNeighbors", nearestIDs)
	return nearestNeighbors, nil
}

func Insert[S ~[]F, F constraints.Float](hnsw *HNSWGraph[S, F], vector S) error {
	// fmt.Println("Insert", vector)

	insertLevel := selectLevel(hnsw.LevelProbabilities, hnsw.Rng)
	newNode := hnsw.EmptyNode()
	newNode.Layer = insertLevel

	// If the graph is empty, set the entry point to the new node, store it, and return
	if hnsw.ep == nil {
		hnsw.ep = &newNode.ID
		err := hnsw.Store(newNode, vector)
		if err != nil {
			return err
		}
		return nil
	}

	ep := *hnsw.ep
	epNode, err := hnsw.NodeStore.Get(ep)
	if err != nil {
		return err
	}
	topLayer := epNode.Layer

	// Search from top layer to insert level
	// var nearestNeighbors *DistancePriorityQueue[S, F]

	var nearestNeighbors *minMaxHeap[*NodeAndVector[S, F]]
	for layer := topLayer; layer > insertLevel; layer-- {
		nearestNeighbors, err = SearchLayer(hnsw, vector, []NodeID{ep}, 1, layer)
		if err != nil {
			return err
		}
		ep = nearestNeighbors.PeekMin().Node.ID
	}

	// Insert node at insert level and find all new neighbors
	eps := []NodeID{ep}
	for layer := min(insertLevel, topLayer); layer >= 0; layer-- {
		nearestNeighbors, err := SearchLayer(hnsw, vector, eps, hnsw.EfConstruction, layer)
		if err != nil {
			return err
		}

		// TODO: replace nearest with a double-ended priority queue.
		// nearest := slices.Clone(nearestNeighbors.PriorityQueue)
		// sort.Slice(nearest, func(i, j int) bool {
		// 	return hnsw.Dist(nearest[i].Vector, vector) < hnsw.Dist(nearest[j].Vector, vector)
		// })

		for nearestNeighbors.Len() > hnsw.M {
			// fmt.Println("throwing away ", nearestNeighbors.PeekMax().Node.ID)
			nearestNeighbors.PopMax()
		}

		// newEps := make([]NodeID, 0)
		for i := nearestNeighbors.Len(); i > 0; i-- {
			neighbor := nearestNeighbors.PopMin()
			// fmt.Println("adding ", neighbor.Node.ID)
			newNode.NeighborsByLevel[layer] = append(newNode.NeighborsByLevel[layer], neighbor.Node.ID)
			// newEps = append(newEps, neighbor.Node.ID)
		}
		// newEps := make([]NodeID, 0)
		// for _, neighbor := range nearest[:min(nearest.Len(), hnsw.M)] {
		// 	newNode.NeighborsByLevel[layer] = append(newNode.NeighborsByLevel[layer], neighbor.Node.ID)
		// 	newEps = append(newEps, neighbor.Node.ID)
		// }
		// eps = newEps
	}

	// Update neighbors of neighbors
	toUpdate := make([]*HNSWNode, 0)
	for layer := min(insertLevel, topLayer); layer >= 0; layer-- {
		neighbors := newNode.GetNeighbors(layer)
		for _, neighborID := range neighbors {
			updatedNeighbor, err := hnsw.UpdateNeighbor(neighborID, newNode, vector, layer)
			if err != nil {
				return err
			}
			toUpdate = append(toUpdate, updatedNeighbor)
		}
	}

	//fmt.Println("Insert", newNode)
	err = hnsw.Store(newNode, vector)
	if err != nil {
		return err
	}
	if insertLevel > topLayer {
		hnsw.ep = &newNode.ID
	}
	for _, node := range toUpdate {
		err := hnsw.NodeStore.Add(node.ID, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func SelectNeighborsSimple[S ~[]F, F constraints.Float](neighbors *minMaxHeap[*NodeAndVector[S, F]], m int, layer int) {
	mLayer := m
	if layer == 0 {
		mLayer += m
	}
	for neighbors.Len() > mLayer {
		neighbors.PopMax()
	}
}

// Search searches the HNSWGraph for the `K` nearest neighbors of `query`. Returns the indexes of the nearest neighbors.
func Search[S ~[]F, F constraints.Float](hnsw *HNSWGraph[S, F], query S, K int, ef int) (*minMaxHeap[*NodeAndVector[S, F]], error) {
	ep := hnsw.EntryPoint()
	// nearestFound := MinDistQueue(hnsw.Dist, []*HNSWNode{}, []S{}, query)
	nearestFound := NewMinMaxHeap(func(a, b *NodeAndVector[S, F]) bool {
		return hnsw.Dist(a.Vector, query) < hnsw.Dist(b.Vector, query)
	}, []*NodeAndVector[S, F]{})
	epNode, err := hnsw.NodeStore.Get(ep)
	if err != nil {
		return nil, err
	}
	maxLayer := epNode.Layer

	for layer := maxLayer; layer >= 0; layer-- {
		nearestFound, err = SearchLayer(hnsw, query, []NodeID{ep}, ef, layer)
		if err != nil {
			return nil, err
		}
		ep = nearestFound.PeekMin().Node.ID
	}
	for nearestFound.Len() > K {
		nearestFound.PopMax()
	}
	return nearestFound, nil
}

// BruteForceSearch searches the HNSWGraph for the `K` nearest neighbors of `query`. Returns the indexes of the nearest neighbors.
func BruteForceSearch[S ~[]F, F constraints.Float](hnsw *HNSWGraph[S, F], query S, K int) ([]*NodeAndVector[S, F], error) {
	// nearestFound := MinDistQueue(hnsw.Dist, []*HNSWNode{}, []S{}, query)
	nearestFound := NewMinMaxHeap(func(a, b *NodeAndVector[S, F]) bool {
		return hnsw.Dist(a.Vector, query) < hnsw.Dist(b.Vector, query)
	}, []*NodeAndVector[S, F]{})

	n := hnsw.NodeStore.NextID()
	for i := NodeID(0); i < n; i++ {
		node, err := hnsw.NodeStore.Get(i)
		if err != nil {
			return nil, err
		}
		vector, err := hnsw.VectorStore.Get(i)
		if err != nil {
			return nil, err
		}
		nearestFound.Push(&NodeAndVector[S, F]{Node: node, Vector: vector})
	}
	found := make([]*NodeAndVector[S, F], 0)
	foundLen := nearestFound.Len()
	for i := 0; i < min(K, foundLen); i++ {
		popped := nearestFound.PopMin()
		found = append(found, popped)
	}
	return found, nil
}
