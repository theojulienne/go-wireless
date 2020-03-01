package wireless

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubscribe(t *testing.T) {
	Convey("given a connection", t, func() {
		c := &Conn{}

		Convey("when Subscribed to a Connected event", func() {
			sub := c.Subscribe(EventConnected)

			Convey("then it should have Connected in it's topics", func() {
				So(sub.topics, ShouldEqual, EventConnected)
			})

			Convey("and a Connected event is pushed", func() {
				c.publishEvent(Event{Name: EventConnected, Arguments: map[string]string{"a": "b"}})

				Convey("then the subscription should get the event", func() {
					So(sub.ch, ShouldHaveLength, 1)
					ev := <-sub.Next()
					So(ev.Name, ShouldEqual, EventConnected)
					So(ev.Arguments, ShouldContainKey, "a")
					So(ev.Arguments["a"], ShouldEqual, "b")
				})
			})

			Convey("and then unsubscribed", func() {
				sub.Unsubscribe()

				Convey("then publishing should not panic", func() {
					So(func() { c.publishEvent(Event{Name: EventConnected}) }, ShouldNotPanic)
					So(func() { sub.publish(Event{Name: EventConnected}) }, ShouldNotPanic)
				})
			})
		})

		Convey("when Subscribed to nothing", func() {
			sub := c.Subscribe()

			Convey("then it should have Connected in it's topics", func() {
				So(sub.topics, ShouldEqual, "")
			})

			Convey("and a Connectd event is pushed", func() {
				c.publishEvent(Event{Name: EventConnected})

				Convey("then the subscription should get the event", func() {
					So(sub.ch, ShouldHaveLength, 1)
					ev := <-sub.Next()
					So(ev.Name, ShouldEqual, EventConnected)
				})
			})
		})
	})
}
