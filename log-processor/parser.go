package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type parser interface {
	Start()
}

type DefaultParser struct {
	ctx      context.Context
	sourceCh <-chan any
	resultCh chan LogInfo
}

func NewDefaultParser(ctx context.Context, ch <-chan any) (*DefaultParser, <-chan LogInfo) {
	resultCh := make(chan LogInfo, 100)
	return &DefaultParser{
		ctx:      ctx,
		sourceCh: ch,
		resultCh: resultCh,
	}, resultCh
}

func (p *DefaultParser) Start() {
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			case data, ok := <-p.sourceCh:
				if !ok {
					close(p.resultCh)
					return
				}
				p.process(data)
			}
		}
	}()
}

func (p *DefaultParser) process(data any) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("parser error,%w\n", err)
		}
	}()
	fmt.Println("parser recived:", data)
	infoSli := strings.Split(data.(string), " ")
	if len(infoSli) != 4 {
		panic(errors.New("log format is wrong"))
	}
	createdat, err := time.Parse("2006-01-02T15:04:05", infoSli[0])
	if err != nil {
		panic(err)
	}
	userId, err := strconv.Atoi(infoSli[2])
	if err != nil {
		panic(err)
	}
	logInfo := LogInfo{
		createdat: createdat,
		userId:    int64(userId),
		level:     levelType(infoSli[1]),
		action:    infoSli[3],
	}
	fmt.Println("parse successed:", logInfo)
	p.resultCh <- logInfo
}
