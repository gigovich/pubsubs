package pubsubs

import (
	"fmt"
)

// ErrNoSuchSubscription error returned when subscription with such id not found
var ErrNoSuchSubscription = fmt.Errorf("no such subscription")

// ErrSubscriptionExists error returned when you try to pass subscrition to broker but it already exists in it
var ErrSubscriptionExists = fmt.Errorf("subscription already exists")

// Broker for publishers and subscribers
type Broker struct {
	subscriptions map[string]*Subscription
}

// New broker instance constructor
func New() *Broker {
	return &Broker{
		subscriptions: make(map[string]*Subscription),
	}
}

// Publish subscription
func (b *Broker) Publish(subscription *Subscription) error {
	if _, ok := b.subscriptions[subscription.id]; !ok {
		b.subscriptions[subscription.id] = subscription
		return nil
	}

	return ErrSubscriptionExists
}

// Unpublish subscription
func (b *Broker) Unpublish(subscription *Subscription) error {
	if _, ok := b.subscriptions[subscription.id]; ok {
		subscription.UnsubscribeAll()
		delete(b.subscriptions, subscription.id)
		return nil
	}
	return ErrNoSuchSubscription
}

// Subscribe for event by hash
func (b *Broker) Subscribe(id string) (*Subscriber, error) {
	if subscription, ok := b.subscriptions[id]; ok {
		return subscription.Subscribe(), nil
	}
	return nil, ErrNoSuchSubscription
}

// Unsubscribe subscriber from broker
func (b *Broker) Unsubscribe(subscriber *Subscriber) error {
	if subscription, ok := b.subscriptions[subscriber.SubsID]; ok {
		subscription.Unsubscribe(subscriber)
		return nil
	}
	return ErrNoSuchSubscription
}
