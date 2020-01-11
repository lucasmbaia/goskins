package request

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"bytes"
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
	Transport   http.RoundTripper
	Username    string
	Password    string

	postBody    io.Reader
}

func Request(method, url string, o *Options) (r Response, err error) {
	var (
		client	http.Client
		req	*http.Request
		resp	*http.Response
		b	[]byte
	)

	if req, err = http.NewRequest(method, url, o.postBody); err != nil {
		return
	}

	client = http.Client{Transport: o.Transport}
	req.Close = true

	if o.Headers != nil {
		for k, v := range o.Headers {
			req.Header.Set(k, v)
		}
	}

	if o.Username != "" && o.Password != "" {
		req.SetBasicAuth(o.Username, o.Password)
	}

	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	r = Response{Header: resp.Header, Code: resp.StatusCode, Body: b}
	return
}

func SetOptions(o *Options) {
	o.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, Renegotiation: tls.RenegotiateOnceAsClient},
	}

	if len(o.Body) > 0 {
		o.postBody = bytes.NewReader(o.Body)
	}

	return
}
