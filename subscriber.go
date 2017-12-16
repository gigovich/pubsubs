package pubsubs

// Subscriber container
type Subscriber struct {
	// SubsID identifier
	SubsID string

	// Notify chanel
	Notify chan interface{}
}
