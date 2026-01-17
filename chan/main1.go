package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	c := make(chan int, 5)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 100 {
			c <- i
			//		fmt.Printf("send: %d\n", i)
		}
		close(c)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range c {
			fmt.Printf("1.recv: %d\n", v)
		}
	}()
	wg.Wait()
	fmt.Println("game over")
}
