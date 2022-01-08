package wireless

import (
	"strconv"
	"testing"
)

func TestConfigRoundTripPlain(t *testing.T) {
	tt := "hello world"
	in := configWriteString(tt)
	out, err := configParseString(in)
	if err != nil {
		t.Errorf("unexpected error in configParseString, %s", err)
	}
	if out != tt {
		t.Errorf("configParseString returned %s (wanted %s)", strconv.QuoteToASCII(out), strconv.QuoteToASCII(tt))
	}
}

func TestConfigRoundTripBackslash(t *testing.T) {
	tt := `hello\world`
	in := configWriteString(tt)
	out, err := configParseString(in)
	if err != nil {
		t.Errorf("unexpected error in configParseString, %s", err)
	}
	if out != tt {
		t.Errorf("configParseString returned %s (wanted %s)", strconv.QuoteToASCII(out), strconv.QuoteToASCII(tt))
	}
}

func TestConfigRoundTripTab(t *testing.T) {
	tt := "hello\tworld"
	in := configWriteString(tt)
	out, err := configParseString(in)
	if err != nil {
		t.Errorf("unexpected error in configParseString, %s", err)
	}
	if out != tt {
		t.Errorf("configParseString returned %s (wanted %s)", strconv.QuoteToASCII(out), strconv.QuoteToASCII(tt))
	}
}

func TestConfigRoundTripUnicode(t *testing.T) {
	tt := "Ë \" ' \\ 😅"
	in := configWriteString(tt)
	out, err := configParseString(in)
	if err != nil {
		t.Errorf("unexpected error in configParseString, %s", err)
	}
	if out != tt {
		t.Errorf("configParseString returned %s (wanted %s)", strconv.QuoteToASCII(out), strconv.QuoteToASCII(tt))
	}
}
