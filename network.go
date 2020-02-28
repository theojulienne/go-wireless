package wireless

import (
	"bytes"
	"encoding/csv"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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

// Network represents a known network
type Network struct {
	ID    int
	SSID  string
	BSSID string
	PSK   string
	Flags []string
}

func parseNetwork(b []byte) ([]Network, error) {
	i := bytes.Index(b, []byte("\n"))
	if i > 0 {
		b = b[i:]
	}

	r := csv.NewReader(bytes.NewReader(b))
	r.Comma = '\t'
	r.FieldsPerRecord = 4

	recs, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	nts := []Network{}
	for _, rec := range recs {
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			return nil, errors.Wrap(err, "parse id")
		}

		nts = append(nts, Network{
			ID:    id,
			SSID:  rec[1],
			BSSID: rec[2],
			Flags: parseFlags(rec[3]),
		})
	}

	return nts, nil
}

func parseFlags(s string) []string {
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	return strings.Split(s, "][")
}

// AP represents an access point seen by the scan networks command
type AP struct {
	ID             int
	Freq           int
	RSSI           int
	BSSID          net.HardwareAddr
	SSID           string
	ESSID          string
	Flags          []string
	SignalStrength int
	Frequency      int
}

func parseAP(b []byte) ([]AP, error) {
	i := bytes.Index(b, []byte("\n"))
	if i > 0 {
		b = b[i:]
	}

	r := csv.NewReader(bytes.NewReader(b))
	r.Comma = '\t'
	r.FieldsPerRecord = 5

	recs, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	aps := []AP{}
	for _, rec := range recs {
		bssid, err := net.ParseMAC(rec[0])
		if err != nil {
			return nil, errors.Wrap(err, "parse mac")
		}

		fr, err := strconv.Atoi(rec[1])
		if err != nil {
			return nil, errors.Wrap(err, "parse frequency")
		}

		ss, err := strconv.Atoi(rec[2])
		if err != nil {
			return nil, errors.Wrap(err, "parse signal strength")
		}

		aps = append(aps, AP{
			BSSID:          bssid,
			SSID:           rec[4],
			Frequency:      fr,
			SignalStrength: ss,
			Flags:          parseFlags(rec[3]),
		})
	}

	return aps, nil
}
