package main

import (
	"fmt"
	"os"

	"github.com/theojulienne/go-wireless"
)

func main() {
	iface, ok := wireless.DefaultInterface()
	if !ok {
		panic("no wifi cards on the system")
	}
	fmt.Printf("Using interface: %s\n", iface)

	wc, err := wireless.NewClient(iface)
	if err != nil {
		panic(err)
	}
	defer wc.Close()

	net, err := wc.Connect(wireless.NewNetwork(os.Args[1], os.Args[2]))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to " + net.SSID)
}
