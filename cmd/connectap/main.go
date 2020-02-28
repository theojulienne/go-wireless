package main

import (
	"fmt"
	"os"

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

	if err := wc.Connect(wireless.NewNetwork(os.Args[1], os.Args[2])); err != nil {
		panic(err)
	}

	fmt.Printf("Connected to " + os.Args[1])
}
