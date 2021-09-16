package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"strings"
)

var tmpl = `
// DANGER! This code was autogenerated from template by github.com/vl4deee11/upool/codegen/.
// For more info check github.com/vl4deee11/upool/codegen/main.go
package {{.Package}}

import (
	"unsafe"
)

const maxInChunks = {{.MaxChunks}}
const size = int(unsafe.Sizeof(new({{.Type}})))

type UPool{{Title .Type}} struct {
	New        func() {{.Type}}
	currChunk  int
	currOffset int
	freeptrs   []uintptr
	memChunks  [][maxInChunks * size]byte
}

func NewUPool{{Title .Type}}(new func() {{.Type}}, pSize int) *UPool{{Title .Type}} {
	return &UPool{{Title .Type}}{
		New:       new,
		memChunks: make([][maxInChunks * size]byte, 0, pSize),
		freeptrs:  make([]uintptr, 0, pSize),
	}
}

// noescape: unused now
func (p *UPool{{Title .Type}}) noescape(up unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(up) ^ 0)
}

// uintptr2t: uintptr to type convert
func (p *UPool{{Title .Type}}) uintptr2t(ptr uintptr) *{{.Type}} {
	return (*{{.Type}})(unsafe.Pointer(ptr))
}

// t2uintptr: type to uintptr convert
func (p *UPool{{Title .Type}}) t2uintptr(v *{{.Type}}) uintptr {
	return uintptr(unsafe.Pointer(v))
}

// malloc: allocate new memory chunk
func (p *UPool{{Title .Type}}) malloc() {
	p.memChunks = append(p.memChunks, [maxInChunks * size]byte{})
}

// Get struct from pool
func (p *UPool{{Title .Type}}) Get() *{{.Type}} {
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
	return (*{{.Type}})(unsafe.Pointer(&p.memChunks[p.currChunk][p.currOffset-size]))
}


//Return struct to pool
func (p *UPool{{Title .Type}}) Return(st *{{.Type}}) {
	p.freeptrs = append(p.freeptrs, p.t2uintptr(st))
}
`

type values struct {
	Package   string
	File      string
	Type      string
	MaxChunks int
}

func main() {
	settings := values{
		Type:      "",
		Package:   "",
		File:      "",
		MaxChunks: 10,
	}
	flag.StringVar(&settings.Package, "package", "", "Package name")
	flag.StringVar(&settings.File, "file", "", "Path to file")
	flag.StringVar(&settings.Type, "type", "", "Type for unsafe struct pool")
	flag.Parse()

	if settings.Package == "" {
		log.Fatal("package not set")
	}

	if settings.File == "" {
		log.Fatal("file not set")
	}
	if settings.Type == "" {
		log.Fatal("type not set")
	}

	file, err := os.Create(settings.File)
	if err == nil {
		defer file.Close()
		err = template.Must(
			template.New(
				"upool",
			).Funcs(
				map[string]interface{}{"Title": strings.Title},
				).Parse(
				tmpl,
			),
		).Execute(
			file,
			&settings,
		)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Unsafe pool generated for package %s in file %s", settings.Package, settings.File)
}
