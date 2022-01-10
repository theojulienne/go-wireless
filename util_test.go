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
				So(a.BSSID.String(), ShouldEqual, "d0:7a:b5:31:23:a0")
			})

			Convey("then it should have AP1 second", func() {
				a := aps[1]
				So(a.SSID, ShouldEqual, "Àndroid 😫♥️👍")
				So(a.Frequency, ShouldEqual, 2442)
				So(a.Signal, ShouldEqual, -37)
				So(a.Flags, ShouldContain, "WPA2-PSK-CCMP")
				So(a.Flags, ShouldContain, "ESS")
				So(a.BSSID.String(), ShouldEqual, "00:1f:1f:37:42:d9")
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
				So(a.BSSID.String(), ShouldEqual, "24:00:ba:f8:65:df")
			})
		})

		Convey("when quoted APs are parsed", func() {
			aps, err := parseAP(scanResultsQuote)
			So(err, ShouldBeNil)

			Convey("then it should have parsed 3 APs", func() {
				So(aps, ShouldHaveLength, 3)
			})

			Convey("the names should be correct", func() {
				So(aps[0].SSID, ShouldEqual, `AP0 "cool"`)
				So(aps[1].SSID, ShouldEqual, `"AP1"`)
				So(aps[2].SSID, ShouldEqual, `"AP2`)
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

			Convey("then it should have AP1 second", func() {
				n := nets[1]
				So(n.ID, ShouldEqual, 1)
				So(n.SSID, ShouldEqual, "Àndroid 😫♥️👍")
				So(n.IsDisabled(), ShouldBeTrue)
				So(n.BSSID, ShouldEqual, "any")
				So(n.Flags, ShouldHaveLength, 1)
				So(n.Flags, ShouldContain, "DISABLED")
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

		Convey("when quoted networks are parsed", func() {
			nets, err := parseNetwork(networkListQuote)
			So(err, ShouldBeNil)

			Convey("then it should have parsed 3 networks", func() {
				So(nets, ShouldHaveLength, 3)
			})

			Convey("the names should be correct", func() {
				So(nets[0].SSID, ShouldEqual, `AP0 "cool"`)
				So(nets[1].SSID, ShouldEqual, `"AP1"`)
				So(nets[2].SSID, ShouldEqual, `"AP2`)
			})
		})

	})
}

var scanResults = []byte(`bssid / frequency / signal level / flags / ssid
d0:7a:b5:31:23:a0	2472	-30	[WPA2-PSK-CCMP][WPS][ESS]	AP0
00:1f:1f:37:42:d9	2442	-37	[WPA2-PSK-CCMP][ESS]	\xc3\x80ndroid \xf0\x9f\x98\xab\xe2\x99\xa5\xef\xb8\x8f\xf0\x9f\x91\x8d
24:00:ba:f8:65:df	2412	-77	[WPA-PSK-CCMP+TKIP][WPA2-PSK-CCMP+TKIP][WPS][ESS]	AP2
`)

var scanResultsQuote = []byte(`bssid / frequency / signal level / flags / ssid
d0:7a:b5:31:23:a0	2472	-30	[WPA2-PSK-CCMP][WPS][ESS]	AP0 \"cool\"
00:1f:1f:37:42:d9	2442	-37	[WPA2-PSK-CCMP][ESS]	\"AP1\"
24:00:ba:f8:65:df	2412	-77	[WPA-PSK-CCMP+TKIP][WPA2-PSK-CCMP+TKIP][WPS][ESS]	\"AP2
`)

var networkList = []byte(`network id / ssid / bssid / flags
0	AP0	any	
1	\xc3\x80ndroid \xf0\x9f\x98\xab\xe2\x99\xa5\xef\xb8\x8f\xf0\x9f\x91\x8d	any	[DISABLED]
2	AP2	any	[DISABLED]`)

var networkListQuote = []byte(`network id / ssid / bssid / flags
0	AP0 \"cool\"	any	
1	\"AP1\"	any	[DISABLED]
2	\"AP2	any	[DISABLED]`)
