package wireless

// Client represents a wireless client
type Client struct {
	conn *Conn
}

// NewClient will create a new client by connecting to the
// given interface in WPA
func NewClient(iface string) (c *Client, err error) {
	c.conn, err = Dial(iface)
	if err != nil {
		return
	}

	return
}

// NewClientFromConn returns a new client from an already established connection
func NewClientFromConn(conn *Conn) (c *Client) {
	c.conn = conn
	return
}

// Close will close the client connection
func (cl *Client) Close() {
	cl.conn.Close()
}

// Scan will scan for networks and return the APs it finds
func (cl *Client) Scan() (nets []AP, err error) {
	err = cl.conn.SendCommandBool(CmdScan)
	if err != nil {
		return
	}

	results := cl.conn.Subscribe(EventScanResults)
	failed := cl.conn.Subscribe(EventScanFailed)

	for {
		select {
		case <-failed.Next():
			err = ErrScanFailed
			return
		case <-results.Next():
			break
		}
	}

	scanned, err := cl.conn.SendCommand(CmdScanResults)
	if err != nil {
		return
	}

	return parseAP([]byte(scanned))
}

// Networks lists the known networks
func (cl *Client) Networks() (nets []Network, err error) {
	data, err := cl.conn.SendCommand(CmdListNetworks)
	if err != nil {
		return nil, err
	}

	return parseNetwork([]byte(data))
}
