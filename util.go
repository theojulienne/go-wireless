package wireless

import (
	"bytes"
	"encoding/csv"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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

	flags := strings.Split(s, "][")
	if len(flags) == 1 && flags[0] == "" {
		return []string{}
	}

	return flags
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
			BSSID:     bssid.String(),
			SSID:      rec[4],
			Frequency: fr,
			Signal:    ss,
			Flags:     parseFlags(rec[3]),
		})
	}

	return aps, nil
}

func quote(s string) string {
	return `"` + s + `"`
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func unquote(s string) string {
	return strings.Trim(s, `"`)
}
