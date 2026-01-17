package main

import (
	"flag"
	"fmt"
	"sync"
)

var name string

func main() {

	flag.StringVar(&name, "name", "liei", "")
	flag.Parse()

	fmt.Println(name)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		wg.Add(1)
		defer wg.Done()

		fmt.Println(1)
	}()

	wg.Wait()
	fmt.Println("hello")
}
