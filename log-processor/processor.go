package main

import (
	"context"
	"sync"
	"time"
)

type levelType string

const (
	INOF  levelType = "INFO"
	WARN  levelType = "WARN"
	ERROR levelType = "ERROR"
)

type LogInfo struct {
	createdat time.Time
	userId    int64
	level     levelType
	action    string
}

type Result struct {
	AggregationByMinutes map[string]int
	AggregationByUsers   map[int64]int
}

type Processor interface {
	Start() (func(*Result), error)
}

type filter func(*LogInfo) bool
type opt func(*DefaultProcessor)

func WithFilePath(paths ...string) opt {
	return func(p *DefaultProcessor) {
		p.filePaths = paths
	}
}

func WithFilter(filters ...filter) opt {
	return func(p *DefaultProcessor) {
		p.filters = filters
	}
}

func WithContext(ctx context.Context) opt {
	return func(p *DefaultProcessor) {
		p.ctx = ctx
	}
}

type DefaultProcessor struct {
	ctx       context.Context
	filePaths []string
	filters   []filter
	callBack  func(*Result)
	wg        *sync.WaitGroup
	result    *Result
}

func NewDefaultProcessor(opts ...opt) *DefaultProcessor {
	p := DefaultProcessor{}
	p.result = &Result{
		AggregationByMinutes: make(map[string]int, 100),
		AggregationByUsers:   make(map[int64]int, 100),
	}
	p.wg = &sync.WaitGroup{}
	for _, opt := range opts {
		opt(&p)
	}
	return &p
}

func (p *DefaultProcessor) Start(callback func(*Result)) error {

	p.callBack = callback
	fileReader, dataCh := NewDefaultFileReader(p.ctx)
	parser, paserCh := NewDefaultParser(p.ctx, dataCh)
	p.aggregation(paserCh)
	parser.Start()
	fileReader.Start(p.filePaths)
	return nil
}

func (p *DefaultProcessor) aggregation(ch <-chan LogInfo) {
	go func() {
		defer func() {
			p.callBack(p.result)
		}()
		for {
			select {
			case <-p.ctx.Done():
				return
			case data, ok := <-ch:
				if !ok {
					return
				}
				p.process(&data)
			}
		}
	}()
}

func (p *DefaultProcessor) process(data *LogInfo) {
	for _, filter := range p.filters {
		if ok := filter(data); !ok {
			return
		}
	}
	t := data.createdat.Truncate(time.Minute)
	// 重新格式化为忽略秒的时间
	newTime := t.Format("2006-01-02T15:04:00")
	p.result.AggregationByUsers[data.userId]++
	p.result.AggregationByMinutes[newTime]++
}
