package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

type FileReader interface {
	Start(paths []string)
}

type DefaultFileReader struct {
	ch   chan any
	ctx  context.Context
	path []string
	wg   *sync.WaitGroup
}

func NewDefaultFileReader(ctx context.Context) (*DefaultFileReader, <-chan any) {
	ch := make(chan any, 1000)
	return &DefaultFileReader{
		ch:  ch,
		ctx: ctx,
		wg:  &sync.WaitGroup{},
	}, ch
}

func (r *DefaultFileReader) Start(paths []string) {
	for _, path := range paths {
		r.read(path)
	}
	r.stateListen()
}

func (r *DefaultFileReader) read(path string) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("文件打开失败,%w\n", err)
			return
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		for {
			select {
			case <-r.ctx.Done():
				fmt.Println("文件处理取消")
				return
			default:
				line, err := reader.ReadString('\n')
				fmt.Printf("read line,%v\n", line)
				if err != nil {
					if err == io.EOF {
						r.ch <- line
						fmt.Printf("%s 已处理完成\n", path)
						return
					} else {
						fmt.Println("文件读取失败")
						return
					}
				}
				r.ch <- line
			}

		}
	}()
}

func (r *DefaultFileReader) stateListen() {
	go func() {
		r.wg.Wait()
		fmt.Println("文件读取已完成")
		close(r.ch)
	}()
}
