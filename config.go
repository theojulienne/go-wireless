package wireless

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

// Decode a configuration value returned from wpa_supplicant.
// wpa_supplicant understands values in 3 formats:
// Raw string, wrapped in quotes:
//  "this is config"
// `printf`-encoded string, wrapped in quotes and starting with `P`:
//  P"config with special char: \t"
// Hex-encoded, without quotes:
//  656e636f646564
// https://github.com/digsrc/wpa_supplicant/blob/515eb37dd1df3f4a05fc2a31c265db6358301988/src/utils/common.c#L631-L685
func configParseString(s string) (string, error) {
	l := len(s)
	if l == 0 {
		return "", nil
	}
	if s[0] == '"' {
		if s[l-1] != '"' {
			return "", errors.New("Quoted string did not end with a quote character")
		}
		return s[1 : l-1], nil
	} else if l > 1 && s[0] == 'P' && s[1] == '"' {
		if s[l-1] != '"' {
			return "", errors.New("printf_format string did not end with a quote character")
		}
		return DecodeSsid(s[2 : l-1])
	} else {
		// If the config value does not start with a quote character, it should be hex-encoded
		v, err := hex.DecodeString(s)
		return string(v), err
	}
}

// Determine if a string contains a "special character", as defined by wpa_supplicant,
// such that it should be encoded as a hex-string rather than literal string.
// https://github.com/digsrc/wpa_supplicant/blob/515eb37dd1df3f4a05fc2a31c265db6358301988/src/utils/common.c#L688-L697
func shouldHexEncode(s string) bool {
	for _, c := range s {
		if c < 32 || c >= 127 {
			return true
		}
	}
	return false
}

// Encode a string for saving or sending to wpa_supplicant configuration
// https://github.com/digsrc/wpa_supplicant/blob/515eb37dd1df3f4a05fc2a31c265db6358301988/wpa_supplicant/config.c#L147-L156
func configWriteString(s string) string {
	if shouldHexEncode(s) {
		return hex.EncodeToString([]byte(s))
	} else {
		return quote(s)
	}
}
