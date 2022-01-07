package conf

import (
	"strings"

	"github.com/plantiga/go-wireless"
)

// Network models a network in a WPA config file
type Network struct {
	wireless.Network
}

func exVal(v string) string {
	return strings.Trim(strings.TrimSpace(v), `"`)
}

// NewNetworkFromLines will parse the lines from a network block and
// turn them into a network object
func NewNetworkFromLines(id int, lines []string) Network {
	net := Network{}
	net.ID = id

	for _, l := range lines {
		bits := strings.Split(l, "=")
		if len(bits) != 2 {
			continue
		}

		f := bits[0]
		v := bits[1]

		switch f {
		case "scan_ssid":
			if exVal(v) == "1" {
				net.ScanSSID = true
			}

		case "ssid":
			net.SSID = exVal(v)

		case "psk":
			net.PSK = exVal(v)

		case "key_mgmt":
			net.KeyMgmt = exVal(v)

			if exVal(v) == "1" {
				net.Disable(true)
			}

		case "id_str":
			net.IDStr = exVal(v)
		}
	}

	return net
}

// Render will render the network as a network block in the config file
func (net Network) Render() string {
	lines := []string{}

	lines = append(lines, "ssid="+quote(net.SSID))
	if net.PSK != "" {
		lines = append(lines, "psk="+quote(net.PSK))
	}

	if net.IsDisabled() {
		lines = append(lines, "disabled=1")
	}

	if net.ScanSSID {
		lines = append(lines, "scan_ssid=1")
	}

	if net.KeyMgmt != "" {
		lines = append(lines, "key_mgmt="+net.KeyMgmt)
	}

	s := "network={\n"
	for _, l := range lines {
		s += "  " + l
	}
	s += "}\n"

	return s
}

func quote(v string) string {
	return `"` + v + `"`
}
