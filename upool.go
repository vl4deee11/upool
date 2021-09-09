package upool

import (
	"unsafe"
)

type Test struct {
	b  string
	bx string
	A  int8
	b2 string
	b3 string
	b4 string
	C  float64
	B  []byte
}

// TODO: use code generation like c++ template class
const maxInChunks = 10
const size = int(unsafe.Sizeof(Test{}))

type UPool struct {
	New        func() Test
	currChunk  int
	currOffset int
	freeptrs   []uintptr
	memChunks  [][maxInChunks * size]byte
}

func NewUPool(new func() Test, pSize int) *UPool {
	return &UPool{
		New:       new,
		memChunks: make([][maxInChunks * size]byte, 0, pSize),
		freeptrs:  make([]uintptr, 0, pSize),
	}
}

// noescape: unused now
func (p *UPool) noescape(up unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(up) ^ 0)
}

// uintptr2t: uintptr to type convert
func (p *UPool) uintptr2t(ptr uintptr) *Test {
	return (*Test)(unsafe.Pointer(ptr))
}

// t2uintptr: type to uintptr convert
func (p *UPool) t2uintptr(v *Test) uintptr {
	return uintptr(unsafe.Pointer(v))
}

// malloc: allocate new memory chunk
func (p *UPool) malloc() {
	p.memChunks = append(p.memChunks, [maxInChunks * size]byte{})
}

// Get struct from pool
func (p *UPool) Get() *Test {
	if len(p.freeptrs) != 0 {
		ptr := p.freeptrs[len(p.freeptrs)-1]
		p.freeptrs = p.freeptrs[:len(p.freeptrs)-1]
		return p.uintptr2t(ptr)
	}

	st := p.New()
	if len(p.memChunks) == 0 {
		p.malloc()
	}

	if p.currOffset == len(p.memChunks[p.currChunk]) {
		p.malloc()
		p.currOffset = 0
		p.currChunk++
	}

	ptr := unsafe.Pointer(&st)

	bs := *(*[size]byte)(ptr)
	for i := range bs {
		p.memChunks[p.currChunk][p.currOffset+i] = bs[i]
	}
	p.currOffset += size
	return (*Test)(unsafe.Pointer(&p.memChunks[p.currChunk][p.currOffset-size]))
}


//Return struct to pool
func (p *UPool) Return(st *Test) {
	p.freeptrs = append(p.freeptrs, p.t2uintptr(st))
}
