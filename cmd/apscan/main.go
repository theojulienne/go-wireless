package main

import (
	"fmt"

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

	aps, err := wc.Scan()
	if err != nil {
		panic(err)
	}

	for _, ap := range aps {
		fmt.Println(ap.SSID)
	}
}
