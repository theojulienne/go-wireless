package wireless

import (
	"context"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

// WPAConn is an interface to the connection
type WPAConn interface {
	SendCommand(...string) (string, error)
	SendCommandBool(...string) error
	SendCommandInt(...string) (int, error)
	SendCommandWithContext(context.Context, ...string) (string, error)
	SendCommandBoolWithContext(context.Context, ...string) error
	SendCommandIntWithContext(context.Context, ...string) (int, error)
	io.Closer
	Subscribe(...string) *Subscription
}

// Client represents a wireless client
type Client struct {
	conn        WPAConn
	ScanTimeout time.Duration
	CmdTimeout  time.Duration
	ctx         context.Context
}

// NewClient will create a new client by connecting to the
// given interface in WPA
func NewClient(iface string) (c *Client, err error) {
	c = new(Client)
	c.conn, err = Dial(iface)
	if err != nil {
		return
	}
	c.CmdTimeout = time.Second
	return
}

// NewClientFromConn returns a new client from an already established connection
func NewClientFromConn(conn WPAConn) (c *Client) {
	c = new(Client)
	c.conn = conn
	return
}

func (cl *Client) getContext() (context.Context, func()) {
	if cl.ctx != nil {
		return cl.ctx, func() {}
	}

	return context.WithTimeout(context.Background(), cl.CmdTimeout)
}

func (cl *Client) WithContext(ctx context.Context, ops ...func(wc *Client)) {
	if len(ops) == 0 {
		cl.ctx = ctx
		return
	}

	wc := NewClientFromConn(cl.conn)
	wc.ctx = ctx
	for _, op := range ops {
		op(wc)
	}
}

// Close will close the client connection
func (cl *Client) Close() {
	err := cl.conn.Close()
	if err != nil {
		log.Println("ERROR: client failed to close conn:", err)
	}
}

// Conn will return the underlying connection
func (cl *Client) Conn() *Conn {
	return cl.conn.(*Conn)
}

// Subscribe will subscribe to certain events that happen in WPA
func (cl *Client) Subscribe(topics ...string) *Subscription {
	return cl.conn.Subscribe(topics...)
}

// Status will return the current state of the WPA
func (cl *Client) Status() (State, error) {
	ctx, cancel := cl.getContext()
	defer cancel()
	data, err := cl.conn.SendCommandWithContext(ctx, CmdStatus)
	if err != nil {
		return State{}, err
	}
	s := NewState(data)
	return s, nil
}

// Scan will scan for networks and return the APs it finds
func (cl *Client) Scan() (nets APs, err error) {
	timeout := cl.ScanTimeout

	if timeout == 0 {
		timeout = 2 * time.Second
	}

	err = cl.conn.SendCommandBool(CmdScan)
	if err != nil {
		return
	}

	results := cl.conn.Subscribe(EventScanResults)
	failed := cl.conn.Subscribe(EventScanFailed)

	defer results.Unsubscribe()
	defer failed.Unsubscribe()

	func() {
		for {
			select {
			case <-failed.Next():
				err = ErrScanFailed
				return
			case <-results.Next():
				return
			case <-time.NewTimer(cl.ScanTimeout).C:
				return
			}
		}
	}()

	scanned, err := cl.conn.SendCommand(CmdScanResults)
	if err != nil {
		return
	}

	return parseAP([]byte(scanned))
}

// Networks lists the known networks
func (cl *Client) Networks() (nets Networks, err error) {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.NetworksWithContext(ctx)
}

func (cl *Client) NetworksWithContext(ctx context.Context) (nets Networks, err error) {
	data, err := cl.conn.SendCommandWithContext(ctx, CmdListNetworks)
	if err != nil {
		return nil, err
	}

	nets, err = parseNetwork([]byte(data))
	if err != nil {
		return nil, err
	}

	for i := range nets {
		nets[i].Known = true
		(&nets[i]).populateAttrs(cl)
	}

	return nets, nil
}

func (cl *Client) Connect(net Network) (Network, error) {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.ConnectWithContext(ctx, net)
}

// Connect to a new or existing network
func (cl *Client) ConnectWithContext(ctx context.Context, net Network) (Network, error) {
	net, err := cl.AddOrUpdateNetworkWithContext(ctx, net)
	if err != nil {
		return net, err
	}

	sub := cl.conn.Subscribe(EventNetworkNotFound, EventAuthReject, EventConnected, EventDisconnected, EventAssocReject)
	defer sub.Unsubscribe()

	if err := cl.EnableNetworkWithContext(ctx, net.ID); err != nil {
		return net, err
	}

	if err := cl.conn.SendCommandBoolWithContext(ctx, "REASSOCIATE"); err != nil {
		return net, err
	}

	if err := cl.conn.SendCommandBoolWithContext(ctx, "RECONNECT"); err != nil {
		return net, err
	}

	select {
	case ev := <-sub.Next():
		switch ev.Name {
		case EventConnected:
			return net, cl.SaveConfig()
		case EventNetworkNotFound:
			return net, ErrSSIDNotFound
		case EventAuthReject:
			return net, ErrAuthFailed
		case EventDisconnected:
			return net, ErrDisconnected
		case EventAssocReject:
			return net, ErrAssocRejected
		default:
			return net, errors.New("failed to catch event " + ev.Name)
		}
	case <-ctx.Done():
		return net, ctx.Err()
	}
}

// AddOrUpdateNetwork will add or, if the network has IDStr set, update it
func (cl *Client) AddOrUpdateNetwork(net Network) (Network, error) {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.AddNetworkWithContext(ctx, net)
}

func (cl *Client) AddOrUpdateNetworkWithContext(ctx context.Context, net Network) (Network, error) {
	if net.IDStr != "" {
		nets, err := cl.NetworksWithContext(ctx)
		if err != nil {
			return net, err
		}

		for _, n := range nets {
			if n.IDStr == net.IDStr {
				return cl.UpdateNetworkWithContext(ctx, net)
			}
		}
	}

	return cl.AddNetworkWithContext(ctx, net)
}

// UpdateNetwork will update the given network, an error will be thrown
// if the network doesn't have IDStr specified
func (cl *Client) UpdateNetwork(net Network) (Network, error) {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.UpdateNetworkWithContext(ctx, net)
}

func (cl *Client) UpdateNetworkWithContext(ctx context.Context, net Network) (Network, error) {
	if net.IDStr == "" {
		return net, ErrNoIdentifier
	}

	for _, cmd := range setCmds(net) {
		if err := cl.conn.SendCommandBoolWithContext(ctx, cmd); err != nil {
			return net, err
		}
	}

	return net, nil
}

// AddNetwork will add a new network
func (cl *Client) AddNetwork(net Network) (Network, error) {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.AddNetworkWithContext(ctx, net)
}

// AddNetwork will add a new network
func (cl *Client) AddNetworkWithContext(ctx context.Context, net Network) (Network, error) {
	nets, err := cl.Networks()
	if err == nil {
		if nw, found := nets.FindBySSID(net.SSID); found {
			return nw, nil
		}
	}

	i, err := cl.conn.SendCommandIntWithContext(ctx, CmdAddNetwork)
	if err != nil {
		return net, err
	}

	net.ID = i

	if net.IDStr == "" {
		net.IDStr = net.SSID
	}

	for _, cmd := range setCmds(net) {
		if err := cl.conn.SendCommandBoolWithContext(ctx, cmd); err != nil {
			return net, err
		}
	}

	net.Known = true
	return net, nil
}

// RemoveNetwork will RemoveNetwork
func (cl *Client) RemoveNetwork(id int) error {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.conn.SendCommandBoolWithContext(ctx, CmdRemoveNetwork, strconv.Itoa(id))
}

// EnableNetwork will EnableNetwork
func (cl *Client) EnableNetwork(id int) error {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.EnableNetworkWithContext(ctx, id)
}

func (cl *Client) EnableNetworkWithContext(ctx context.Context, id int) error {
	return cl.conn.SendCommandBoolWithContext(ctx, CmdEnableNetwork, strconv.Itoa(id))
}

// SelectNetwork will SelectNetwork
func (cl *Client) SelectNetwork(id int) error {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.SelectNetworkWithContext(ctx, id)
}

func (cl *Client) SelectNetworkWithContext(ctx context.Context, id int) error {
	return cl.conn.SendCommandBoolWithContext(ctx, CmdSelectNetwork, strconv.Itoa(id))
}

// Disconnect will Disconnect
func (cl *Client) Disconnect() error {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.DisconnectWithContext(ctx)
}

func (cl *Client) DisconnectWithContext(ctx context.Context) error {
	return cl.conn.SendCommandBoolWithContext(ctx, CmdDisconnect)
}

// DisableNetwork will DisableNetwork
func (cl *Client) DisableNetwork(id int) error {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.conn.SendCommandBoolWithContext(ctx, CmdDisableNetwork+" "+strconv.Itoa(id))
}

// SaveConfig will SaveConfig
func (cl *Client) SaveConfig() error {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.conn.SendCommandBoolWithContext(ctx, CmdSaveConfig)
}

// LoadConfig will LoadConfig
func (cl *Client) LoadConfig() error {
	ctx, cancel := cl.getContext()
	defer cancel()
	return cl.conn.SendCommandBoolWithContext(ctx, CmdReconfigure)
}

// GetNetworkAttr will get the given attribute of the given network
func (cl *Client) GetNetworkAttr(id int, attr string) (string, error) {
	ctx, cancel := cl.getContext()
	defer cancel()
	s, err := cl.conn.SendCommandWithContext(ctx, CmdGetNetwork, strconv.Itoa(id), attr)
	if err != nil {
		return s, err
	}

	return strings.TrimSpace(s), nil
}
