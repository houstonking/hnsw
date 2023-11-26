package hnsw

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

type VectorStore[S ~[]F, F constraints.Float] interface {
	Add(id NodeID, vector S) error
	Get(id NodeID) (S, error)
	BatchGet(ids []NodeID) ([]S, error)
}

type MapVectorStore[S ~[]F, F constraints.Float] struct {
	Vectors map[NodeID]S
}

func NewMapVectorStore[S ~[]F, F constraints.Float]() *MapVectorStore[S, F] {
	return &MapVectorStore[S, F]{
		Vectors: make(map[NodeID]S),
	}
}

func (m *MapVectorStore[S, F]) Add(id NodeID, vector S) error {
	//fmt.Println("MapVectorStore.Add", id)
	m.Vectors[id] = vector
	return nil
}

func (m *MapVectorStore[S, F]) Get(id NodeID) (S, error) {
	vector, ok := m.Vectors[id]
	if !ok {
		return nil, fmt.Errorf("vector %d not found", id)
	}
	return vector, nil
}

func (m *MapVectorStore[S, F]) BatchGet(ids []NodeID) ([]S, error) {
	vectors := make([]S, len(ids))
	for i, id := range ids {
		vector, ok := m.Vectors[id]
		if !ok {
			return nil, fmt.Errorf("vector %d not found", id)
		}
		vectors[i] = vector
	}
	return vectors, nil
}

type SliceVectorStore[S ~[]F, F constraints.Float] struct {
	Vectors []S
}

func NewSliceVectorStore[S ~[]F, F constraints.Float]() *SliceVectorStore[S, F] {
	return &SliceVectorStore[S, F]{
		Vectors: make([]S, 0),
	}
}

func (s *SliceVectorStore[S, F]) Add(id NodeID, vector S) error {
	if id != NodeID(len(s.Vectors)) {
		return fmt.Errorf("invalid id %d", id)
	}
	s.Vectors = append(s.Vectors, vector)
	return nil
}

func (s *SliceVectorStore[S, F]) Get(id NodeID) (S, error) {
	if id >= NodeID(len(s.Vectors)) {
		return nil, fmt.Errorf("vector %d not found", id)
	}
	return s.Vectors[id], nil
}

func (s *SliceVectorStore[S, F]) BatchGet(ids []NodeID) ([]S, error) {
	vectors := make([]S, len(ids))
	for i, id := range ids {
		if id >= NodeID(len(s.Vectors)) {
			return nil, fmt.Errorf("vector %d not found", id)
		}
		vectors[i] = s.Vectors[id]
	}
	return vectors, nil
}
