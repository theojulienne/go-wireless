package wireless

import "strings"

// This file contains components from github.com/brlbil/wpaclient
//
// Copyright (c) 2017 Birol Bilgin
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// nd/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// UT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// NewNamedNetwork will create a new network with the given parameters
func NewNamedNetwork(name, ssid, psk string) Network {
	n := Network{IDStr: name, SSID: ssid, PSK: psk}
	if psk == "" {
		n.KeyMgmt = "NONE"
	}
	return n
}

// NewNetwork will create a new network with the given parameters
func NewNetwork(ssid, psk string) Network {
	return NewNamedNetwork(ssid, ssid, psk)
}

// NewOpenNetwork will return a new open network
func NewOpenNetwork(ssid string) Network {
	return NewNamedNetwork(ssid, ssid, "")
}

// NewDisabledNetwork will create a new disabled network with the given parameters
func NewDisabledNetwork(ssid, psk string) Network {
	n := NewNamedNetwork(ssid, ssid, psk)
	n.Flags = append(n.Flags, "DISABLED")
	return n
}

// Network represents a known network
type Network struct {
	Known bool `json:"known"`

	ID       int      `json:"id"`
	IDStr    string   `json:"id_str"`
	KeyMgmt  string   `json:"key_mgmt"`
	SSID     string   `json:"ssid"`
	BSSID    string   `json:"bssid"`
	ScanSSID bool     `json:"scan_ssid"`
	PSK      string   `json:"psk"`
	Flags    []string `json:"flags"`
}

// IsDisabled will return true if the network is disabled
func (net Network) IsDisabled() bool {
	for _, f := range net.Flags {
		if f == "DISABLED" {
			return true
		}
	}
	return false
}

// IsCurrent will return true if the network is the currently active one
func (net Network) IsCurrent() bool {
	for _, f := range net.Flags {
		if f == "CURRENT" {
			return true
		}
	}
	return false
}

type attributeGetter interface {
	GetNetworkAttr(int, string) (string, error)
}

func (net *Network) populateAttrs(cl attributeGetter) error {
	v, err := cl.GetNetworkAttr(net.ID, "ssid")
	if err != nil {
		return err
	}
	net.SSID = unquote(v)

	v, err = cl.GetNetworkAttr(net.ID, "id_str")
	if err != nil {
		return err
	}
	net.IDStr = unquote(v)

	v, err = cl.GetNetworkAttr(net.ID, "key_mgmt")
	if err != nil {
		return err
	}
	net.KeyMgmt = v

	v, err = cl.GetNetworkAttr(net.ID, "scan_ssid")
	if err != nil {
		return err
	}
	if v == "1" {
		net.ScanSSID = true
	}

	v, err = cl.GetNetworkAttr(net.ID, "disabled")
	if err != nil {
		return err
	}
	if v == "1" {
		net.Flags = append(net.Flags, "DISABLED")
	}

	return nil
}

// Disable or enabled the network
func (net *Network) Disable(on bool) {
	var idx int
	var found bool
	for i, f := range net.Flags {
		if f == "DISABLED" {
			found = true
			idx = i
			break
		}
	}

	if on && !found {
		net.Flags = append(net.Flags, "DISABLED")
		return
	}

	net.Flags = append(net.Flags[:idx], net.Flags[idx:]...)
}

// Networks models a collection of networks
type Networks []Network

// FindByIDStr will find a network by the given ID Str or return false
func (nets Networks) FindByIDStr(idStr string) (Network, bool) {
	for _, n := range nets {
		if n.IDStr == idStr {
			return n, true
		}
	}

	return Network{}, false
}

// FindBySSID will find a network by the given SSID or return false
func (nets Networks) FindBySSID(ssid string) (Network, bool) {
	for _, n := range nets {
		if n.SSID == ssid {
			return n, true
		}
	}

	return Network{}, false
}

// FindCurrent will find the current network or return false
func (nets Networks) FindCurrent() (Network, bool) {
	for _, n := range nets {
		if n.IsCurrent() {
			return n, true
		}
	}

	return Network{}, false
}

// Attributes return the attributes of the network as a list of strings,
// with the ability to set the separator or indentation
func (net Network) Attributes(sep, indent string) []string {
	lines := []string{}

	lines = append(lines, indent+"ssid"+sep+quote(net.SSID))
	lines = append(lines, indent+"id_str"+sep+quote(net.SSID))

	if net.PSK != "" {
		lines = append(lines, indent+"psk"+sep+quote(net.PSK))
	}

	if net.IsDisabled() {
		lines = append(lines, indent+"disabled"+sep+"1")
	}

	if net.ScanSSID {
		lines = append(lines, indent+"scan_ssid"+sep+"1")
	}

	// switch {
	// case net.KeyMgmt == "" && net.PSK == "":
	// 	lines = append(lines, indent+"key_mgmt"+sep+"NONE")

	// case net.KeyMgmt != "":
	// 	lines = append(lines, indent+"key_mgmt"+sep+net.KeyMgmt)
	// }

	return lines
}

func setCmds(net Network) []string {

	cmds := []string{}

	for _, attr := range net.Attributes(" ", "") {
		cmds = append(cmds, setCmdJoin(net.ID, attr))
	}

	return cmds
}

func setCmdJoin(idx int, bits ...string) string {
	return strings.Join(append([]string{CmdSetNetwork, itoa(idx)}, bits...), " ")
}
