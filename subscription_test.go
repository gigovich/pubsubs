package pubsubs

import (
	"testing"
	"time"
)

func TestSubscription(t *testing.T) {
	s := NewSubscription("test")
	sb := s.Subscribe()
	go s.Publish("message")
	select {
	case <-time.After(time.Millisecond):
		t.Error("published message not received")
		return
	case msg := <-sb.Notify:
		if msg != "message" {
			t.Errorf("we got wrong message '%v', expected 'message'", msg)
			return
		}
	}
	s.Unsubscribe(sb)

	_, ok := <-sb.Notify
	if ok {
		t.Error("Notify channel of subscriber should be closed, after unsubscribe")
		return
	}

	if len(s.Subscribers) != 0 {
		t.Error("subscribers list should be empty")
	}

	finish := make(chan bool)
	go func() {
		s.Publish("message2")
		finish <- true
	}()

	select {
	case <-time.After(time.Millisecond):
		t.Errorf("publish blocked")
		return
	case <-finish:
		break
	}
}
