# go-wireless

[![Go Report Card](https://goreportcard.com/badge/github.com/theojulienne/go-wireless)](https://goreportcard.com/report/github.com/theojulienne/go-wireless) [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/theojulienne/go-wireless) ![Go](https://github.com/theojulienne/go-wireless/workflows/Go/badge.svg)

A way to interact with the Wireless interfaces on a Linux machine using WPA Supplicant.

## Requirements

Requires a running wpa_supplicant with control interface at `/var/run/wpa_supplicant` (which is usually 
a symlink to `/run/wpa_supplicant`).  This requires the config file to contain the line:

```
ctrl_interface=DIR=/run/wpa_supplicant GROUP=wheel
```

Or for the `wpa_supplicant` instance to be running with the `-O /run/wpa_supplicant` argument.

You will probably also need to be running as root unless you are in the specified group (`wheel` in the above example).

# Usage

There are two main objects to interact with:

* `Conn` - allows running arbitrary commands and listening to events
* `Client` - uses `Conn` to do things for you like scanning for and connecting to APs

Examples of the usage can be found in the `cmd` directory as standalone commands:

* `apscan` - similar to your `iwlist wlan0 scan`
* `connectap` - connect to an access point
* `currentap` - get the access point that is currently connected
* `ifaces` - show the wireless interfaces
* `wifistate` - dump the current wifi state as JSON
* `wpalogs` - print logs as they happen
* `wpaspy` - print events as they happen
* `wpapi` - an HTTP API that can scan and connect to APs

Get a list of wifi cards attached:

```golang
ifaces := wireless.Interfaces()
```

From there you can use the client:

```golang
wc, err := wireless.NewClient("wlan0")
defer wc.Close()
```

Get a list of APs that are in range:

```golang
aps, err := wc.Scan()
fmt.Println(aps, err)
ap, ok := wireless.APs(aps).FindBySSID("CIA Predator Drone 237A")
```

Get a list of known networks (**note:** the password cannot be retrieved so are not populated):

```golang
nets, err := wc.Networks()
fmt.Println(nets, err)
```

Connect to networks:

```golang
net := NewNetwork("FBI Surveillance Van #4", "secretpass")
net, err := wc.Connect(net)
```

Disable networks:

```golang
nets, err:= wc.Networks()
net, err := net, ok := wireless.Networks(nets).FindBySSID("FBI Surveillance Van #4")
net.Disable(true)
net, err := wc.UpdateNetwork(net)
```

Subscribe to events by getting a raw connection object:

```golang
conn, _ := wireless.Dial("wlp2s0")
sub := conn.Subscribe(wireless.EventConnected, wireless.EventAuthReject, wireless.EventDisconnected)

ev := <-sub.Next()
switch ev.Name {
	case wireless.EventConnected:
		fmt.Println(ev.Arguments)
	case wireless.EventAuthReject:
		fmt.Println(ev.Arguments)
	case wireless.EventDisconnected:
		fmt.Println(ev.Arguments)
}
```

Check the status of the connection:

```golang
st, err := wc.Status()
fmt.Printf("%+v\n", st)
```

Use a timeout for scanning APs and general command timeouts:

```golang
wc.ScanTimeout = time.Second * 3
wc.CmdTimeout = time.Second
```

You can also set a context to use for the latter instead:

```golang
ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
wc.WithContext(ctx)
```

However that may not be useful to you unless you have a shortlived client.  If you want to use a context for groups of operations you can give a closure:

```golang
wc.WithContext(ctx, func(wc *wireless.Client) {
	err := wc.Disconnect()
	err = wc.RemoveNetwork(0)
	net := NewNetwork("DHS Stingray #27", "secretpass")
	net, err := wc.Connect(net)
})
```

Or if you want to use it for a single operation most have `*WithContext` alternatives:

```golang
wc.ConnectWithContext(ctx, net)
```

## API

There is an API that can be used with [gin](https://github.com/gin-gonic/gin):

```golang
r := gin.Default()
api.SetupRoutes(r)
r,Serve(":8080")
```

## Endpoints

- [x] `GET /interfaces`
- [ ] `GET /interfaces/:iface`
- [ ] `PUT /interfaces/:iface`
- [x] `GET /interfaces/:iface/aps`
- [x] `GET /interfaces/:iface/networks`
- [ ] `POST /interfaces/:iface/networks`
- [ ] `PUT /interfaces/:iface/networks/:id_or_idstr`
- [ ] `GET /interfaces/:iface/networks/:id_or_idstr`
- [ ] `DELETE /interfaces/:iface/networks/:id_or_idstr`
