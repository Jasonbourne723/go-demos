package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	taskRunner := NewTaskRunner(WithTimeout(3*time.Second), WithLimit(3))
	taskRunner.Run(func() error {
		for i := range 5 {
			fmt.Println(i)
			time.Sleep(time.Millisecond * 1000)
		}
		return nil
	}, func() error {
		for i := 5; i < 10; i++ {
			fmt.Println(i)
			time.Sleep(time.Millisecond * 200)
		}
		return nil
	})

	if err := taskRunner.Wait(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("task is finished")
	}

}

type task func() error
type opt func(*TaskRunner)

type TaskRunner struct {
	limit    int
	ctx      context.Context
	timeout  time.Duration
	tasks    []task
	resultCh chan error
	group    *errgroup.Group
}

func WithLimit(limit int) opt {
	return func(t *TaskRunner) {
		t.limit = limit
	}
}

func WithContext(ctx context.Context) opt {
	return func(t *TaskRunner) {
		t.ctx = ctx
	}
}

func WithTimeout(duration time.Duration) opt {
	return func(t *TaskRunner) {
		t.timeout = duration
	}
}

func NewTaskRunner(opts ...opt) *TaskRunner {

	tr := TaskRunner{}
	for _, opt := range opts {
		opt(&tr)
	}
	if tr.ctx == nil {
		tr.ctx = context.Background()
	}
	if tr.timeout > 0 {
		tr.ctx, _ = context.WithTimeout(tr.ctx, tr.timeout)
	}
	tr.group, tr.ctx = errgroup.WithContext(tr.ctx)
	if tr.limit > 0 {
		tr.group.SetLimit(tr.limit)
	}
	return &tr
}

func (t *TaskRunner) Run(tasks ...task) error {
	t.tasks = tasks
	t.resultCh = make(chan error, len(t.tasks))
	t.work()

	return nil
}

func (t *TaskRunner) work() {
	for i, task := range t.tasks {
		t.group.Go(func() error {
			go func() {
				t.resultCh <- task()
			}()
			select {
			case <-t.ctx.Done():
				fmt.Println("task is cancel")
				return t.ctx.Err()
			case err := <-t.resultCh:
				fmt.Printf("the %dth task is completed", i)
				return err
			}
			return nil
		})
	}
}

func (t *TaskRunner) Wait() error {
	return t.group.Wait()
}
