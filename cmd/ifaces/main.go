package main

import (
	"flag"
	"fmt"

	"github.com/plantiga/go-wireless"
)

func main() {
	var sysfs, wpa bool
	flag.BoolVar(&sysfs, "s", false, "find wireless interfaces by SysFS")
	flag.BoolVar(&wpa, "w", false, "find wireless interfaces by open WPA sockets")
	flag.Parse()

	var ifaces []string
	switch {
	case sysfs:
		ifaces = wireless.InterfacesFromSysfs()
	case wpa:
		ifaces = wireless.InterfacesFromWPARunDir()
	default:
		ifaces = wireless.Interfaces()
	}

	for _, i := range ifaces {
		fmt.Println(i)
	}
}
