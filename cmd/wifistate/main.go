package main

import (
	"encoding/json"
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

	st, err := wc.Status()
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
