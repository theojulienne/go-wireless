package api

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theojulienne/go-wireless"
)

func notImplemented(c *gin.Context) {
	c.AbortWithStatus(501)
}

func json(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}

func errStatus(err error) int {
	switch err {
	case wireless.ErrFailBusy:
		return 409
	case wireless.ErrNoIdentifier, wireless.ErrAssocRejected, wireless.ErrSSIDNotFound:
		return 400
	default:
		return 500
	}
}

func listInterfaces(c *gin.Context) {
	c.JSON(200, wireless.Interfaces())
}

func listAccesPoints(c *gin.Context) {
	iface := c.Param("iface")

	wc, err := wireless.NewClient(iface)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(errStatus(err), json(err))
		return
	}
	defer wc.Close()

	wc.ScanTimeout = 3 * time.Second

	aps, err := wc.Scan()
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(errStatus(err), json(err))
		return
	}

	c.JSON(200, aps)
}

func listNetworks(c *gin.Context) {
	iface := c.Param("iface")

	wc, err := wireless.NewClient(iface)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(errStatus(err), json(err))
		return
	}
	defer wc.Close()

	nets, err := wc.Networks()
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(errStatus(err), json(err))
		return
	}

	c.JSON(200, nets)
}

func getInterface(c *gin.Context) {
	iface := c.Param("iface")

	wc, err := wireless.NewClient(iface)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(errStatus(err), json(err))
		return
	}
	defer wc.Close()

	state, err := wc.Status()
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(errStatus(err), json(err))
		return
	}

	c.JSON(200, state)
}

func addNetwork(c *gin.Context) {
	iface := c.Param("iface")
	if iface == "" {
		c.AbortWithStatusJSON(400, json(errors.New("must specify interface name")))
		return
	}

	disable := false
	if v, ok := c.GetQuery("disable"); ok && v == "1" {
		disable = true
	}
	if v, ok := c.GetQuery("disabled"); ok && v == "1" {
		disable = true
	}

	force := false
	if v, ok := c.GetQuery("force"); ok && v == "1" {
		force = true
	}

	hidden := false
	if v, ok := c.GetQuery("hidden"); ok && v == "1" {
		hidden = true
	}

	connect := false
	if v, ok := c.GetQuery("connect"); ok && v == "1" {
		connect = true
	}
	if v, ok := c.GetQuery("connected"); ok && v == "1" {
		connect = true
	}

	wc, err := wireless.NewClient(iface)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(errStatus(err), json(err))
		return
	}
	defer wc.Close()

	var nw wireless.Network
	if err := c.BindJSON(&nw); err != nil {
		c.AbortWithStatusJSON(400, json(err))
		return
	}

	if hidden {
		nw.ScanSSID = true
	}

	if connect {
		disable = false
	}
	nw.Disable(disable)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wc.WithContext(ctx)

	wc.LoadConfig()
	newNet, err := wc.AddNetwork(nw)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(errStatus(err), json(err))
		return
	}

	if connect {
		// disable all but the SSID we are trying to connect to
		disableMap := map[int]bool{}
		nets, _ := wc.Networks()
		for _, net := range nets {
			if net.SSID != newNet.SSID {
				disableMap[net.ID] = net.IsDisabled()
				wc.DisableNetwork(net.ID)
			}
		}

		// Ensure the network is selected
		time.Sleep(time.Second / 10)
		wc.SelectNetwork(newNet.ID)

		newNet, err = wc.Connect(newNet)
		if err != nil {
			c.Error(err)
			reEnable(disableMap, wc)

			if err == wireless.ErrSSIDNotFound && force {
				wc.SaveConfig()
				c.JSON(200, newNet)
				return
			}

			wc.RemoveNetwork(newNet.ID)
			c.AbortWithStatusJSON(errStatus(err), json(err))
			return
		}
	}

	nets, err := wc.Networks()
	if err != nil {
		c.AbortWithStatusJSON(500, json(errors.New("network wasn't added")))
		return
	}

	curr, connected := nets.FindCurrent()
	if !connected {
		c.AbortWithStatusJSON(503, json(errors.New("not connected to any SSID")))
		return
	}

	if curr.SSID != newNet.SSID {
		c.AbortWithStatusJSON(502, json(errors.New("connected to unexpected SSID "+curr.SSID)))
		return
	}

	wc.SaveConfig()
	c.JSON(200, newNet)
}

func reEnable(disableMap map[int]bool, wc *wireless.Client) {
	for id, disabled := range disableMap {
		if !disabled {
			wc.EnableNetwork(id)
		}
	}
}
