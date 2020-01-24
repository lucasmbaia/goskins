package steam

import (
	"regexp"
	"errors"
	"github.com/lucasmbaia/goskins/steam-api/request"
)

const (
	keyURL		= "https://steamcommunity.com/dev/apikey"
	accessDenied	= "<h2>Access Denied</h2>"
)

var (
	keyRegExp = regexp.MustCompile("<p>Key: ([0-9A-F]+)</p>")
)

func (s *Session) parseKey(resp request.Response) (key string, err error) {
	var (
		ok	bool
		subkey	[]string
	)

	if ok, err = regexp.Match(accessDenied, resp.Body); err != nil {
		return
	} else if ok {
		err = errors.New("Access Denied")
		return
	}

	subkey = keyRegExp.FindStringSubmatch(string(resp.Body))
	if len(subkey) != 2 {
		err = errors.New("Key not found")
		return
	}

	key = subkey[1]
	s.apiKey = key
	return
}

func (s *Session) GetWebApiKey() (key string, err error) {
	var (
		resp	request.Response
	)

	if resp, err = s.client.Request(request.GET, keyURL, request.Options{}); err != nil {
		return
	}

	return s.parseKey(resp)
}
