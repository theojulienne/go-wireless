iwlib
-----

Requires:
```
apt-get install libiw-dev
````

Example:
``` go
	networks, err := iwlib.GetWirelessNetworks("wlan0")
	if err != nil {
		fmt.Printf("Error retrieve wireless networks:", err)
		return
	}
	
	for _, network := range networks {
		fmt.Printf("SSID: %v\n", network.SSID)
	}
```

wpactl
------

Requires a running wpa_supplicant with control interface at `/var/run/wpa_supplicant`.

Example:
``` go
wpa_ctl, err := NewController("wlan0")
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer wpa_ctl.Cleanup()

	err = wpa_ctl.ReloadConfiguration()
	if err != nil {
		log.Fatal("Error:", err)
	}

	networks, err := wpa_ctl.ListNetworks()
	if err != nil {
		log.Fatal("Error retrieving networks:", err)
	}
	for _,network := range networks {
		log.Println("NET", network)
	}

	i, _ := wpa_ctl.AddNetwork()
	wpa_ctl.SetNetworkSettingString(i, "ssid", "helloworld")
	wpa_ctl.SetNetworkSettingString(i, "psk", "thisisnotarealpsk")
	wpa_ctl.SetNetworkSettingRaw(i, "scan_ssid", "1")
	wpa_ctl.SetNetworkSettingRaw(i, "key_mgmt", "WPA-PSK")
	wpa_ctl.SelectNetwork(i)
	wpa_ctl.SaveConfiguration()
	//

	for {
		event := <- wpa_ctl.EventChannel
		//log.Println(event)
		switch event.name {
			case "CTRL-EVENT-DISCONNECTED":
				log.Println("Disconnected")
			case "CTRL-EVENT-CONNECTED":
				log.Println("Connected")
			case "CTRL-EVENT-SSID-TEMP-DISABLED":
				log.Println("InvalidKey")
		}
	}
```