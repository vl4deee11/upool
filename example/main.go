package main

import "fmt"

//go:generate ./upoolcodegen.bin --package=main --file=upoolint.go --type=int
func main() {
	p := NewUPoolInt(
		func() int {
			return 6
		},
		10,
	)

	v := p.Get()
	fmt.Println(v)
	fmt.Println(*v)
	*v = 7
	p.Return(v)

	v2 := p.Get()
	fmt.Println(v2)
	fmt.Println(*v2)
}
