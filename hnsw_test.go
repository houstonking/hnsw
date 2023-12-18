package hnsw

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/chewxy/math32"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

/*
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
			NodeStore:          NewMapNodeStore(),
			VectorStore:        NewMapVectorStore[S, F](),
		}
	}
*/
func TestHNSWGraph(t *testing.T) {
	hnsw := NewHNSWGraph(1536, 32, CosineDistance32, 100)

	// populate with random vectors
	rng := rand.New(rand.NewSource(12345))
	for i := 0; i < 5000; i++ {
		vector := make([]float32, hnsw.D)
		for j := range vector {
			vector[j] = rng.Float32()
		}

		err := Insert(hnsw, vector)
		assert.NoError(t, err)
	}
}

//func BenchmarkHNSWGraph_Search(b *testing.B) {
//	dim := 2048
//	hnsw := NewHNSWGraph(dim, 32, CosineDistance32, 100)
//	n := 10000
//	// generate random vectors
//	rng := rand.New(rand.NewSource(12345))
//
//	fmt.Println("Generating random vectors")
//	vecs := generateRandomVectors(rng, n, dim)
//
//	fmt.Println("Inserting vectors")
//	for _, vec := range vecs {
//		err := Insert(hnsw, vec)
//		assert.NoError(b, err)
//	}
//
//	fmt.Println("Searching")
//	b.ResetTimer()
//	b.Run("Search", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			vec := generateRandomVectors(rng, 1, dim)[0]
//			_, err := Search(hnsw, vec, 10, 100)
//			assert.NoError(b, err)
//		}
//	})
//}

func BenchmarkHNSWGraph_Insert(b *testing.B) {
	dims := []int{8}
	//dims := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512}
	//ns := []int{1, 10, 100, 1_000, 10_000, 100_000}
	// ns := []int{1_000, 2_000, 3_000, 4_000, 5_000, 6_000, 7_000, 8_000, 9_000, 10_000}
	ns := []int{1_000, 2_000, 4_000, 8_000}
	for _, dim := range dims {
		for _, n := range ns {
			b.Run(fmt.Sprintf("dim=%d, n=%d", dim, n), func(b *testing.B) {
				hnsw := NewHNSWGraph(dim, 32, CosineDistance32, 100)
				rng := rand.New(rand.NewSource(12345))
				vecs := generateRandomVectors(rng, n, dim)
				for _, vec := range vecs[:n-1] {
					err := Insert(hnsw, vec)
					assert.NoError(b, err)
				}
				RunInsert(hnsw, vecs[0], b)
			})
		}
	}
}

func BenchmarkHNSWGraph_Search(b *testing.B) {
	dims := []int{1536}
	//dims := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512}
	// ns := []int{1, 10, 100, 1_000, 10_000, 100_000, 1_000_000}
	// ns := []int{1000}
	ns := []int{1_000, 2_000, 3_000, 4_000, 5_000, 6_000, 7_000, 8_000, 9_000, 10_000}
	// ns := []int{1_000, 2_000, 4_000, 8_000}

	// efSearches := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	efSearches := []int{64}
	for _, dim := range dims {
		for _, n := range ns {
			for _, efSearch := range efSearches {
				hnsw := NewHNSWGraph(dim, 32, L2Distance32, 100)
				rng := rand.New(rand.NewSource(time.Now().UnixNano()))
				vecs := generateRandomVectors(rng, n, dim)
				for _, vec := range vecs {
					err := Insert(hnsw, vec)
					assert.NoError(b, err)
				}
				RunSearch(hnsw, 10, efSearch, b)

			}
		}
	}
}

func TestDeterministicHNSWGraph_Search(t *testing.T) {
	vecs := [][]float32{}
	for i := 0; i < 100; i++ {
		vecs = append(vecs, []float32{float32(i)})
	}
	hnsw := NewHNSWGraph(1, 5, L2Distance32, 100)
	for _, vec := range vecs[:len(vecs)-1] {
		err := Insert(hnsw, vec)
		assert.NoError(t, err)
	}

	err := Insert(hnsw, vecs[len(vecs)-1])
	assert.NoError(t, err)
	fmt.Println("Searching")

	for _, vec := range vecs {
		res, err := Search(hnsw, vec, 10, 100)
		assert.NoError(t, err)
		println("\nvec", int(vec[0]))
		for res.Len() > 0 {
			nodeAndVector := res.PopMin()
			dist := hnsw.Dist(nodeAndVector.Vector, vec)
			t.Logf("dist: %d, id: %d", int(dist), nodeAndVector.Node.ID)
		}
	}
}

func TestDeterministicHNSWGraph_Search2D(t *testing.T) {
	vecs := [][]float32{}
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			vecs = append(vecs, []float32{float32(i), float32(j)})
		}
	}
	vecs = lo.Shuffle(vecs)

	hnsw := NewHNSWGraph(1, 5, L2Distance32, 100)
	for _, vec := range vecs[:len(vecs)-1] {
		err := Insert(hnsw, vec)
		assert.NoError(t, err)
	}

	err := Insert(hnsw, vecs[len(vecs)-1])
	assert.NoError(t, err)
	fmt.Println("Searching")

	for _, vec := range vecs {
		res, err := Search(hnsw, vec, 10, 100)
		assert.NoError(t, err)
		t.Logf("\nvec: %.2f, %.2f", vec[0], vec[1])
		for res.Len() > 0 {
			nodeAndVector := res.PopMin()
			dist := hnsw.Dist(nodeAndVector.Vector, vec)
			t.Logf("dist: %.2f, vec: %.2f, %.2f", dist, nodeAndVector.Vector[0], nodeAndVector.Vector[1])
		}
	}
}

