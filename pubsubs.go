package pubsubs

import (
	"fmt"
)

// ErrNoSuchTopic error returned when topic with such id not found
var ErrNoSuchTopic = fmt.Errorf("no such topic")

// ErrTopicExists error returned when you try to add topic to broker but it already exists in it
var ErrTopicExists = fmt.Errorf("topic already exists")

// Broker for publishers and subscribers
type Broker struct {
	topics map[string]*Topic
}

// New broker instance constructor
func New() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

// Add topic
func (b *Broker) Add(topic *Topic) error {
	if _, ok := b.topics[topic.id]; !ok {
		b.topics[topic.id] = topic
		return nil
	}

	return ErrTopicExists
}

// Remove topic
func (b *Broker) Remove(topic *Topic) error {
	if _, ok := b.topics[topic.id]; ok {
		topic.UnsubscribeAll()
		delete(b.topics, topic.id)
		return nil
	}
	return ErrNoSuchTopic
}

// Subscribe for event by hash
func (b *Broker) Subscribe(id string) (*Subscription, error) {
	if topic, ok := b.topics[id]; ok {
		return topic.Subscribe(), nil
	}
	return nil, ErrNoSuchTopic
}

// Unsubscribe subscriber from broker
func (b *Broker) Unsubscribe(subscription *Subscription) error {
	if topic, ok := b.topics[subscription.topicID]; ok {
		topic.Unsubscribe(subscription)
		return nil
	}
	return ErrNoSuchTopic
}
