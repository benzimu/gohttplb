package gohttplb

import (
	"bytes"
	"gohttplb/strategy"
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
	MaxIdleConnsPerHost: 3,
	IdleConnTimeout:     90 * time.Second,
	TLSHandshakeTimeout: 5 * time.Second,
}

// R is Request struct
type R struct {
	servers         []string
	serverWeighteds []strategy.ServerItem
	scheduler       strategy.Scheduler
	*LBConfig
	*http.Client
}

// NewR new R
func NewR(servers []string, serverWeighteds []strategy.ServerItem, conf *LBConfig) *R {
	r := &R{
		servers:         servers,
		serverWeighteds: serverWeighteds,
		LBConfig:        conf,
		Client:          conf.client,
	}

	r.scheduler = strategy.NewScheduler(conf.Strategy, servers, serverWeighteds)
	return r
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

func (r *R) doRetry(rA *rArgs) (resp *http.Response, err error) {
	serversSize := len(r.servers)
	for i := 0; i < r.Retry*serversSize; i++ {
		if serversSize == 1 {
			rA.url = r.servers[0] + rA.path
		} else {
			server := r.scheduler.Make()
			rA.url = server + rA.path
		}
		// TODO: check response status code 5xx 4xx
		// TODO: failed server handle
		resp, err = r.do(rA)
		if err != nil {
			continue
		}
		return
	}
	return
}

type rArgs struct {
	url     string
	method  string
	path    string
	params  map[string]string
	headers map[string]string
	body    []byte
}

func (r *R) doRequest(method, path string, params, headers map[string]string, body []byte) (resp *http.Response, err error) {
	rA := &rArgs{
		method:  method,
		path:    path,
		params:  params,
		headers: headers,
		body:    body,
	}
	return r.doRetry(rA)
}

func (r *R) get(method, path string, params map[string]string, headers map[string]string) (resp *http.Response, err error) {
	return r.doRequest(method, path, params, headers, nil)
}

func (r *R) post(method, path string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	return r.doRequest(method, path, params, headers, body)
}

func (r *R) delete(method, path string, params map[string]string, headers map[string]string) (resp *http.Response, err error) {
	return r.doRequest(method, path, params, headers, nil)
}

func (r *R) put(method, path string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	return r.doRequest(method, path, params, headers, body)
}

func (r *R) patch(method, path string, params map[string]string, headers map[string]string, body []byte) (resp *http.Response, err error) {
	return r.doRequest(method, path, params, headers, body)
}

func (r *R) parseDo(method, path string, params map[string]string, headers map[string]string, body []byte) (statusCode int, data []byte, err error) {
	response, err := r.doRequest(method, path, params, headers, body)
	if err != nil {
		return response.StatusCode, nil, err
	}

	if r.ResponseParser == nil {
		r.ResponseParser = &DefaultResponseParser{}
	}
	return r.ResponseParser.Parse(response)
}

func (r *R) parseGet(method, path string, params map[string]string, headers map[string]string) (statusCode int, data []byte, err error) {
	return r.parseDo(method, path, params, headers, nil)
}

func (r *R) parsePost(method, path string, params map[string]string, headers map[string]string, body []byte) (statusCode int, data []byte, err error) {
	return r.parseDo(method, path, params, headers, body)
}

func (r *R) parseDelete(method, path string, params map[string]string, headers map[string]string) (statusCode int, data []byte, err error) {
	return r.parseDo(method, path, params, headers, nil)
}

func (r *R) parsePut(method, path string, params map[string]string, headers map[string]string, body []byte) (statusCode int, data []byte, err error) {
	return r.parseDo(method, path, params, headers, body)
}

func (r *R) parsePatch(method, path string, params map[string]string, headers map[string]string, body []byte) (statusCode int, data []byte, err error) {
	return r.parseDo(method, path, params, headers, body)
}
