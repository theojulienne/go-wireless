package main

import "github.com/theojulienne/go-iwlib/iwlib"
//import "fmt"

func main() {
	iwlib.GetWirelessNetworks("wlan0")
}
