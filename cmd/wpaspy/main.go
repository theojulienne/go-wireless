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

	sub := wc.Subscribe()

	for {
		ev := <-sub.Next()
		fmt.Println(ev.Name, ev.Arguments)
	}
}
