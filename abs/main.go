package main

import "fmt"

func main() {

	man := &Man{
		Person{
			Name: "lilei",
		},
	}
	fmt.Println(man.SayName())

}

type Person struct {
	Name string
}

func (p *Person) SayName() string {
	return p.Name
}

type Man struct {
	Person
}

func (m *Man) SayName() string {
	return m.Name
}

type Woman struct {
	Person
}
