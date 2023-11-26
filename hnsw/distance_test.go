package hnsw

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
	"testing"
)

func SF32(fs ...float32) []float32 {
	return fs
}

func toArgs[S ~[]F, F constraints.Float](a, b S) args[S, F] {
	return args[S, F]{a: a, b: b}
}

func f32Args(a, b []float32) args[[]float32, float32] {
	return toArgs(a, b)
}

type args[S ~[]F, F constraints.Float] struct {
	a S
	b S
}

func TestCosineSimilarityAndDistance(t *testing.T) {

	type testCase[S ~[]F, F constraints.Float] struct {
		name           string
		args           args[S, F]
		wantSimilarity F
		wantDistance   F
	}
	tests := []testCase[[]float32, float32]{
		{
			name:           "same 2D vector",
			args:           f32Args(SF32(0.1, 0.9), SF32(0.1, 0.9)),
			wantSimilarity: 1,
			wantDistance:   0,
		},
		{
			name:           "orthogonal 2D vector",
			args:           f32Args(SF32(1, 0), SF32(0, 1)),
			wantSimilarity: 0,
			wantDistance:   1,
		},
		{
			name:           "opposite 2D vector",
			args:           f32Args(SF32(1, 0), SF32(-1, 0)),
			wantSimilarity: -1,
			wantDistance:   2,
		},
		{
			name:           "same 10D vector",
			args:           f32Args(SF32(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), SF32(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
			wantSimilarity: 1,
			wantDistance:   0,
		},
		{
			name:           "orthogonal 10D vector",
			args:           f32Args(SF32(1, 0, 0, 0, 0, 0, 0, 0, 0, 0), SF32(0, 1, 0, 0, 0, 0, 0, 0, 0, 0)),
			wantSimilarity: 0,
			wantDistance:   1,
		},
		{
			name:           "opposite 10D vector",
			args:           f32Args(SF32(1, 0, 0, 0, 0, 0, 0, 0, 0, 0), SF32(-1, 0, 0, 0, 0, 0, 0, 0, 0, 0)),
			wantSimilarity: -1,
			wantDistance:   2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CosineSimilarity32(tt.args.a, tt.args.b)
			assert.InDelta(t, tt.wantSimilarity, got, 0.0001)
			got = CosineDistance32(tt.args.a, tt.args.b)
			assert.InDelta(t, tt.wantDistance, got, 0.0001)
		})
	}
}
