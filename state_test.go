package wireless

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseState(t *testing.T) {
	Convey("given a network state", t, func() {
		Convey("when it is parsed", func() {
			s := NewState(stateResults)

			So(s.BSSID, ShouldEqual, "74:ac:b9:bd:2f:4a")
			So(s.SSID, ShouldEqual, "yurt")
			So(s.ID, ShouldEqual, "0")
			So(s.Mode, ShouldEqual, "station")
			So(s.PairwiseCipher, ShouldEqual, "CCMP")
			So(s.GroupCipher, ShouldEqual, "CCMP")
			So(s.KeyManagement, ShouldEqual, "WPA2-PSK")
			So(s.WpaState, ShouldEqual, "COMPLETED")
			So(s.IPAddress, ShouldEqual, "192.168.0.108")
			So(s.Address, ShouldEqual, "b8:27:eb:71:3d:f5")
			So(s.UUID, ShouldEqual, "9cea0881-719c-5e06-9049-fe64900ef0d3")
		})

		Convey("when unicode is parsed", func() {
			s := NewState(unicodeStateResults)

			So(s.SSID, ShouldEqual, "Àndroid=😫♥️👍")
		})
	})
}

var unicodeStateResults = `
bssid=74:ac:b9:bd:2f:4a
freq=2437
ssid=\xc3\x80ndroid=\xf0\x9f\x98\xab\xe2\x99\xa5\xef\xb8\x8f\xf0\x9f\x91\x8d
id=0
mode=station
pairwise_cipher=CCMP
group_cipher=CCMP
key_mgmt=WPA2-PSK
wpa_state=COMPLETED
ip_address=192.168.0.108
p2p_device_address=ce:b4:c0:a1:31:67
address=b8:27:eb:71:3d:f5
uuid=9cea0881-719c-5e06-9049-fe64900ef0d3
`

var stateResults = `
bssid=74:ac:b9:bd:2f:4a
freq=2437
ssid=yurt
id=0
mode=station
pairwise_cipher=CCMP
group_cipher=CCMP
key_mgmt=WPA2-PSK
wpa_state=COMPLETED
ip_address=192.168.0.108
p2p_device_address=ce:b4:c0:a1:31:67
address=b8:27:eb:71:3d:f5
uuid=9cea0881-719c-5e06-9049-fe64900ef0d3
`
