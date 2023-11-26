package hnsw

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
	"testing"
)

func TestKNN(t *testing.T) {
	type args[S interface{ ~[]F }, F constraints.Float] struct {
		candidates []S
		D          DistanceFunc[S, F]
		q          S
		k          int
	}
	type testCase[S interface{ ~[]F }, F constraints.Float] struct {
		name string
		args args[S, F]
		want []int
	}
	tests := []testCase[[]float32, float32]{
		{
			name: "same 2D vector",
			args: args[[]float32, float32]{
				candidates: [][]float32{
					{0.1, 0.9},
					{0.1, 0.9},
				},
				D: CosineDistance32,
				q: []float32{0.100001, 0.9}, // slight offset to make tie-breaks deterministic in top k order
				k: 2,
			},
			want: []int{1, 0},
		},
		{
			name: "orthogonal 2D vector",
			args: args[[]float32, float32]{
				candidates: [][]float32{
					{1, 0},
					{0, 1},
				},
				D: CosineDistance32,
				q: []float32{1, 0},
				k: 1,
			},
			want: []int{0},
		},
		{
			name: "opposite 2D vector",
			args: args[[]float32, float32]{
				candidates: [][]float32{
					{1, 0},
					{-1, 0},
				},
				D: CosineDistance32,
				q: []float32{1, 0},
				k: 1,
			},
			want: []int{0},
		},
		{
			name: "10 2D vectors, between {1, 0} and {0, 1}",
			args: args[[]float32, float32]{
				candidates: [][]float32{
					{1, 0},
					{0.9, 0.1},
					{0.8, 0.2},
					{0.7, 0.3},
					{0.6, 0.4},
					{0.5, 0.5},
					{0.4, 0.6},
					{0.3, 0.7},
					{0.2, 0.8},
					{0.1, 0.9},
				},
				D: CosineDistance32,
				q: []float32{0.0001, 1}, // slight offset to make tie-breaks deterministic in top k order
				k: 5,
			},
			want: []int{9, 8, 7, 6, 5},
		},
		{
			name: "10 2D vectors, between {1, 0} and {0, 1}",
			args: args[[]float32, float32]{
				candidates: [][]float32{
					{1, 0},
					{0.9, 0.1},
					{0.8, 0.2},
					{0.7, 0.3},
					{0.6, 0.4},
					{0.5, 0.5},
					{0.4, 0.6},
					{0.3, 0.7},
					{0.2, 0.8},
					{0.1, 0.9},
				},
				D: CosineDistance32,
				q: []float32{10001, 10000}, // slight offset to make tie-break deterministic in top k order
				k: 5,
			},
			want: []int{5, 4, 6, 3, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, KNN(tt.args.candidates, tt.args.D, tt.args.q, tt.args.k), "KNN(%v, %v, %v, %v)", tt.args.candidates, tt.args.D, tt.args.q, tt.args.k)
		})
	}
}
