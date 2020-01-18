package request

import (
	"testing"
)

func Test_Request(t *testing.T) {
	var c = NewClient()

	if _, err := c.Request("GET", "https://www.google.com.br/", &Options{}); err != nil {
		t.Fatal(err)
	}
}
