package main

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	source := make(chan rxgo.Item)
	observable := rxgo.FromEventSource(source, rxgo.WithBackPressureStrategy(rxgo.Block))

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	var received []int

	wg := &sync.WaitGroup{}

	source <- rxgo.Of(1)
	source <- rxgo.Of(2)

	wg.Add(1)
	go func() {
		<-observable.ForEach(func(i interface{}) {
			println(i.(int))
			mu.Lock()
			defer mu.Unlock()
			received = append(received, i.(int))
		}, func(err error) {
			panic(err)
		}, func() {
			println("finish foreach")
		}, rxgo.WithContext(ctx))

		println("finish receive")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		<-time.After(500 * time.Millisecond)

		source <- rxgo.Of(3)
		source <- rxgo.Of(4)

		cancel()

		source <- rxgo.Of(5)
		source <- rxgo.Of(6)

		wg.Done()
	}()

	wg.Wait()

	if len(received) != 2 {
		panic(received)
	}

	if received[0] != 3 || received[1] != 4 {
		panic(received)
	}
}
