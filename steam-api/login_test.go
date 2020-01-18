package steam

import (
	"testing"
)

func Test_LoginRequest(t *testing.T) {
	var (
		s   *Session
		err error
	)

	if s, err = NewSession(); err != nil {
		panic(err)
	}

	if _, err = s.loginRequest("", ""); err != nil {
		t.Fatal(err)
	}
}
