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

	ErrSSIDNotFound  = errors.New("SSID not found")
	ErrAuthFailed    = errors.New("auth failed")
	ErrDisconnected  = errors.New("disconnected")
	ErrAssocRejected = errors.New("assocation rejected")
	ErrNoIdentifier  = errors.New("no id_str field found")
	ErrInvalidEvent  = errors.New("invalid event message")
)

// IsUseOfClosedNetworkConnectionError will return true if the error is about use of
// closed network connection
func IsUseOfClosedNetworkConnectionError(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
