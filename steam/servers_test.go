package steam

import (
	"testing"
	"fmt"
)

func Test_GetServersSteam(t *testing.T) {
	if s, err := GetServersSteam(); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(s)
	}
}
