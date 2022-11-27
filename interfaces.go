package wireless

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Interfaces is a shortcut to the best known method for gathering the wireless
// interfaces from the current system
var Interfaces = InterfacesFromWPARunDir
var DefaultCtrlDir = "/var/run/wpa_supplicant"

// DefaultInterface will return the default wireless interface, being the first
// one returned from the Interfaces method
func DefaultInterface() (string, bool) {
	ifs := Interfaces()
	if len(ifs) == 0 {
		return "", false
	}

	return ifs[0], true
}

// InterfacesFromWPARunDir returns the interfaces that WPA Supplicant is currently running on
// by checking the sockets available in the run directory (/var/run/wpa_supplicant)
// however a different run directory can be specified as the basePath parameter
func InterfacesFromWPARunDir(basePath ...string) []string {
	s := []string{}
	base := DefaultCtrlDir
	if len(basePath) > 0 {
		base = basePath[0]
	}
	matches, _ := os.ReadDir(base)
	for _, iface := range matches {
		if !strings.HasPrefix(iface.Name(), "p2p") {
			s = append(s, iface.Name())
		}
	}

	return s
}

// InterfacesFromSysfs returns the wireless interfaces found in the SysFS (/sys/class/net)
func InterfacesFromSysfs() []string {
	s := []string{}
	base := "/sys/class/net"
	matches, _ := filepath.Glob(path.Join(base, "*"))

	//  look for the wireless folder in each interfces directory to determine if it is a wireless device
	for _, iface := range matches {
		if stat, err := os.Stat(path.Join(iface, "wireless")); err == nil && stat.IsDir() {
			s = append(s, path.Base(iface))
		}
	}

	return s
}
