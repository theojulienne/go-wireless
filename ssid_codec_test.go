package wireless

import (
	"strconv"
	"testing"
)

var tests = []struct {
	in  string
	out string
}{
	{`hello world`, "hello world"},
	{`hello\\world`, `hello\world`},
	{`hello\134world`, `hello\world`},
	{`\xc3\x8b \" ' \\ \xf0\x9f\x98\x85`, "Ë \" ' \\ 😅"},
}

func TestDecodeSsid(t *testing.T) {
	for _, tt := range tests {
		out, err := DecodeSsid(tt.in)
		if err != nil {
			t.Errorf("unexpected error in DecodeSsid, %s", err)
		}
		if out != tt.out {
			t.Errorf("DecodeSsid returned %s (wanted %s)", strconv.QuoteToASCII(out), strconv.QuoteToASCII(tt.out))
		}
	}
}

func TestEncodeDecodeSsid(t *testing.T) {
	for _, tt := range tests {
		in := EncodeSsid(tt.out)
		out, err := DecodeSsid(in)
		if err != nil {
			t.Errorf("unexpected error in DecodeSsid, %s", err)
		}
		if out != tt.out {
			t.Errorf("DecodeSsid returned %s (wanted %s)", strconv.QuoteToASCII(out), strconv.QuoteToASCII(tt.out))
		}
	}
}
