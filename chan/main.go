package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	fmt.Println("vim-go")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	eventProcessor := NewDefaultEventProcessor(WithContext(ctx), WithWorkerLimit(10))
	eventProcessor.Start()

	eg, ctx := errgroup.WithContext(ctx)
	for i := range 10 {
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("producer is cancel")
					return ctx.Err()
				case <-time.After(time.Millisecond * 200):
					eventProcessor.Submit(Event{
						Data:     "test",
						Priority: i,
					})
				}
			}
			return nil
		})
	}

	<-ctx.Done()
	eventProcessor.Stop(time.Second * 2)

}

const (
	running int = iota
	draining
	stopped
)

type Event struct {
	Priority int
	Data     any
}

type EventProcessor interface {
	Start() error
	Stop() error
	Submit(Event) error
}

type Worker interface {
	Consume(<-chan Event) error
}

type DefaultWorker struct {
	id  int
	ctx context.Context
	ch  chan Event
}

func (w *DefaultWorker) Consume(ch <-chan Event) error {
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				fmt.Println("worker was stopped")
				return
			case event := <-w.ch:
				time.Sleep(time.Millisecond * 500)
				fmt.Printf("event was alreadly processed,%v\n", event)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-w.ctx.Done():
				fmt.Println("worker was stopped")
				return
			case event := <-ch:
				w.ch <- event
			}
		}
	}()
	return nil
}

type opt func(*DefaultEventProcessor)

func WithContext(ctx context.Context) opt {
	return func(e *DefaultEventProcessor) {
		e.ctx, e.cancelFunc = context.WithCancel(ctx)
	}
}
func WithWorkerLimit(limit int) opt {
	return func(e *DefaultEventProcessor) {
		e.limit = limit
	}
}

func NewDefaultEventProcessor(opts ...opt) *DefaultEventProcessor {
	ep := &DefaultEventProcessor{
		limit:     10,
		startOnce: sync.Once{},
		closeOnce: sync.Once{},
		ch:        make(chan Event, 5),
		status:    atomic.Value{},
		demoteCh:  make(chan Event, 100),
	}
	for _, opt := range opts {
		opt(ep)
	}
	return ep
}

type DefaultEventProcessor struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	ch         chan Event
	limit      int
	workers    []Worker
	startOnce  sync.Once
	closeOnce  sync.Once
	status     atomic.Value
	demoteCh   chan Event
}

func (e *DefaultEventProcessor) Start() {
	e.startOnce.Do(func() {
		e.status.Store(running)
		e.demoteProcess()
		for i := range e.limit {
			worker := &DefaultWorker{
				id:  i,
				ctx: e.ctx,
				ch:  make(chan Event, 5),
			}
			worker.Consume(e.ch)
			e.workers = append(e.workers, worker)
		}
	})
}

func (e *DefaultEventProcessor) Stop(timeout time.Duration) {
	if ok := e.status.CompareAndSwap(running, draining); !ok {
		fmt.Println("repeat stop")
		return
	}
	fmt.Println("srv is stopping")
	timeoutCtx, _ := context.WithTimeout(context.Background(), timeout)
	<-timeoutCtx.Done()
	e.cancelFunc()
	if ok := e.status.CompareAndSwap(draining, stopped); ok {
		fmt.Println("srv was stopped")
	}
}

func (e *DefaultEventProcessor) demoteProcess() {
	go func() {
		for {
			select {
			case event := <-e.demoteCh:
				func() {
					select {
					case e.ch <- event:
						fmt.Printf("reenqueue successed,%v\n", event)
					case <-time.After(time.Millisecond * 200):
						fmt.Println("reenqueue failed,throw ")
					}
				}()
			}
		}
	}()
}

func (e *DefaultEventProcessor) Submit(event Event) error {
	if status := e.status.Load().(int); status != running {
		return errors.New("srv is stopping")
	}
	select {
	case e.ch <- event:
		fmt.Printf("enqueue successed,%v\n", event)
	case <-time.After(time.Millisecond * 50):
		if event.Priority >= 5 {
			e.demoteCh <- event
		}
	}
	return nil
}
