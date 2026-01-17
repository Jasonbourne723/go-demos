package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	fmt.Println("vim-go")
	// done
	once := sync.Once{}
	for _ = range 10 {
		once.Do(func() {
			fmt.Println("100")
		})

	}
	// waitgroup
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 1)
		fmt.Println(1)
	}()

	wg.Wait()

	var eg errgroup.Group
	eg.SetLimit(3)
	eg.Go(func() error {
		time.Sleep(time.Second)
		fmt.Println("11")
		return nil
	})
	eg.Go(func() error {
		time.Sleep(time.Second)
		fmt.Println("11")
		return errors.New("test error")
	})
	if err := eg.Wait(); err != nil {
		fmt.Printf("%v\n", err)
	}

	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()
	fmt.Println("lock")

	type person struct {
		name string
	}

	pool := sync.Pool{
		New: func() any {
			return person{}
		},
	}

	for i := 0; i < 3; i++ {
		val := pool.Get().(person)
		fmt.Println(val.name)
		val.name = strconv.Itoa(i)
		pool.Put(val)
	}

	var config atomic.Value
	config.Store(person{
		name: "lilie",
	})
	p := config.Load().(person)
	fmt.Println(p.name)

	c1, c2 := make(chan int, 1), make(chan int, 1)
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for _ = range 3 {
			c2 <- 0
			fmt.Println(0)
			<-c1
		}
	}()
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for _ = range 3 {
			<-c2
			fmt.Println(1)
			c1 <- 0
		}
	}()
	wg1.Wait()
}
