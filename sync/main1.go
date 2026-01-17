package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mut := &sync.Mutex{}
	cond := sync.NewCond(mut)

	for i := range 3 {
		go func() {
			mut.Lock()
			defer mut.Unlock()
			cond.Wait()
			fmt.Println(i)
		}()
	}

	time.Sleep(time.Second)
	for _ = range 3 {
		cond.Signal()
	}
	<-time.After(time.Second * 3)
}
