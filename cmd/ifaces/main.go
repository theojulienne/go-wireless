package main

import (
	"flag"
	"fmt"

	"github.com/theojulienne/go-wireless"
)

func main() {
	var sysfs bool
	flag.BoolVar(&sysfs, "s", false, "find wireless interfaces by SysFS")
	flag.Parse()

	var ifaces []string
	switch {
	case sysfs:
		ifaces = wireless.SysFSInterfaces()
	default:
		ifaces = wireless.Interfaces()
	}

	for _, i := range ifaces {
		fmt.Println(i)
	}
}
