package wireless

import (
	"strconv"
	"testing"
)

var configTests = []string{
	"hello world",
	`hello\world`,
	"hello\tworld",
	"Ë \" ' \\ 😅",
}

func TestConfigRoundTrip(t *testing.T) {
	for _, tt := range configTests {
		in := configWriteString(tt)
		out, err := configParseString(in)
		if err != nil {
			t.Errorf("unexpected error in configParseString, %s", err)
		}
		if out != tt {
			t.Errorf("configParseString returned %s (wanted %s)", strconv.QuoteToASCII(out), strconv.QuoteToASCII(tt))
		}
	}
}
