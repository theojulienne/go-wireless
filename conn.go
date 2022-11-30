package wireless

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const CONN_MAX_LISTEN_BUFF = 3 * 1024 // Allow 3kB of buffer for listening to events

func init() {
	rand.Seed(time.Now().Unix())
}

// Conn represents a connection to a WPA supplicant control interface
type Conn struct {
	Interface string

	lsockname              string
	conn                   *net.UnixConn
	currentCommandResponse chan string

	subs []*Subscription
	log  *log.Logger

	quit chan bool
}

// Dial will dial the WPA control interface with the given
// interface name
func Dial(iface string) (*Conn, error) {
	c := &Conn{Interface: iface, log: log.New(ioutil.Discard, "", log.LstdFlags)}
	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// WithLogOutput will set the log output of the connection.  By default it is
// set to ioutil.Discard
func (c *Conn) WithLogOutput(w io.Writer) {
	c.log.SetOutput(w)
}

// Close will close the connection to the WPA control interface
func (c *Conn) Close() error {
	close(c.quit)
	c.conn.Close()
	return os.Remove(c.lsockname)
}

func (c *Conn) listen() {
	buf := make([]byte, CONN_MAX_LISTEN_BUFF)
	for {
		select {
		case <-c.quit:
			return
		default:
			bytesRead, err := c.conn.Read(buf)
			c.processMessage(bytesRead, buf, err)
		}
	}
}

func (c *Conn) processMessage(bytesRead int, data []byte, err error) {
	if err != nil {
		if !IsUseOfClosedNetworkConnectionError(err) {
			c.log.Println("Error:", err)
		}
		return
	}

	msg := string(data[:bytesRead])
	if msg[0] != '<' {
		c.currentCommandResponse <- msg
		return
	}

	var ev Event
	if strings.Index(msg, "<3>CTRL-") == 0 {
		ev, err = NewEventFromMsg(msg)
		if err != nil {
			c.log.Println("Error:", err)
			return
		}
	} else {
		ev.Name = "logs"
		ev.Arguments = map[string]string{"msg": msg[3:]}
	}

	c.publishEvent(ev)
}

func (c *Conn) init() error {
	addr, err := net.ResolveUnixAddr("unixgram", "/var/run/wpa_supplicant/"+c.Interface)
	if err != nil {
		return err
	}

	if v := os.Getenv("LOG"); v != "" {
		c.WithLogOutput(os.Stderr)
	}

	c.lsockname = fmt.Sprintf("/tmp/wpa_ctrl_%d_%d", os.Getpid(), rand.Intn(10000))
	laddr, err := net.ResolveUnixAddr("unixgram", c.lsockname)
	if err != nil {
		return err
	}

	c.conn, err = net.DialUnix("unixgram", laddr, addr)
	if err != nil {
		return err
	}

	c.log.Println("Local addr: ", c.conn.LocalAddr())

	c.currentCommandResponse = make(chan string, 1)
	c.quit = make(chan bool)

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
	c.log.Println("<<<", command)
	_, err := c.conn.Write([]byte(strings.Join(command, " ")))
	if err != nil {
		return "", err
	}

	for {
		select {
		case resp := <-c.currentCommandResponse:
			c.log.Println(">>>", resp)
			return resp, nil
		case <-ctx.Done():
			return "", ErrCmdTimeout

		}
	}
}

// SendCommandBool will send a command and return an error
// if the response was not OK
func (c *Conn) SendCommandBool(command ...string) error {
	return c.SendCommandBoolWithContext(context.Background(), command...)
}

func (c *Conn) SendCommandBoolWithContext(ctx context.Context, command ...string) error {
	resp, err := c.SendCommandWithContext(ctx, command...)
	if err != nil {
		return err
	}
	if resp != "OK\n" {
		return errors.New(strings.TrimSpace(resp))
	}
	return nil
}

// SendCommandInt will send a command where the response is expected to be an integer
func (c *Conn) SendCommandInt(command ...string) (int, error) {
	return c.SendCommandIntWithContext(context.Background(), command...)
}

func (c *Conn) SendCommandIntWithContext(ctx context.Context, command ...string) (int, error) {
	resp, err := c.SendCommandWithContext(ctx, command...)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(strings.TrimSpace(resp))
	if err != nil {
		return 0, err
	}
	return i, nil
}
