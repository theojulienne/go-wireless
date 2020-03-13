package main

import (
	"flag"
	"fmt"

	"github.com/theojulienne/go-wireless"
)

func main() {
	var iface string
	flag.StringVar(&iface, "i", "", "interface to use")
	flag.Parse()

	if iface == "" {
		var ok bool
		iface, ok = wireless.DefaultInterface()
		if !ok {
			panic("no wifi cards on the system")
		}
	}
	fmt.Printf("Using interface: %s\n", iface)

	wc, err := wireless.NewClient(iface)
	if err != nil {
		panic(err)
	}
	defer wc.Close()

	nets, err := wc.Networks()
	if err != nil {
		panic(err)
	}

	curr, ok := wireless.Networks(nets).FindCurrent()
	if !ok {
		fmt.Println("not connected")
		return
	}

	fmt.Println(curr.SSID)
}
