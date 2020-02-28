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
	// if _, ok := <-s.ch; !ok {
	// 	return errors.New("Topic has been closed")
	// }
	s.ch <- ev
	return nil
}

// Next will return a cahnnel that returns events
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
		if sub.ch == nil {
			continue
		}

		// don't let slow consumers cause goroutine leak, drop
		// events if the channel is full
		if len(sub.ch) == cap(sub.ch) {
			continue
		}

		if strings.Contains(sub.topics, ev.Name) {
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
