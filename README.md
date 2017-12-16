# Golang PUB/SUBS package

Golang publish/subscribe library

## Intall
Use [dep](https://github.com/golang/dep) tool:
```bash
$ dep ensure -add github.com/gigovich/dep
```

## Example
This program posts 10 times tick, every 1 second, for 5 subscribers:

```go
package main

import (
	"fmt"
	"github.com/gigovich/pubsubs"
	"sync"
	"time"
)

const tickSubscriptionName = "tickSubscription"

// broker can used as global registry of subscriptions
var broker = pubsubs.New()

// subscribe function receive values from publisher by subscription name
func subscribe(wg *sync.WaitGroup, num int) {
	defer wg.Done()

	subsc, err := broker.Subscribe(tickSubscriptionName)
	if err != nil {
		panic(err.Error())
	}

	// iterate over published values
	for tick := range subsc.Notify {
		now, ok := tick.(time.Time)
		if !ok {
			panic("we expect time as tick")
		}
		fmt.Printf("Goroutine #%v -> receive tick: %v\n", num+1, now.Format("15:04:05.999"))
	}
}

func main() {
	// create subscription by
	subsc := pubsubs.NewSubscription(tickSubscriptionName)
	if err := broker.Publish(subsc); err != nil {
		panic(err.Error())
	}

	// in example application in first order we run subscribers, because they lost all messages which
	// will be published until subscribe done
	wg := &sync.WaitGroup{}
	wg.Add(5)
	// run 5 subscribe functions
	for i := 0; i < 5; i++ {
		go subscribe(wg, i)
	}

	// run publish function
	go func() {
		// publish tick message every second
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			fmt.Println("--- New tick ---")

			// publish tick value
			subsc.Publish(time.Now())
		}

		if err := broker.Unpublish(subsc); err != nil {
			panic(err.Error())
		}
	}()

	wg.Wait()
}
```
