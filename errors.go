package wireless

import (
	"errors"
	"strings"
)

var (

	// ErrCmdTimeout is an error that happens when the command times out
	ErrCmdTimeout = errors.New("timeout while waiting for command response")

	// ErrScanFailed is an error that happens when scanning for wifi networks fails
	ErrScanFailed = errors.New("scan failed")

	ErrSSIDNotFound = errors.New("SSID not found")

	// ErrAuthFailed this is not actually talking about WPA but the raw 802.11 management frames.  This
	// was originally called WEP but when security holes were found it became disused, instead passing
	// to higher layer protocols like 802.1X (EAP) or 802.11i (WPA2).  You should never see this error.
	ErrAuthFailed = errors.New("raw 802.11 auth failed")

	ErrDisconnected = errors.New("disconnected")

	// ErrAssocRejected will be given should the station fail to associate with the AP, usually due
	// to higher level authentication protocols like WPA2 failing to authenticate a password.  You can
	// use this to detect invalid passwords.
	ErrAssocRejected = errors.New("assocation rejected")

	ErrNoIdentifier = errors.New("no id_str field found")
	ErrInvalidEvent = errors.New("invalid event message")
	ErrFailBusy     = errors.New("FAIL-BUSY")
)

// IsUseOfClosedNetworkConnectionError will return true if the error is about use of
// closed network connection
func IsUseOfClosedNetworkConnectionError(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
