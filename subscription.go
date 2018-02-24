package pubsubs

// Subscription container
type Subscription struct {
	// topicID identifier
	topicID string

	// Notify chanel
	Notify chan interface{}
}
