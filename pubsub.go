package wireless

import (
	"strings"
)

// Subscription represents a subscription to one or events
type Subscription struct {
	ch     chan Event
	topics string
}

func (s *Subscription) publish(ev Event) error {
	if s.ch == nil {
		return nil
	}

	// don't let slow consumers cause goroutine leak, drop
	// events if the channel is full
	if len(s.ch) == cap(s.ch) {
		return nil
	}

	s.ch <- ev
	return nil
}

// Next will return a channel that returns events
func (s *Subscription) Next() chan Event {
	return s.ch
}

// Unsubscribe closes the channel and sets it to nil
func (s *Subscription) Unsubscribe() {
	close(s.ch)
	s.ch = nil
}

func (c *Conn) publishEvent(ev Event) {
	for _, sub := range c.subs {
		if strings.Contains(sub.topics, ev.Name) || sub.topics == "" {
			sub.publish(ev)
		}
	}
}

// Subscribe to one or more events and return the subscription
func (c *Conn) Subscribe(eventNames ...string) *Subscription {
	sub := &Subscription{
		topics: strings.Join(eventNames, " "),
		ch:     make(chan Event, 99),
	}

	c.subs = append(c.subs, sub)
	return sub
}
