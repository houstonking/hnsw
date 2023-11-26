package hnsw

import (
	"github.com/viterin/vek/vek32"
	"golang.org/x/exp/constraints"
)

type DistanceFunc[S ~[]F, F constraints.Float] func(a, b S) F

func CosineSimilarity32(a, b []float32) float32 {
	return vek32.Dot(a, b) / (vek32.Norm(a) * vek32.Norm(b))
}

func CosineDistance32(a, b []float32) float32 {
	return 1 - CosineSimilarity32(a, b)
}

func Less[F constraints.Float](a F, b F) bool {
	return a < b
}

func GreaterEqual[F constraints.Float](a F, b F) bool {
	return a >= b
}

func L2Distance32(a, b []float32) float32 {
	return vek32.Distance(a, b)
}
