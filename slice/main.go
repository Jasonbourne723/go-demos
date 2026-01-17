package main

import "fmt"

func main() {

	var i any = 3

	switch t := i.(type) {
	case int:
		fmt.Println("int")
	default:
		fmt.Println(t)
	}

	fmt.Println("------------")

	a := []int{1, 2}
	b := a[:1]
	c := a
	c = append(c, 100)
	fmt.Println(a, b, c)
	fmt.Println(cap(a), cap(b), cap(c))
	add(a)

	fmt.Println(a, b, c)

	fmt.Println("--------")

	fmt.Println(-5 % 4)
}

func add(s []int) {
	s = append(s, 100)
}
