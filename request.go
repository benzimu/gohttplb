package gohttplb

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"
)

var defaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	},
	Timeout: 10 * time.Second,
}

// R is Request struct
type R struct {
	servers []string
	*LBClientConfig
	*http.Client
}

// NewR new R
func NewR(servers []string, conf *LBClientConfig) *R {
	return &R{
		servers:        servers,
		LBClientConfig: conf,
		Client:         defaultHTTPClient,
	}
}

func (r *R) do(rargs *rArgs) (resp *http.Response, err error) {
	var paramStr string
	if len(rargs.params) != 0 {
		paramStr += "?"
		for key, val := range rargs.params {
			paramStr += key + "=" + val
		}
	}
	rargs.url += paramStr

	var bodyReader io.Reader
	if rargs.body != nil {
		bodyReader = bytes.NewReader(rargs.body)
	}

	req, err := http.NewRequest(rargs.method, rargs.url, bodyReader)
	if err != nil {
		return nil, err
	}

	for key, val := range rargs.headers {
		req.Header.Set(key, val)
	}

	resp, err = r.Do(req)
	return
}

type rArgs struct {
	method  string
	url     string
	params  map[string]string
	headers map[string]string
	body    []byte
}

func (r *R) get(method, url string, params map[string]string, headers map[string]string) (resp *http.Response, err error) {
	rA := &rArgs{
		method:  method,
		url:     url,
		params:  params,
		headers: headers,
	}
	rScheduler := NewRScheduler(r)
	return rScheduler.schedule(rA)
}

func (r *R) post(method, url string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	rA := &rArgs{
		method:  method,
		url:     url,
		params:  params,
		headers: headers,
		body:    body,
	}
	rScheduler := NewRScheduler(r)
	return rScheduler.schedule(rA)
}

func (r *R) delete(method, url string, params map[string]string, headers map[string]string) (resp *http.Response, err error) {
	rA := &rArgs{
		method:  method,
		url:     url,
		params:  params,
		headers: headers,
	}
	rScheduler := NewRScheduler(r)
	return rScheduler.schedule(rA)
}

func (r *R) put(method, url string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	rA := &rArgs{
		method:  method,
		url:     url,
		params:  params,
		headers: headers,
		body:    body,
	}
	rScheduler := NewRScheduler(r)
	return rScheduler.schedule(rA)
}

func (r *R) patch(method, url string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	rA := &rArgs{
		method:  method,
		url:     url,
		params:  params,
		headers: headers,
		body:    body,
	}
	rScheduler := NewRScheduler(r)
	return rScheduler.schedule(rA)
}
