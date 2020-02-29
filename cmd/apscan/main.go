package main

import (
	"fmt"

	"github.com/theojulienne/go-wireless"
)

func main() {
	ifaces := wireless.Interfaces()
	var iface string
	switch len(ifaces) {
	case 0:
		panic("no wifi cards on the system")
	default:
		iface = ifaces[0]
	}

	fmt.Printf("Using interface: " + iface)

	wc, err := wireless.NewClient(iface)
	if err != nil {
		panic(err)
	}
	defer wc.Close()

	aps, err := wc.Scan()
	if err != nil {
		panic(err)
	}

	for _, ap := range aps {
		fmt.Println(ap.SSID)
	}
}
