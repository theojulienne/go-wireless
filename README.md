# go-wireless

A way to interact with the Wireless interfaces on a Linux machine using WPA Supplicant.

## Requirements

Requires a running wpa_supplicant with control interface at `/var/run/wpa_supplicant`.

# Usage

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/theojulienne/go-wireless)

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
```

Get a list of known networks:

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
net, err := nets.Find("FBI Surveillance Van #4")
net.Disable(true)
net, err := wc.UpdateNetwork(net)
```

Subscribe to events:

```golang
sub := wc.Subscribe(wireless.EventConnected, wireless.EventAuthReject, wireless.EventDisconnected)

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

# API

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