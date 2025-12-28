package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	fmt.Println("vim-go")

	//	testforeach()
	//	testEmptyRun()
	//	testchan()
	//	testwaitgroup()
	// 	testerrgroup()
	// 	testmutex()
	//	testerrgroupwithcontext()
}

func testmutex() {

	mu := sync.Mutex{}
	m := make(map[int]int, 10)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 10 {
			func() {
				mu.Lock()
				defer mu.Unlock()
				m[i] = i
			}()
		}
	}()
	time.Sleep(time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 10 {
			func() {
				mu.Lock()
				defer mu.Unlock()
				fmt.Println(m[i])
			}()
		}
	}()
	wg.Wait()
}

func testerrgroupwithcontext() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	eg, ctx1 := errgroup.WithContext(ctx)
	eg.SetLimit(2)

	eg.Go(func() error {
		for {
			select {
			case <-ctx1.Done():
				fmt.Println("work is cancel")
				return ctx1.Err()
			case <-time.After(time.Millisecond * 300):
				fmt.Println("work is running")
			}
		}
		return nil
	})

	eg.Go(func() error {
		fmt.Println("work start")
		return nil
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("work is finish")
	}

}

func testerrgroup() {
	var eg errgroup.Group
	eg.Go(func() error {
		for i := range 3 {
			fmt.Println(i)
		}
		return nil
	})
	eg.Go(func() error {
		fmt.Println("4")
		return errors.New("failed")
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err.Error())
	}

}

func testwaitgroup() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 3 {
			fmt.Println(i)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 3; i < 6; i++ {
			fmt.Println(i)
		}
	}()
	wg.Wait()
}

func testchan() {
	// 无缓冲chan，写入时，如果没有消费者，会被阻塞
	// 有缓存chan，当chan满时，插入也会被阻塞
	ch := make(chan int, 1)
	ch <- 1

	quit := make(chan int, 2)
	quit <- 10
	quit <- 9
	close(quit)

	v1, ok1 := <-quit
	v2, ok2 := <-quit
	// 当chan关闭时ok为false
	fmt.Println(v1, ok1)
	fmt.Println(v2, ok2)

	// chan关闭时，循环会自动结束
	for i := range quit {
		fmt.Println(i)
	}
}

func testEmptyRun() {
	// 1.14版本后gmp模型使用抢占式调度模型，cpu空跑不会完全占用线程资源
	go func() {
		for {
		}
	}()

	//time.Sleep(time.Second)
	fmt.Println("end")
}

func testforeach() {
	// for循环内执行goroutine闭包问题已被修复
	for i := range 10 {
		go func() {
			fmt.Println(i)
		}()
	}
	//time.Sleep(time.Second)
}
