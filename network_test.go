package wireless

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var pskNet = Network{ID: 0, IDStr: "h", SSID: "gg", PSK: "xx"}
var pskNet2 = Network{ID: 0, IDStr: "h", SSID: "gg", PSK: "xx", KeyMgmt: "WPA2-TKIP", ScanSSID: true}
var knwonPSKNet = Network{Known: true, ID: 0, IDStr: "h", SSID: "gg", KeyMgmt: "WPA2-TKIP", ScanSSID: true}
var openNet = Network{ID: 0, IDStr: "x", SSID: "gg"}

// func TestNetwork(t *testing.T) {
// 	Convey("given a blank network", t, func() {
// 	})
// }

func TestSetCmds(t *testing.T) {
	Convey("given a private network", t, func() {
		n := pskNet

		Convey("when the set commands are rendered", func() {
			cmds := setCmds(n)

			Convey("then it should have the basic fields", func() {
				So(cmds, ShouldContain, `SET_NETWORK 0 psk "xx"`)
				So(cmds, ShouldContain, `SET_NETWORK 0 ssid "gg"`)
			})
		})
	})

	Convey("given a disabled custom private network", t, func() {
		n := pskNet2
		n.Disable(true)

		Convey("when the set commands are rendered", func() {
			cmds := setCmds(n)
			Println(cmds)

			Convey("then it should have the basic fields", func() {
				So(cmds, ShouldContain, `SET_NETWORK 0 psk "xx"`)
				So(cmds, ShouldContain, `SET_NETWORK 0 ssid "gg"`)
				So(cmds, ShouldContain, `SET_NETWORK 0 key_mgmt WPA2-TKIP`)
				So(cmds, ShouldContain, `SET_NETWORK 0 disabled 1`)
				So(cmds, ShouldContain, `SET_NETWORK 0 scan_ssid 1`)
			})
		})
	})

	Convey("given a known private network", t, func() {
		n := knwonPSKNet

		Convey("when the set commands are rendered", func() {
			cmds := setCmds(n)
			Println(cmds)

			Convey("then it should not contain the PSK", func() {
				So(len(cmds), ShouldEqual, 3)
				So(cmds, ShouldContain, `SET_NETWORK 0 ssid "gg"`)
				So(cmds, ShouldContain, `SET_NETWORK 0 key_mgmt WPA2-TKIP`)
				So(cmds, ShouldContain, `SET_NETWORK 0 scan_ssid 1`)
			})

		})

		Convey("when the password is reset", func() {
			n.PSK = "horsewaffle"
			Convey("when the set commands are rendered", func() {
				cmds := setCmds(n)
				Println(cmds)

				Convey("then it should not contain the PSK", func() {
					So(len(cmds), ShouldEqual, 4)
					So(cmds, ShouldContain, `SET_NETWORK 0 psk "horsewaffle"`)
					So(cmds, ShouldContain, `SET_NETWORK 0 ssid "gg"`)
					So(cmds, ShouldContain, `SET_NETWORK 0 key_mgmt WPA2-TKIP`)
					So(cmds, ShouldContain, `SET_NETWORK 0 scan_ssid 1`)
				})
			})
		})
	})

	Convey("given an open network", t, func() {
		n := openNet

		Convey("when the set commands are rendered", func() {
			cmds := setCmds(n)

			Convey("then it should have the proper fields", func() {
				So(cmds, ShouldContain, `SET_NETWORK 0 key_mgmt NONE`)
				So(cmds, ShouldContain, `SET_NETWORK 0 ssid "gg"`)
			})
		})
	})
}
