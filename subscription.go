package pubsubs

import (
	"log"
	"sync"
	"time"
)

// Subscription message
type Subscription struct {
	id          string
	Subscribers map[*Subscriber]*Subscriber
	LastMessage time.Time

	mutex sync.Mutex
}

// NewSubscription instance constructor
func NewSubscription(id string) *Subscription {
	return &Subscription{
		id:          id,
		Subscribers: make(map[*Subscriber]*Subscriber),
	}
}

// Publish message
func (s *Subscription) Publish(message interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.LastMessage = time.Now()
	for _, subscriber := range s.Subscribers {
		select {
		case subscriber.Notify <- message:
		default:
			log.Println("subscriber channel filled,", subscriber, "drop message")
		}
	}
}

// Subscribe subscriber
func (s *Subscription) Subscribe() *Subscriber {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	subscriber := &Subscriber{
		Notify: make(chan interface{}, 10),
	}

	s.Subscribers[subscriber] = subscriber
	return subscriber
}

// Unsubscribe subscriber
func (s *Subscription) Unsubscribe(subscriber *Subscriber) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.Subscribers[subscriber]; ok {
		close(subscriber.Notify)
		delete(s.Subscribers, subscriber)
	}
}

// UnsubscribeAll subscriptions
func (s *Subscription) UnsubscribeAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, subscriber := range s.Subscribers {
		close(subscriber.Notify)
		delete(s.Subscribers, id)
	}
}
