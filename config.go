package wireless

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

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

func isHex(s string) bool {
	for _, c := range s {
		if c < 32 || c >= 127 {
			return true
		}
	}
	return false
}

// https://github.com/digsrc/wpa_supplicant/blob/515eb37dd1df3f4a05fc2a31c265db6358301988/wpa_supplicant/config.c#L147-L156
func configWriteString(s string) string {
	if isHex(s) {
		return hex.EncodeToString([]byte(s))
	} else {
		return quote(s)
	}
}
