package request

import (
	"testing"
)

func Test_Request(t *testing.T) {
	if _, err := Request("GET", "https://www.google.com.br/", &Options{}); err != nil {
		t.Fatal(err)
	}
}
