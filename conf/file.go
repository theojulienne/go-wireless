package conf

import (
	"io/ioutil"
	"strings"
)

// DefaultPath is the default path for WPA supplicant config file
const DefaultPath = "/etc/wpa_supplicant/wpa_supplicant.conf"

// File represents the config file
type File struct {
	path     string
	Preamble []string
	Networks []Network
}

// Open will open and parse the config file at the given path
func Open(path string) (*File, error) {
	f := &File{path: path}
	return f, f.Load()
}

// Path will return the path of the file
func (f *File) Path() string {
	return f.path
}

// Load will parse the file contents into Preamble and Networks
func (f *File) Load() error {
	data, err := ioutil.ReadFile(f.path)
	if err != nil {
		return err
	}

	f.Preamble, f.Networks = parseFile(data)
	return nil
}

// Save will save the config to the file path
func (f *File) Save() error {
	return ioutil.WriteFile(f.path, renderFile(f), 0600)
}

func renderFile(file *File) []byte {
	s := ""
	for _, l := range file.Preamble {
		s += l
	}

	for _, n := range file.Networks {
		s += n.Render()
	}

	return []byte(s)
}

func parseFile(data []byte) (pre []string, nets []Network) {
	var inNetworkSection bool
	var inPreamble = true
	var netID = -1
	var netLines []string

	for _, l := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(l, "network={") {
			inPreamble = false
			inNetworkSection = true
			netID++
			continue
		}

		if strings.HasPrefix(l, "}") && inNetworkSection {
			inNetworkSection = false
			nets = append(nets, NewNetworkFromLines(netID, netLines))
			continue
		}

		if inPreamble {
			pre = append(pre, l)
		}

		if inNetworkSection {
			netLines = append(netLines, l)
		}
	}

	return
}
