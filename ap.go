package wireless

// AP represents an access point seen by the scan networks command
type AP struct {
	ID        int      `json:"id"`
	RSSI      int      `json:"rssi"`
	BSSID     string   `json:"bssid"`
	SSID      string   `json:"ssid"`
	ESSID     string   `json:"essid"`
	Flags     []string `json:"flags"`
	Signal    int      `json:"signal"`
	Frequency int      `json:"frequency"`
}

// APs models a collection of access points
type APs []AP

// FindBySSID will find an AP by the given SSID or return false
func (nets APs) FindBySSID(ssid string) (AP, bool) {
	for _, n := range nets {
		if n.SSID == ssid {
			return n, true
		}
	}

	return AP{}, false
}
