package hnsw

import "fmt"

type NodeStore interface {
	Add(id NodeID, node *HNSWNode) error
	Get(id NodeID) (*HNSWNode, error)
	BatchGet(ids []NodeID) ([]*HNSWNode, error)
	NextID() NodeID
}

type MapNodeStore struct {
	Nodes map[NodeID]*HNSWNode
}

func NewMapNodeStore() *MapNodeStore {
	return &MapNodeStore{
		Nodes: make(map[NodeID]*HNSWNode),
	}
}

func (m *MapNodeStore) Add(id NodeID, node *HNSWNode) error {
	//fmt.Println("MapNodeStore.Add", id, node)
	m.Nodes[id] = node
	return nil
}

func (m *MapNodeStore) Get(id NodeID) (*HNSWNode, error) {
	node, ok := m.Nodes[id]
	if !ok {
		return nil, fmt.Errorf("node %d not found", id)
	}
	return node, nil
}

func (m *MapNodeStore) BatchGet(ids []NodeID) ([]*HNSWNode, error) {
	nodes := make([]*HNSWNode, len(ids))
	for i, id := range ids {
		node, ok := m.Nodes[id]
		if !ok {
			return nil, fmt.Errorf("node %d not found", id)
		}
		nodes[i] = node
	}
	return nodes, nil
}

func (m *MapNodeStore) NextID() NodeID {
	return NodeID(len(m.Nodes))
}

type SliceNodeStore struct {
	Nodes []*HNSWNode
}

func NewSliceNodeStore() *SliceNodeStore {
	return &SliceNodeStore{
		Nodes: make([]*HNSWNode, 0),
	}
}

func (s *SliceNodeStore) Add(id NodeID, node *HNSWNode) error {
	//fmt.Println("SliceNodeStore.Add", id, node)
	if id > NodeID(len(s.Nodes)) {
		return fmt.Errorf("node %d not found", id)
	}
	if id == NodeID(len(s.Nodes)) {
		s.Nodes = append(s.Nodes, node)
	}
	s.Nodes[id] = node
	return nil
}

func (s *SliceNodeStore) Get(id NodeID) (*HNSWNode, error) {
	if id >= NodeID(len(s.Nodes)) {
		return nil, fmt.Errorf("node %d not found", id)
	}
	return s.Nodes[id], nil
}

func (s *SliceNodeStore) BatchGet(ids []NodeID) ([]*HNSWNode, error) {
	nodes := make([]*HNSWNode, len(ids))
	for i, id := range ids {
		if id >= NodeID(len(s.Nodes)) {
			return nil, fmt.Errorf("node %d not found", id)
		}
		nodes[i] = s.Nodes[id]
	}
	return nodes, nil
}

func (s *SliceNodeStore) NextID() NodeID {
	return NodeID(len(s.Nodes))
}
