package wireless

import (
	"strings"
)

// Subscription represents a subscription to one or events
type Subscription struct {
	ch     chan<- Event
	topics string
}

func (s *Subscription) publish(ev Event) error {
	// if _, ok := <-s.ch; !ok {
	// 	return errors.New("Topic has been closed")
	// }
	s.ch <- ev
	return nil
}

// Next will return a cahnnel that returns the next event
func (s *Subscription) Next(ev Event) chan<- Event {
	return s.ch
}

func (c *Conn) publishEvent(ev Event) {
	for _, sub := range c.subs {
		if strings.Contains(sub.topics, ev.Name) {
			sub.publish(ev)
		}
	}
}

// Subscribe to one or more events and return the subscription
func (c *Conn) Subscribe(eventNames ...string) Subscription {
	sub := Subscription{
		topics: strings.Join(eventNames, " "),
		ch:     make(chan Event, 99),
	}

	c.subs = append(c.subs, sub)
	return sub
}
