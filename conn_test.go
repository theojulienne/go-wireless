package wireless

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProcessMessage(t *testing.T) {
	c := &Conn{}

	Convey("given an event message", t, func() {
		data := []byte("<3>" + EventConnected + " ssid=watership")
		n := len(data)

		Convey("and a subscription to that event", func() {
			sub := c.Subscribe(EventConnected)

			Convey("when the message is processed", func() {
				c.processMessage(n, data, nil)

				Convey("then the event should be published", func() {
					So(sub.ch, ShouldHaveLength, 1)
					ev := <-sub.Next()
					So(ev.Name, ShouldEqual, EventConnected)
				})
			})
		})
	})

	Convey("given a log message", t, func() {
		data := []byte("<3>hey hey hey")
		n := len(data)

		Convey("and a subscription to that logs", func() {
			sub := c.Subscribe("logs")

			Convey("when the message is processed", func() {
				c.processMessage(n, data, nil)

				Convey("then the event should be published", func() {
					So(sub.ch, ShouldHaveLength, 1)
					ev := <-sub.Next()
					So(ev.Name, ShouldEqual, "logs")
					So(ev.Arguments["msg"], ShouldEqual, "hey hey hey")
				})
			})
		})
	})

	Convey("given a command response", t, func() {
		data := []byte("OK\n")
		n := len(data)
		c.currentCommandResponse = make(chan string, 1)

		Convey("when the message is processed", func() {
			c.processMessage(n, data, nil)

			Convey("then it should be put in the command response channel", func() {
				So(c.currentCommandResponse, ShouldHaveLength, 1)
				So(<-c.currentCommandResponse, ShouldEqual, "OK\n")
			})
		})
	})
}
