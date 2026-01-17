package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
	bus := NewBus()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		subscribeChan, doneChan := bus.Subscribe("topic-1", 100)
		for {
			select {
			case msg, ok := <-subscribeChan:
				if ok {
					fmt.Printf("consumer msg:%v\n", msg)
				}
			case <-doneChan:
				return
			default:

			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		subscribeChan, doneChan := bus.Subscribe("topic-1", 100)
		for {
			select {
			case msg, ok := <-subscribeChan:
				if ok {
					fmt.Printf("consumer msg:%v\n", msg)
				}
			case <-doneChan:
				return
			default:

			}
		}
	}()

	time.Sleep(time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 100 {
			bus.Publish("topic-1", i)
		}
	}()
	wg.Wait()
	fmt.Println("over")
}

func NewBus() *Bus {
	return &Bus{
		topics: make(map[string]*Topic, 10),
		mu:     sync.RWMutex{},
	}
}

type Bus struct {
	topics map[string]*Topic
	mu     sync.RWMutex
}

type Topic struct {
	subscribers map[int]*Subscriber
	mu          sync.RWMutex
	nextId      int
}

type Subscriber struct {
	id   int
	ch   chan any
	done chan struct{}
}

func (b *Bus) getTopic(topicName string) *Topic {
	b.mu.RLock()
	if t, exist := b.topics[topicName]; exist {
		return t
	} else {
		b.mu.RLocker().Unlock()
		b.mu.Lock()
		defer b.mu.Unlock()
		if t, exist := b.topics[topicName]; exist {
			return t
		} else {
			b.topics[topicName] = &Topic{
				subscribers: make(map[int]*Subscriber, 3),
				mu:          sync.RWMutex{},
				nextId:      1,
			}
			return b.topics[topicName]
		}
	}
}

func (b *Bus) Publish(topicName string, message any) error {
	topic := b.getTopic(topicName)
	return topic.publish(message)
}

func (b *Bus) Subscribe(topicName string, bufferSize int) (<-chan any, <-chan struct{}) {
	topic := b.getTopic(topicName)
	return topic.subscribe(bufferSize)
}

func (t *Topic) subscribe(bufferSize int) (<-chan any, <-chan struct{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	subscriber := &Subscriber{
		id:   t.nextId,
		ch:   make(chan any, bufferSize),
		done: make(chan struct{}, 1),
	}
	t.nextId++
	t.subscribers[subscriber.id] = subscriber
	return subscriber.ch, subscriber.done
}

func (t *Topic) publish(message any) error {
	var err error
	t.mu.RLock()
	defer t.mu.RLocker().Unlock()
	for _, subscriber := range t.subscribers {
		select {
		case subscriber.ch <- message:
			fmt.Printf("message %v is published to subscriber %d\n", message, subscriber.id)
		default:
			err = errors.New("subscriber consumer delay")
		}
	}
	if err != nil {
		return err
	}
	return nil
}
