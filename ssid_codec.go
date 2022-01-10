package wireless

import (
	"fmt"
	"strconv"
)

func DecodeSsid(ssid string) (string, error) {
	decoded := []byte{}
	for i := 0; i < len(ssid); i++ {
		switch ssid[i] {
		case '\\':
			i++
			switch ssid[i] {
			case '\\':
				decoded = append(decoded, '\\')
			case '"':
				decoded = append(decoded, '"')
			case 'n':
				decoded = append(decoded, '\n')
			case 'r':
				decoded = append(decoded, '\r')
			case 't':
				decoded = append(decoded, '\t')
			case 'e':
				decoded = append(decoded, '\033')
			case 'x':
				// Hexadecimal
				x, err := strconv.ParseInt(ssid[i+1:i+3], 16, 16)
				if err != nil {
					return string(decoded), err
				}
				decoded = append(decoded, byte(x))
				i++
				i++
			case '0', '1', '2', '3', '4', '5', '6', '7':
				// Octal
				x, err := strconv.ParseInt(ssid[i:i+3], 8, 24)
				if err != nil {
					return string(decoded), err
				}
				decoded = append(decoded, byte(x))
				i++
				i++
			default:
				return string(decoded), fmt.Errorf("unexpected escape sequence \\%c", ssid[i])
			}
		default:
			decoded = append(decoded, ssid[i])
		}
	}

	return string(decoded), nil
}

func EncodeSsid(ssid string) string {
	data := []byte(ssid)
	out := []byte{}
	for i := 0; i < len(data); i++ {
		switch data[i] {
		case '"':
			out = append(out, `\"`...)
		case '\\':
			out = append(out, `\\`...)
		case '\033':
			out = append(out, `\e`...)
		case '\n':
			out = append(out, `\n`...)
		case '\r':
			out = append(out, `\r`...)
		case '\t':
			out = append(out, `\t`...)
		default:
			if data[i] >= 32 && data[i] <= 127 {
				out = append(out, data[i])
			} else {
				out = append(out, fmt.Sprintf(`\x%02x`, data[i])...)
			}
		}
	}

	return string(out)
}
