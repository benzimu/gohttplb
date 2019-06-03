package gohttplb

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

// Errors
var (
	ErrInvalidAddr = errors.New("Invalid addr")
)

// Headers
var (
	HeaderContentType  = "Content-Type"
	DefaultContentType = "application/json; charset=utf-8"
	HeaderAccept       = "Accept"
	DefaultAccept      = "application/json"
)

// Default config
var (
	DefaultSchedPolicy   = PolicyRandom
	DefaultRetry         = 1
	DefaultSeparator     = ","
	DefaultTransport     = defaultTransport
	DefaultClientTimeout = 10 * time.Second
)

func setDefaultConf(conf *LBConfig) {
	if conf.SchedPolicy == 0 {
		conf.SchedPolicy = DefaultSchedPolicy
	}
	if conf.Retry == 0 {
		conf.Retry = DefaultRetry
	}
	if conf.Separator == "" {
		conf.Separator = DefaultSeparator
	}
	if conf.Transport == nil {
		conf.Transport = DefaultTransport
	}
	if conf.ClientTimeout == 0 {
		conf.ClientTimeout = DefaultClientTimeout
	}
	conf.Client = &http.Client{
		Transport: conf.Transport,
		Timeout:   conf.ClientTimeout,
	}
}

// LBConfig for init LBClient config
type LBConfig struct {
	// SchedPolicy request schedule policy
	// Default PolicyRandom
	SchedPolicy SchedPolicy
	// Retry request retry if return err, will retry all servers if retry set 1
	// Most retries: len(servers) * Retry
	// Default 1
	Retry int
	// Client for http request
	// Default defaultHTTPClient
	Client *http.Client
	// Separator split addr
	// Default ","
	Separator string
	// ResponseParser response parser
	// Will auto parse response if set, and must use ParseGet, ParsePost...
	ResponseParser ResponseParser
	// Transport for http client
	Transport *http.Transport
	// ClientTimeout for `http.Client.Timeout`
	ClientTimeout time.Duration
}

// LBClient ...
type LBClient struct {
	r *R
}

// New ...
func New(addr string, config ...*LBConfig) (*LBClient, error) {
	if addr == "" {
		return nil, ErrInvalidAddr
	}

	conf := &LBConfig{}
	if len(config) > 0 {
		conf = config[0]
	}
	setDefaultConf(conf)

	addrs := AddSchemeSlice(RemoveDuplicateElement(TrimStringSlice(strings.Split(addr, conf.Separator))))
	if len(addrs) == 0 {
		return nil, ErrInvalidAddr
	}

	r := NewR(addrs, conf)

	return &LBClient{r}, nil
}

func (lbc *LBClient) parseParamsHeaders(paramsHeaders []map[string]string) (params map[string]string, headers map[string]string) {
	if len(paramsHeaders) == 0 {
		return
	} else if len(paramsHeaders) == 1 {
		params = paramsHeaders[0]
		return
	}

	params = paramsHeaders[0]
	headers = paramsHeaders[1]
	return
}

// Get method request
func (lbc *LBClient) Get(url string, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.get(http.MethodGet, url, params, headers)
}

// Post method request
func (lbc *LBClient) Post(url string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.post(http.MethodPost, url, params, headers, body)
}

// Delete method request
func (lbc *LBClient) Delete(url string, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.delete(http.MethodDelete, url, params, headers)
}

// Put method request
func (lbc *LBClient) Put(url string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.put(http.MethodPut, url, params, headers, body)
}

// Patch method request
func (lbc *LBClient) Patch(url string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.patch(http.MethodPatch, url, params, headers, body)
}

// PGet get method request and parse response
func (lbc *LBClient) PGet(url string, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.parseGet(http.MethodGet, url, params, headers)
}

// PPost post method request and parse response
func (lbc *LBClient) PPost(url string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.parsePost(http.MethodPost, url, params, headers, body)
}

// PDelete delete method request and parse response
func (lbc *LBClient) PDelete(url string, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.parseDelete(http.MethodDelete, url, params, headers)
}

// PPut put method request and parse response
func (lbc *LBClient) PPut(url string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.parsePut(http.MethodPut, url, params, headers, body)
}

// PPatch patch method request and parse response
func (lbc *LBClient) PPatch(url string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.r.parsePatch(http.MethodPatch, url, params, headers, body)
}

func (lbc *LBClient) setJSONHeader(headers map[string]string) {
	headers[HeaderContentType] = DefaultContentType
	headers[HeaderAccept] = DefaultAccept
}

func (lbc *LBClient) jsonParamsHeaders(paramsHeaders []map[string]string) (params map[string]string, headers map[string]string) {
	params, headers = lbc.parseParamsHeaders(paramsHeaders)
	if headers == nil {
		headers = make(map[string]string)
	}
	lbc.setJSONHeader(headers)
	return
}

// JGet json get method request
func (lbc *LBClient) JGet(url string, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.get(http.MethodGet, url, params, headers)
}

// JPost json post method request
func (lbc *LBClient) JPost(url string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.post(http.MethodPost, url, params, headers, body)
}

// JDelete json delete method request
func (lbc *LBClient) JDelete(url string, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.delete(http.MethodDelete, url, params, headers)
}

// JPut json put method request
func (lbc *LBClient) JPut(url string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.put(http.MethodPut, url, params, headers, body)
}

// JPatch json patch method request
func (lbc *LBClient) JPatch(url string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.patch(http.MethodPatch, url, params, headers, body)
}

// JPGet json get method request and parse response
func (lbc *LBClient) JPGet(url string, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.parseGet(http.MethodGet, url, params, headers)
}

// JPPost json post method request and parse response
func (lbc *LBClient) JPPost(url string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.parsePost(http.MethodPost, url, params, headers, body)
}

// JPDelete json delete method request and parse response
func (lbc *LBClient) JPDelete(url string, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.parseDelete(http.MethodDelete, url, params, headers)
}

// JPPut json put method request and parse response
func (lbc *LBClient) JPPut(url string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.parsePut(http.MethodPut, url, params, headers, body)
}

// JPPatch json patch method request and parse response
func (lbc *LBClient) JPPatch(url string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.r.parsePatch(http.MethodPatch, url, params, headers, body)
}

// PResponse parse response use custom or default ResponseParser
func (lbc *LBClient) PResponse(resp *http.Response, rps ...ResponseParser) (int, []byte, error) {
	var rp ResponseParser
	if len(rps) > 0 {
		rp = rps[0]
	} else {
		rp = &DefaultResponseParser{}
	}

	return rp.Parse(resp)
}
