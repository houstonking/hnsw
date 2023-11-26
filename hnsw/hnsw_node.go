package hnsw

type HNSWNode struct {
	ID               NodeID
	Layer            int
	NeighborsByLevel [][]NodeID
}

func (n *HNSWNode) GetNeighbors(level int) []NodeID {
	return n.NeighborsByLevel[level]
}

// AddNeighbor adds a neighbor to the node at the given level if it does not already exist.
func (n *HNSWNode) AddNeighbor(level int, neighbor NodeID) {
	for _, node := range n.NeighborsByLevel[level] {
		if node == neighbor {
			return
		}
	}
	n.NeighborsByLevel[level] = append(n.NeighborsByLevel[level], neighbor)
}
