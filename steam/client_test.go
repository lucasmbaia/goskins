package steam

import (
	"testing"
)

func Test_Connect(t *testing.T) {
	var (
		c	*Client
		err	error
		done	= make(chan struct{})
	)

	c, _ = NewClient()

	if err = c.Connect(); err != nil {
		t.Fatal(err)
	}

	<-done
}
