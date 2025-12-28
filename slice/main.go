package main

import "fmt"

func main() {

	a := []int{1, 2}
	b := a[:1]
	c := a
	c = append(c, 100)
	fmt.Println(a, b, c)
	fmt.Println(cap(a), cap(b), cap(c))
	add(a)

	fmt.Println(a, b, c)
}

func add(s []int) {
	s = append(s, 100)
}
