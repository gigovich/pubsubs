package pubsubs

import (
	"log"
	"sync"
	"time"
)

// Topic message
type Topic struct {
	id           string
	subscription map[*Subscription]*Subscription
	lastMessage  time.Time

	mutex sync.Mutex
}

// NewTopic instance constructor
func NewTopic(id string) *Topic {
	return &Topic{
		id:           id,
		subscription: make(map[*Subscription]*Subscription),
	}
}

// Publish message
func (s *Topic) Publish(message interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.lastMessage = time.Now()
	for _, subscriber := range s.subscription {
		select {
		case subscriber.Notify <- message:
		default:
			log.Println("subscriber channel filled,", subscriber, "drop message")
		}
	}
}

// Subscribe subscriber
func (s *Topic) Subscribe() *Subscription {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	subscription := &Subscription{
		topicID: s.id,
		Notify:  make(chan interface{}, 10),
	}

	s.subscription[subscription] = subscription
	return subscription
}

// Unsubscribe subscriber
func (s *Topic) Unsubscribe(subscription *Subscription) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.subscription[subscription]; ok {
		close(subscription.Notify)
		delete(s.subscription, subscription)
	}
}

// UnsubscribeAll subscriptions
func (s *Topic) UnsubscribeAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, subscription := range s.subscription {
		close(subscription.Notify)
		delete(s.subscription, id)
	}
}
