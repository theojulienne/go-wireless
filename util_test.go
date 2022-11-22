package wireless

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseAP(t *testing.T) {
	Convey("given some scan results", t, func() {
		Convey("when they are parsed", func() {
			aps, err := parseAP(scanResults)
			So(err, ShouldBeNil)

			Convey("it should have parsed 3 APs", func() {
				So(aps, ShouldHaveLength, 3)
			})

			Convey("then it should have AP0 first", func() {
				a := aps[0]
				So(a.SSID, ShouldEqual, "AP0")
				So(a.Frequency, ShouldEqual, 2472)
				So(a.Signal, ShouldEqual, -30)
				So(a.Flags, ShouldContain, "WPS")
				So(a.Flags, ShouldContain, "ESS")
				So(a.Flags, ShouldContain, "WPA2-PSK-CCMP")
				So(a.BSSID, ShouldEqual, "d0:7a:b5:31:23:a0")
			})

			Convey("then it should have AP2 last", func() {
				a := aps[2]
				So(a.SSID, ShouldEqual, "AP2")
				So(a.Frequency, ShouldEqual, 2412)
				So(a.Signal, ShouldEqual, -77)
				So(a.Flags, ShouldContain, "WPS")
				So(a.Flags, ShouldContain, "ESS")
				So(a.Flags, ShouldContain, "WPA2-PSK-CCMP+TKIP")
				So(a.Flags, ShouldContain, "WPA-PSK-CCMP+TKIP")
				So(a.BSSID, ShouldEqual, "24:00:ba:f8:65:df")
			})
		})
	})
}

func TestParseNetwork(t *testing.T) {
	Convey("given some networks", t, func() {
		Convey("when they are parsed", func() {
			nets, err := parseNetwork(networkList)
			So(err, ShouldBeNil)

			Convey("then it should have parsed 3 networks", func() {
				So(nets, ShouldHaveLength, 3)
			})

			Convey("then it should have AP0 first", func() {
				n := nets[0]
				So(n.ID, ShouldEqual, 0)
				So(n.SSID, ShouldEqual, "AP0")
				So(n.IsDisabled(), ShouldBeFalse)
				So(n.BSSID, ShouldEqual, "any")
				So(n.Flags, ShouldHaveLength, 0)
			})

			Convey("then it should have AP2 last", func() {
				n := nets[2]
				So(n.ID, ShouldEqual, 2)
				So(n.SSID, ShouldEqual, "AP2")
				So(n.IsDisabled(), ShouldBeTrue)
				So(n.BSSID, ShouldEqual, "any")
				So(n.Flags, ShouldHaveLength, 1)
				So(n.Flags, ShouldContain, "DISABLED")
			})
		})
	})
}

var scanResults = []byte(`bssid / frequency / signal level / flags / ssid
d0:7a:b5:31:23:a0	2472	-30	[WPA2-PSK-CCMP][WPS][ESS]	AP0
00:1f:1f:37:42:d9	2442	-37	[WPA2-PSK-CCMP][ESS]	AP1
24:00:ba:f8:65:df	2412	-77	[WPA-PSK-CCMP+TKIP][WPA2-PSK-CCMP+TKIP][WPS][ESS]	AP2
`)

var networkList = []byte(`network id / ssid / bssid / flags
0	AP0	any	
1	AP1	any	[DISABLED]
2	AP2	any	[DISABLED]`)
