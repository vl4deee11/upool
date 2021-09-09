package upool

import "testing"

var x *Test

// go test -bench=. -gcflags '-l -N' -benchmem -cpu=1
// goos: linux
// goarch: amd64
// Benchmark_PoolGetOnly           10000000               759.1 ns/op           627 B/op          1 allocs/op
// Benchmark_PoolGetReturn         10000000                12.23 ns/op            0 B/op          0 allocs/op
func Benchmark_PoolGetOnly(b *testing.B) {
	f := func() Test {
		te := Test{}
		te.A = 123
		te.b = "NOTUSED"
		te.B = []byte{1, 2, 3, 4}
		return te
	}

	var v *Test
	p := NewUPool(f, 1000)
	for j := 0; j < b.N; j++ {
		v = p.Get()
	}
	x = v
}

func Benchmark_PoolGetReturn(b *testing.B) {
	f := func() Test {
		te := Test{}
		te.A = 123
		te.b = "NOTUSED"
		te.B = []byte{1, 2, 3, 4}
		return te
	}

	var v *Test
	p := NewUPool(f, 1000)
	for _i := 0; _i < b.N; _i++ {
		v = p.Get()
		p.Return(v)
	}
	x = v
}

func Test_Pool(t *testing.T) {
	f := func() Test {
		te := Test{}
		te.A = 123
		te.b = "NOTUSED"
		te.B = []byte{1, 2, 3, 4}
		return te
	}

	p := NewUPool(f, 10)
	var objects []*Test
	for i := 0; i < 10000; i++ {
		res := p.Get()
		resV := res
		objects = append(objects, resV)
		if i%2 == 0 {
			p.Return(objects[i])
		}
	}
}
