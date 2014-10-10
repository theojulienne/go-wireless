// +build !linux

package iwlib

type WirelessScanResult struct {
	SSID string
}

func GetWirelessNetworks(iface string) ([]WirelessScanResult, error) {
	return []WirelessScanResult{
		WirelessScanResult{"FakeNetwork"},
	}, nil
}
