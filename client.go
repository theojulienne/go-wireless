package wireless

// Client represents a wireless client
type Client struct {
	conn *Conn
}

// Scan will scan for networks and return the APs it finds
func (cl *Client) Scan() (nets []AP, err error) {
	err = cl.conn.SendCommmandBool(CmdScan)
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
