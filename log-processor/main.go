package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	paths := []string{"./log/1.log", "./log/2.log", "./log/3.log"}
	processor := NewDefaultProcessor(WithFilePath(paths...), WithContext(context.Background()), WithFilter(func(log *LogInfo) bool {
		return log.userId != 111
	}))
	var wg sync.WaitGroup
	wg.Add(1)
	callback := func(result *Result) {
		defer wg.Done()
		if result == nil {
			fmt.Println("empty result")
			return
		}
		fmt.Println("createdat\tcount")
		for k, v := range result.AggregationByMinutes {
			fmt.Printf("%v\t%d\n", k, v)
		}
		fmt.Println("userId\tcount")
		for k, v := range result.AggregationByUsers {
			fmt.Printf("%d\t%d\n", k, v)
		}
	}
	err := processor.Start(callback)
	if err != nil {
		fmt.Println(err)
	}
	wg.Wait()
	fmt.Println("game over")
}
