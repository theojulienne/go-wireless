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

type Conn struct {
	Interface    string
	EventChannel chan Event

	lsockname              string
	conn                   *net.UnixConn
	currentCommandResponse chan string
}

type Event struct {
	Name      string
	Arguments map[string]string
}

type AP struct {
	Id    int
	Freq  int
	RSSI  int
	BSSID string
	SSID  string
	ESSID string
	Flags string
}

type Network struct {
	Id    int
	SSID  string
	BSSID string
	PSK   string
}

func Dial(iface string) (*Conn, error) {
	c := &Conn{Interface: iface}
	err := c.Initialise()
	if err != nil {
		return nil, err
	} else {
		return c, nil
	}
}

func (c *Conn) Close() {
	c.conn.Close()
	os.Remove(c.lsockname)
}

func (c *Conn) listen() {
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

					event := Event{Name: parts[0][3:], Arguments: make(map[string]string)}
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
}

func (c *Conn) init() error {
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

	go c.listen()

	err = c.SendCommandBool("ATTACH")
	if err != nil {
		return err
	}

	return nil // ok
}

func (c *Conn) SendCommand(command string) (string, error) {
	log.Println("<<<", command)
	_, err := c.conn.Write([]byte(command))
	if err != nil {
		return "", err
	}

	resp := <-c.currentCommandResponse
	log.Println(">>>", resp)
	return resp, nil
}

func (c *Conn) SendCommandBool(command string) error {
	resp, err := c.SendCommand(command)
	if err != nil {
		return err
	}
	if resp != "OK\n" {
		return errors.New(resp)
	}
	return nil
}

func (c *Conn) SendCommandInt(command string) (int, error) {
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
