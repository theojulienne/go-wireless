package main

import "github.com/theojulienne/go-wireless/iwlib"
//import "fmt"

func main() {
	iwlib.GetWirelessNetworks("wlan0")
}
