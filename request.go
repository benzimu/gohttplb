package gohttplb

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"time"
)

var defaultTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:        100,
	IdleConnTimeout:     90 * time.Second,
	TLSHandshakeTimeout: 5 * time.Second,
}

// R is Request struct
type R struct {
	servers []string
	*LBConfig
	*http.Client
}

// NewR new R
func NewR(servers []string, conf *LBConfig) *R {
	return &R{
		servers:  servers,
		LBConfig: conf,
		Client:   conf.Client,
	}
}

func (r *R) do(rargs *rArgs) (resp *http.Response, err error) {
	var bodyReader io.Reader
	if rargs.body != nil {
		bodyReader = bytes.NewReader(rargs.body)
	}

	req, err := http.NewRequest(rargs.method, rargs.url, bodyReader)
	if err != nil {
		return nil, err
	}

	if len(rargs.params) != 0 {
		q := req.URL.Query()
		for key, val := range rargs.params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
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

func (r *R) doSchedule(method, url string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	rA := &rArgs{
		method:  method,
		url:     url,
		params:  params,
		headers: headers,
		body:    body,
	}
	// TODO: filter last request failed server
	if len(r.servers) == 1 {
		rA.url = r.servers[0] + rA.url
		return r.do(rA)
	}
	rScheduler := NewRScheduler(r)
	return rScheduler.schedule(rA)
}

func (r *R) get(method, url string, params map[string]string, headers map[string]string) (resp *http.Response, err error) {
	return r.doSchedule(method, url, params, headers, nil)
}

func (r *R) post(method, url string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	return r.doSchedule(method, url, params, headers, body)
}

func (r *R) delete(method, url string, params map[string]string, headers map[string]string) (resp *http.Response, err error) {
	return r.doSchedule(method, url, params, headers, nil)
}

func (r *R) put(method, url string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	return r.doSchedule(method, url, params, headers, body)
}

func (r *R) patch(method, url string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	return r.doSchedule(method, url, params, headers, body)
}

func (r *R) parseDo(method, url string, params map[string]string, headers map[string]string, body []byte) (statusCode int, data []byte, err error) {
	response, err := r.doSchedule(method, url, params, headers, body)
	if err != nil {
		return response.StatusCode, nil, err
	}

	if r.ResponseParser == nil {
		r.ResponseParser = &DefaultResponseParser{}
	}
	return r.ResponseParser.Parse(response)
}

func (r *R) parseGet(method, url string, params map[string]string, headers map[string]string) (statusCode int, data []byte, err error) {
	return r.parseDo(method, url, params, headers, nil)
}

func (r *R) parsePost(method, url string, params map[string]string, headers map[string]string, body []byte) (statusCode int, data []byte, err error) {
	return r.parseDo(method, url, params, headers, body)
}

func (r *R) parseDelete(method, url string, params map[string]string, headers map[string]string) (statusCode int, data []byte, err error) {
	return r.parseDo(method, url, params, headers, nil)
}

func (r *R) parsePut(method, url string, params map[string]string, headers map[string]string, body []byte) (statusCode int, data []byte, err error) {
	return r.parseDo(method, url, params, headers, body)
}

func (r *R) parsePatch(method, url string, params map[string]string, headers map[string]string, body []byte) (statusCode int, data []byte, err error) {
	return r.parseDo(method, url, params, headers, body)
}