func RunSearch(hnsw *HNSWGraph[[]float32, float32], k int, ef int, b *testing.B) {
	var err error
	var searchRes *minMaxHeap[*NodeAndVector[[]float32, float32]]
	var lastVec []float32
	b.Run(fmt.Sprintf("Search dim= %d, n= %d, k=%d, ef=%d", hnsw.D, hnsw.NodeStore.NextID(), k, ef), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			vec := generateRandomVectors(rand.New(rand.NewSource(time.Now().UnixNano())), 1, hnsw.D)[0]
			lastVec = vec
			searchRes, err = Search(hnsw, vec, k, ef)
			assert.NoError(b, err)
		}
	})

	before := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			randVec := generateRandomVectors(rand.New(rand.NewSource(time.Now().UnixNano())), 1, hnsw.D)[0]
			_, err := Search(hnsw, randVec, k, ef)
			assert.NoError(b, err)
			wg.Done()
		}()
	}
	wg.Wait()
	after := time.Since(before)
	b.Logf("Parallel search time: %v", after)

	if hnsw.D < 10 {
		b.Logf("queryVec: %+v", lastVec)
	}
	data := make([]*NodeAndVector[[]float32, float32], searchRes.Len())
	searchLen := searchRes.Len()
	for i := 0; i < min(k, searchLen); i++ {
		nodeAndVector := searchRes.PopMin()
		data[i] = nodeAndVector
	}
	// sort by distance
	sort.Slice(data, func(i, j int) bool {
		return hnsw.Dist(data[i].Vector, lastVec) < hnsw.Dist(data[j].Vector, lastVec)
	})
	for _, nodeAndVector := range data {
		dist := hnsw.Dist(nodeAndVector.Vector, lastVec)

		if hnsw.D < 10 {
			b.Logf("dist: %.4f, id: %d, vec: %+v", dist, nodeAndVector.Node.ID, nodeAndVector.Vector)
		} else {
			b.Logf("dist: %.4f, id: %d", dist, nodeAndVector.Node.ID)
		}
	}
	b.Logf("Brute force search results")
	bfRes, err := BruteForceSearch(hnsw, lastVec, k)
	assert.NoError(b, err)
	for _, nodeAndVector := range bfRes {
		dist := hnsw.Dist(nodeAndVector.Vector, lastVec)

		if hnsw.D < 10 {
			b.Logf("dist: %.4f, id: %d, vec: %+v", dist, nodeAndVector.Node.ID, nodeAndVector.Vector)
		} else {
			b.Logf("dist: %.4f, id: %d", dist, nodeAndVector.Node.ID)
		}
	}
}

func RunInsert(hnsw *HNSWGraph[[]float32, float32], vec []float32, b *testing.B) {
	b.Run("Insert %d", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := Insert(hnsw, vec)
			assert.NoError(b, err)
		}
	})
}

func generateRandomVectors(rng *rand.Rand, n int, dim int) [][]float32 {
	vecs := make([][]float32, n)
	for i := range vecs {
		vecs[i] = make([]float32, dim)
		for j := range vecs[i] {
			vecs[i][j] = rng.Float32()
		}
	}
	return vecs

}

func Test_calculateProbabilitiesAndNeighborsPerLevel(t *testing.T) {
	type args struct {
		M  int
		ml float32
	}
	tests := []struct {
		name                  string
		args                  args
		wantProbabilities     []float32
		wantNeighborsPerLevel []int
	}{
		{
			name: "M=32, ml=1 / log(M)",
			args: args{
				M:  32,
				ml: 1 / math32.Log(32),
			},
			wantProbabilities: []float32{
				0.96875,
				0.030273437499999986,
				0.0009460449218749991,
				2.956390380859371e-05,
				9.23871994018553e-07,
				2.887099981307982e-08,
			},
			wantNeighborsPerLevel: []int{64, 96, 128, 160, 192, 224},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProbabilities, gotNeighborsPerLevel := calculateProbabilitiesAndNeighborsPerLevel(tt.args.M, tt.args.ml)
			assert.InDeltaSlice(t, tt.wantProbabilities, gotProbabilities, 0.0001, "calculateProbabilitiesAndNeighborsPerLevel(%v, %v)", tt.args.M, tt.args.ml)
			assert.Equalf(t, tt.wantNeighborsPerLevel, gotNeighborsPerLevel, "calculateProbabilitiesAndNeighborsPerLevel(%v, %v)", tt.args.M, tt.args.ml)
		})
	}
}

func Test_selectLevelDistribution(t *testing.T) {
	type args struct {
		probabilities []float32
		rng           *rand.Rand
	}
	tests := []struct {
		name             string
		args             args
		numSelections    int
		wantDistribution []int
	}{
		{
			name: "selectLevel",
			args: args{
				probabilities: []float32{
					0.96875,
					0.030273437499999986,
					0.0009460449218749991,
					2.956390380859371e-05,
					9.23871994018553e-07,
					2.887099981307982e-08,
				},
				rng: rand.New(rand.NewSource(12345)), // deterministic seed
			},
			numSelections:    100000,
			wantDistribution: []int{96905, 2972, 116, 6, 1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distribution := make([]int, len(tt.args.probabilities))
			for i := 0; i < tt.numSelections; i++ {
				level := selectLevel(tt.args.probabilities, tt.args.rng)
				distribution[level]++
			}
			assert.Equalf(t, tt.wantDistribution, distribution, "selectLevelDistribution(%v, %v)", tt.args.probabilities, tt.args.rng)
		})
	}
}
