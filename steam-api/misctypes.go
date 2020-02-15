package steam

import (
	"encoding/json"
	"strconv"
	"time"
)

type seconds struct {
	time.Duration
}

func (s *seconds) UnmarshalJSON(b []byte) error {
	i, err := unmarshalStringyInt(b)
	if err != nil {
		return err
	}
	s.Duration = time.Duration(i) * time.Second
	return nil
}

type timestamp struct {
	time.Time
}

func (t *timestamp) UnmarshalJSON(b []byte) error {
	s, err := unmarshalStringyValue(b)
	if err != nil {
		return err
	}
	t.Time, err = time.Parse(time.RFC3339, s)
	if _, ok := err.(*time.ParseError); ok {
		// ok, that didn't work lets try parsing an int
		var i int64
		if i, err = strconv.ParseInt(s, 10, 64); err != nil {
			return err
		}
		t.Time = time.Unix(i, 0)
	}
	return err
}

func (t *timestamp) String() string {
	return strconv.FormatInt(t.Unix(), 10)
}

func unmarshalStringyValue(b []byte) (s string, err error) {
	err = json.Unmarshal(b, &s)
	return
}

// see previous statement x2
func unmarshalStringyInt(b []byte) (i int64, err error) {
	s := ""
	s, err = unmarshalStringyValue(b)
	if err != nil {
		return
	}
	i, err = strconv.ParseInt(s, 10, 64)
	return
}

