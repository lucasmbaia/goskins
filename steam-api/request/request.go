package request

import (
	"net/http/cookiejar"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"bytes"
	"net/url"
)

const (
	GET	= "GET"
	POST	= "POST"
	PUT	= "PUT"
	PATCH	= "PATCH"
	DELETE	= "DELETE"
)

type Response struct {
	Header	http.Header
	Code	int
	Body	[]byte
}

type Options struct {
	Body	    []byte
	Headers	    map[string]string
}

type Client struct {
	client	*http.Client
}

func NewClient() (c *Client, err error) {
	var jar *cookiejar.Jar

	if jar, err = cookiejar.New(nil); err != nil {
		return
	}

	return &Client{
		client:	&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true, Renegotiation: tls.RenegotiateOnceAsClient},
			},
			Jar:	jar,
		},
	}, nil
}

func (c *Client) SetCookies(u *url.URL, cookies	[]*http.Cookie) {
	c.client.Jar.SetCookies(u, cookies)
}

func (c *Client) GetCookies(u *url.URL) []*http.Cookie {
	return c.client.Jar.Cookies(u)
}

func (c *Client) Request(method, url string, o Options) (r Response, err error) {
	var (
		req	*http.Request
		resp	*http.Response
		b	[]byte
		pb	io.Reader
	)

	if len(o.Body) > 0 {
		pb = bytes.NewReader(o.Body)
	}

	if req, err = http.NewRequest(method, url, pb); err != nil {
		return
	}
	req.Close = true

	if o.Headers != nil {
		for k, v := range o.Headers {
			req.Header.Set(k, v)
		}
	}

	if resp, err = c.client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	r = Response{Header: resp.Header, Code: resp.StatusCode, Body: b}
	return
}
