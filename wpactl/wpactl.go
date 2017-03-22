package wpactl

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type WPAController struct {
	Interface    string
	EventChannel chan WPAEvent

	lsockname              string
	conn                   *net.UnixConn
	currentCommandResponse chan string
}

type WPAEvent struct {
	Name      string
	Arguments map[string]string
}

type WPANetwork struct {
	Id    int
	Freq  int
	RSSI  int
	BSSID string
	SSID  string
	ESSID string
	Flags string
}

func NewController(iface string) (*WPAController, error) {
	c := &WPAController{Interface: iface}
	err := c.Initialise()
	if err != nil {
		return nil, err
	} else {
		return c, nil
	}
}

func (c *WPAController) Cleanup() {
	c.conn.Close()
	os.Remove(c.lsockname)
}

func (c *WPAController) Initialise() error {
	addr, err := net.ResolveUnixAddr("unixgram", "/var/run/wpa_supplicant/"+c.Interface)
	if err != nil {
		return err
	}

	c.lsockname = fmt.Sprintf("/tmp/wpa_ctrl_%d", os.Getpid())
	laddr, err := net.ResolveUnixAddr("unixgram", c.lsockname)
	if err != nil {
		return err
	}

	c.conn, err = net.DialUnix("unixgram", laddr, addr)
	if err != nil {
		return err
	}

	log.Println("Local addr: ", c.conn.LocalAddr())

	c.EventChannel = make(chan WPAEvent, 128)
	c.currentCommandResponse = make(chan string, 1)

	go func() {
		buf := make([]byte, 2048)
		for {
			bytesRead, err := c.conn.Read(buf)
			if err != nil {
				log.Println("Error:", err)
			} else {
				msg := string(buf[:bytesRead])
				if msg[0] == '<' {
					// event message

					if strings.Index(msg, "<3>CTRL-") == 0 {
						// control event, sent to the channel
						reader := csv.NewReader(strings.NewReader(msg))
						reader.Comma = ' '
						reader.LazyQuotes = true
						reader.TrimLeadingSpace = false
						parts, err := reader.Read()
						if err != nil {
							log.Println("Error during parsing:", err)
						}
						if len(parts) == 0 {
							continue
						}

						event := WPAEvent{Name: parts[0][3:], Arguments: make(map[string]string)}
						for _, record := range parts[1:] {
							if strings.Index(record, "=") != -1 {
								nvs := strings.SplitN(record, "=", 2)
								event.Arguments[nvs[0]] = nvs[1]
							}
						}

						c.EventChannel <- event
					}
				} else {
					c.currentCommandResponse <- msg
				}
			}
		}
	}()

	err = c.SendCommandBool("ATTACH")
	if err != nil {
		return err
	}

	return nil // ok
}

func (c *WPAController) SendCommand(command string) (string, error) {
	log.Println("<<<", command)
	_, err := c.conn.Write([]byte(command))
	if err != nil {
		return "", err
	}

	resp := <-c.currentCommandResponse
	log.Println(">>>", resp)
	return resp, nil
}

func (c *WPAController) SendCommandBool(command string) error {
	resp, err := c.SendCommand(command)
	if err != nil {
		return err
	}
	if resp != "OK\n" {
		return errors.New(resp)
	}
	return nil
}

func (c *WPAController) SendCommandInt(command string) (int, error) {
	resp, err := c.SendCommand(command)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(strings.TrimSpace(resp))
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (c *WPAController) ListNetworks() ([]WPANetwork, error) {
	resp, err := c.SendCommand("LIST_NETWORKS")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(resp, "\n")

	num_networks := len(lines) - 1
	networks := make([]WPANetwork, num_networks)
	valid_networks := 0

	for _, line := range lines[1:] {
		fields := strings.Split(line, "\t")
		id, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != 4 {
			continue
		}
		networks[valid_networks].Id = id
		networks[valid_networks].SSID = fields[1]
		networks[valid_networks].ESSID = fields[2]
		networks[valid_networks].Flags = fields[3]
		valid_networks += 1
	}

	return networks[:valid_networks], nil
}

func (c *WPAController) ScanNetworks() ([]WPANetwork, error) {
	resp, err := c.SendCommand("SCAN")
	if err != nil {
		return nil, err
	}
	if resp != "OK\n" {
		return nil, errors.New(resp)
	}
	resp, err = c.SendCommand("SCAN_RESULTS")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(resp, "\n")

	num_networks := len(lines) - 1
	networks := make([]WPANetwork, num_networks)
	valid_networks := 0

	for _, line := range lines[1:] {
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			continue
		}
		freq, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		rssi, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}
		networks[valid_networks].BSSID = fields[0]
		networks[valid_networks].Freq = freq
		networks[valid_networks].RSSI = rssi
		networks[valid_networks].Flags = fields[3]
		networks[valid_networks].SSID = fields[4]
		valid_networks += 1
	}

	return networks[:valid_networks], nil
}

func (c *WPAController) AddNetwork() (int, error) {
	return c.SendCommandInt("ADD_NETWORK")
}

func (c *WPAController) SetNetworkSettingRaw(networkId int, variable string, value string) error {
	return c.SendCommandBool(fmt.Sprintf("SET_NETWORK %d %s %s", networkId, variable, value))
}

func (c *WPAController) SetNetworkSettingString(networkId int, variable string, value string) error {
	return c.SetNetworkSettingRaw(networkId, variable, fmt.Sprintf("\"%s\"", value))
}

func (c *WPAController) GetNetworkSetting(networkId int, variable string) (string, error) {
	return c.SendCommand(fmt.Sprintf("GET_NETWORK %d %s", networkId, variable))
}

func (c *WPAController) SelectNetwork(networkId int) error {
	return c.SendCommandBool(fmt.Sprintf("SELECT_NETWORK %d", networkId))
}

func (c *WPAController) EnableNetwork(networkId int) error {
	return c.SendCommandBool(fmt.Sprintf("ENABLE_NETWORK %d", networkId))
}

func (c *WPAController) DisableNetwork(networkId int) error {
	return c.SendCommandBool(fmt.Sprintf("DISABLE_NETWORK %d", networkId))
}

func (c *WPAController) RemoveNetwork(networkId int) error {
	return c.SendCommandBool(fmt.Sprintf("REMOVE_NETWORK %d", networkId))
}

func (c *WPAController) ReloadConfiguration() error {
	return c.SendCommandBool(fmt.Sprintf("RECONFIGURE"))
}

func (c *WPAController) SaveConfiguration() error {
	return c.SendCommandBool(fmt.Sprintf("SAVE_CONFIG"))
}

/*
func main() {
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

	networks, err := wpa_ctl.ScanNetworks()
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

	for {
		event := <- wpa_ctl.EventChannel
		//log.Println(event)
		switch event.Name {
			case "CTRL-EVENT-DISCONNECTED":
				log.Println("Disconnected")
			case "CTRL-EVENT-CONNECTED":
				log.Println("Connected")
			case "CTRL-EVENT-SSID-TEMP-DISABLED":
				log.Println("InvalidKey")
		}
	}
}
*/
