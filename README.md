## HNSW 

I mostly just wanted to knock out a quick implementation of a HNSW as a side-project. 
While mostly unoptimized, the asymptotics seem about right and distance calculations on larger dimensional vectors dominate.

https://arxiv.org/abs/1603.09320

### Benchmark Results

```
Running tool: /snap/bin/go test -benchmem -run=^$ -bench ^BenchmarkHNSWGraph_Search$ hnsw/hnsw

goos: linux
goarch: amd64
pkg: hnsw/hnsw
cpu: AMD Ryzen 9 7950X 16-Core Processor
=== RUN   BenchmarkHNSWGraph_Search
BenchmarkHNSWGraph_Search
=== RUN   BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_1,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_1,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_1,_k=10,_ef=64-32                  103090             10903 ns/op           12008 B/op         19 allocs/op
    /home/king/git/hnsw/hnsw/hnsw_test.go:210: Parallel search time: 9.737899ms
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 16.2566, id: 0
    /home/king/git/hnsw/hnsw/hnsw_test.go:232: Brute force search results
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 16.2566, id: 0
=== RUN   BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_10,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_10,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_10,_k=10,_ef=64-32                  86139             13900 ns/op           13352 B/op         45 allocs/op
    /home/king/git/hnsw/hnsw/hnsw_test.go:210: Parallel search time: 14.16301ms
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.6609, id: 5
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.7277, id: 6
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.8561, id: 2
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.8709, id: 3
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.8881, id: 9
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.9571, id: 7
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.9837, id: 0
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 16.0239, id: 8
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 16.0548, id: 1
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 16.1674, id: 4
    /home/king/git/hnsw/hnsw/hnsw_test.go:232: Brute force search results
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.6609, id: 5
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.7277, id: 6
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.8561, id: 2
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.8709, id: 3
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.8881, id: 9
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.9571, id: 7
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.9837, id: 0
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 16.0239, id: 8
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 16.0548, id: 1
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 16.1674, id: 4
=== RUN   BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_100,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_100,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_100,_k=10,_ef=64-32                 21238             56626 ns/op           25101 B/op        227 allocs/op
    /home/king/git/hnsw/hnsw/hnsw_test.go:210: Parallel search time: 43.707259ms
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.5514, id: 61
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.5546, id: 11
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.5815, id: 64
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.5895, id: 22
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.6361, id: 40
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.7063, id: 90
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.7128, id: 42
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.7364, id: 63
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.7634, id: 15
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.7741, id: 5
    /home/king/git/hnsw/hnsw/hnsw_test.go:232: Brute force search results
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.5514, id: 61
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.5546, id: 11
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.5815, id: 64
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.5895, id: 22
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.6361, id: 40
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.7063, id: 90
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.7128, id: 42
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.7364, id: 63
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.7634, id: 15
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.7741, id: 5
=== RUN   BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_1000,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_1000,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_1000,_k=10,_ef=64-32                 5150            228420 ns/op           44995 B/op        571 allocs/op
    /home/king/git/hnsw/hnsw/hnsw_test.go:210: Parallel search time: 62.65939ms
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2246, id: 791
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2361, id: 563
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.3534, id: 854
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.3711, id: 878
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.4081, id: 714
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.4184, id: 778
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.4262, id: 317
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.4634, id: 735
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.4647, id: 251
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.4736, id: 539
    /home/king/git/hnsw/hnsw/hnsw_test.go:232: Brute force search results
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.2246, id: 791
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.2361, id: 563
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.3534, id: 854
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.3711, id: 878
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.4081, id: 714
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.4184, id: 778
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.4262, id: 317
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.4634, id: 735
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.4647, id: 251
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.4736, id: 539
=== RUN   BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_10000,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_10000,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_10000,_k=10,_ef=64-32                1146           1025654 ns/op           61235 B/op        773 allocs/op
    /home/king/git/hnsw/hnsw/hnsw_test.go:210: Parallel search time: 101.708914ms
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2134, id: 2850
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2531, id: 284
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2863, id: 1441
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.3609, id: 212
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.3713, id: 3015
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.3745, id: 4586
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.3815, id: 1003
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.3927, id: 1687
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.4021, id: 1961
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.4156, id: 6749
    /home/king/git/hnsw/hnsw/hnsw_test.go:232: Brute force search results
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.0918, id: 6960
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.2134, id: 2850
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.2145, id: 7935
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.2531, id: 284
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.2688, id: 4302
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.2723, id: 7977
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.2863, id: 1441
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.3066, id: 5115
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.3092, id: 8423
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.3108, id: 4261
=== RUN   BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_100000,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_100000,_k=10,_ef=64
BenchmarkHNSWGraph_Search/Search_dim=_1536,_n=_100000,_k=10,_ef=64-32                640           1680860 ns/op          124862 B/op         901 allocs/op
    /home/king/git/hnsw/hnsw/hnsw_test.go:210: Parallel search time: 264.365378ms
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.1807, id: 5938
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.1910, id: 14867
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2080, id: 21314
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2108, id: 8142
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2188, id: 1040
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2245, id: 1961
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2248, id: 5829
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2464, id: 7942
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2558, id: 1594
    /home/king/git/hnsw/hnsw/hnsw_test.go:229: dist: 15.2602, id: 4518
    /home/king/git/hnsw/hnsw/hnsw_test.go:232: Brute force search results
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 14.9686, id: 47518
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 14.9693, id: 4040
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.0154, id: 30650
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.0232, id: 85075
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.0369, id: 62532
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.0416, id: 64564
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.0432, id: 18721
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.0458, id: 79555
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.0516, id: 82357
    /home/king/git/hnsw/hnsw/hnsw_test.go:239: dist: 15.0587, id: 85868
PASS
ok      hnsw/hnsw       294.786s
```