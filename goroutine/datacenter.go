package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	dc := NewDefaultDataCenter(ctx)
	group, ctx := errgroup.WithContext(ctx)
	// 并发测试getall方法
	for _ = range 10 {
		group.Go(func() error {
			for _ = range 10 {
				time.Sleep(time.Second)
				select {
				case <-ctx.Done():
					fmt.Println("get test is cancel")
					return ctx.Err()
				default:
					data := dc.GetAll()
					fmt.Printf("GetAll: %v\n", data)
				}
			}
			return nil
		})
	}
	// 并发Loadconfig
	for _ = range 10 {
		group.Go(func() error {
			for _ = range 5 {
				select {
				case <-ctx.Done():
					fmt.Println("get test is cancel")
					return ctx.Err()
				case <-time.After(time.Second * 2):
					dc.LoadFromFile("./config.json")
				}
			}
			return nil
		})
	}
	// 测试订阅

	for _ = range 10 {
		group.Go(func() error {

			subscribe := dc.Subscribe()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("subscribe is cancel")
					return ctx.Err()
				case data := <-subscribe:
					fmt.Printf("订阅更新:%v\n", data)
				}
			}

			return nil
		})
	}
	// 模拟修改配置文件

	group.Go(func() error {
		for i := range 10 {

			f, err := os.OpenFile("./config.json", os.O_RDWR, os.ModePerm)
			if err != nil {
				return err
			}
			decoder := json.NewDecoder(f)
			var m map[string]string
			if err := decoder.Decode(&m); err != nil {
				return err
			}
			m["key"+strconv.Itoa(i)] = "val" + strconv.Itoa(i)
			if _, err := f.Seek(0, 0); err != nil {
				return err
			}
			encoder := json.NewEncoder(f)
			if err := encoder.Encode(m); err != nil {
				return err
			}
			<-time.After(time.Second)
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Printf("err: %w\n", err)
	}

	fmt.Println("task was finished, system will close after 2s")

	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
	<-timeoutCtx.Done()
	fmt.Println("system is closing")

}

var d DataCenter = &defaulDataCenter{}

type DataCenter interface {
	Get(string) (string, bool)
	GetAll() map[string]string
	Subscribe() <-chan map[string]string
	UnSubscribe(ch chan map[string]string)
	LoadFromFile(path string) error
}

func NewDefaultDataCenter(ctx context.Context) *defaulDataCenter {
	atomicval := atomic.Value{}
	atomicval.Store(map[string]string{})
	dc := &defaulDataCenter{
		ctx:        ctx,
		data:       atomicval,
		subscriber: make(map[chan map[string]string]struct{}, 10),
		mutex:      sync.RWMutex{},
		notify:     make(chan struct{}, 10),
	}
	go dc.broadcast()
	return dc
}

type defaulDataCenter struct {
	ctx        context.Context
	data       atomic.Value // map[string]string
	subscriber map[chan map[string]string]struct{}
	mutex      sync.RWMutex
	notify     chan struct{}
}

func (d *defaulDataCenter) LoadFromFile(path string) error {
	if len(path) == 0 {
		return errors.New("invalid path")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	var data map[string]string
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	d.data.Store(data)
	select {
	case d.notify <- struct{}{}:
	default:
	}
	return nil
}

func (d *defaulDataCenter) broadcast() {

	for {

		select {
		case <-d.ctx.Done():
			fmt.Println("broadcast is cancel")
			return
		case <-d.notify:
			d.mutex.RLock()
			subscribers := make(map[chan map[string]string]struct{}, len(d.subscriber))
			for subscriber, _ := range d.subscriber {
				subscribers[subscriber] = struct{}{}
			}
			d.mutex.RUnlock()
			data := d.data.Load().(map[string]string)
			for subscriber := range subscribers {
				m := make(map[string]string, len(data))
				for k, v := range data {
					m[k] = v
				}
				select {
				case subscriber <- m:
					fmt.Println("broadcast")
				case <-time.After(time.Millisecond * 200):
					d.UnSubscribe(subscriber)
				}
			}
		}
	}

}

func (d *defaulDataCenter) Get(key string) (string, bool) {
	data := d.data.Load().(map[string]string)
	v, exist := data[key]
	return v, exist
}

func (d *defaulDataCenter) GetAll() map[string]string {
	data := d.data.Load().(map[string]string)
	m := make(map[string]string, len(data))
	for k, v := range data {
		m[k] = v
	}
	return m
}

func (d *defaulDataCenter) Subscribe() <-chan map[string]string {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	ch := make(chan map[string]string, 10)
	d.subscriber[ch] = struct{}{}
	return ch
}

func (d *defaulDataCenter) UnSubscribe(ch chan map[string]string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.subscriber, ch)
	close(ch)
}
