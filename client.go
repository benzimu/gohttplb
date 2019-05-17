package gohttplb

import (
	"errors"
	"net/http"
	"strings"

	"gohttplb/utils"
)

// Errors
var (
	ErrInvalidAddr = errors.New("Invalid addr")
)

var defaultLBClientConfig = &LBClientConfig{
	SchedPolicy: PolicyRandomRetry,
	Retry:       3,
}

// LBClientConfig for init LBClient config
type LBClientConfig struct {
	// SchedPolicy request schedule policy
	// Default PolicyRandomRetry
	SchedPolicy SchedPolicy
	// Retry request retry count if return err
	// PolicyRandom and PolicyOrder no retry
	// Default 3
	Retry int
}

// LBClient ...
type LBClient struct {
	*R
}

// NewLBClient ...
func NewLBClient(addr string, config ...*LBClientConfig) (*LBClient, error) {
	if addr == "" {
		return nil, ErrInvalidAddr
	}

	addrs := utils.AddSchemeSlice(utils.TrimStringSlice(strings.Split(addr, ",")))
	if len(addrs) == 0 {
		return nil, ErrInvalidAddr
	}

	var conf = defaultLBClientConfig
	if len(config) > 0 {
		conf = config[0]
	}

	r := NewR(addrs, conf)

	return &LBClient{r}, nil
}

// Get ...
func (lbc *LBClient) Get(method, url string, params map[string]string, header map[string]string) (resp *http.Response, err error) {
	return lbc.get(http.MethodGet, url, params, header)
}

// Post ...
func (lbc *LBClient) Post(method, url string, params map[string]string, header map[string]string, body []byte) (resp *http.Response, err error) {
	return lbc.post(http.MethodPost, url, params, header, body)
}

// Delete ...
func (lbc *LBClient) Delete(method, url string, params map[string]string, header map[string]string) (resp *http.Response, err error) {
	return lbc.delete(http.MethodDelete, url, params, header)
}

// Put ...
func (lbc *LBClient) Put(method, url string, params map[string]string, header map[string]string, body []byte) (resp *http.Response, err error) {
	return lbc.put(http.MethodPut, url, params, header, body)
}

// Patch ...
func (lbc *LBClient) Patch(method, url string, params map[string]string, header map[string]string, body []byte) (resp *http.Response, err error) {
	return lbc.patch(http.MethodPatch, url, params, header, body)
}
