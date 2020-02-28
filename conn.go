package wireless

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Conn represents a connection to a WPA supplicant control interface
type Conn struct {
	Interface string

	lsockname              string
	conn                   *net.UnixConn
	currentCommandResponse chan string

	subs []Subscription
}

// Dial will dial the WPA control interface with the given
// interface name
func Dial(iface string) (*Conn, error) {
	c := &Conn{Interface: iface}
	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Close will close the connection to the WPA control interface
func (c *Conn) Close() {
	c.conn.Close()
	os.Remove(c.lsockname)
}

func (c *Conn) listen() {
	buf := make([]byte, 2048)
	for {
		bytesRead, err := c.conn.Read(buf)
		if err != nil {
			if IsUseOfClosedNetworkConnectionError(err) {
				continue
			}
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

					go c.publishEvent(event)
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

	c.currentCommandResponse = make(chan string, 1)

	go c.listen()

	err = c.SendCommandBool("ATTACH")
	if err != nil {
		return err
	}

	return nil // ok
}

// SendCommand will call SendCommandWithContext with a 2 second timeout
func (c *Conn) SendCommand(command ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	return c.SendCommandWithContext(ctx, command...)
}

// SendCommandWithContext will send the command with a context
func (c *Conn) SendCommandWithContext(ctx context.Context, command ...string) (string, error) {
	log.Println("<<<", command)
	_, err := c.conn.Write([]byte(strings.Join(command, " ")))
	if err != nil {
		return "", err
	}

	for {
		select {
		case resp := <-c.currentCommandResponse:
			log.Println(">>>", resp)
			return resp, nil
		case <-ctx.Done():
			return "", ErrCmdTimeout

		}
	}
}

// SendCommandBool will send a command and return an error
// if the response was not OK
func (c *Conn) SendCommandBool(command ...string) error {
	resp, err := c.SendCommand(command...)
	if err != nil {
		return err
	}
	if resp != "OK\n" {
		return errors.New(resp)
	}
	return nil
}

// SendCommandInt will send a command where the response is expected to be an integer
func (c *Conn) SendCommandInt(command ...string) (int, error) {
	resp, err := c.SendCommand(command...)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(strings.TrimSpace(resp))
	if err != nil {
		return 0, err
	}
	return i, nil
}
