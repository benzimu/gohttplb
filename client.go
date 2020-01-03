package gohttplb

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Errors
var (
	ErrInvalidAddr         = errors.New("invalid addr")
	ErrInvalidAddrWeighted = errors.New("invalid addr weighted")
)

// Default config
var (
	DefaultAddrsSeparator        = ","
	DefaultAddrWeightedSeparator = "@@"
	DefaultStrategy              = StrategyRoundRobin

	DefaultRetry         = 1
	DefaultTransport     = defaultTransport
	DefaultClientTimeout = 10 * time.Second
)

// Headers
var (
	HeaderContentType  = "Content-Type"
	DefaultContentType = "application/json; charset=utf-8"
	HeaderAccept       = "Accept"
	DefaultAccept      = "application/json"
)

func setDefaultConf(conf *LBConfig) {
	if conf.Retry == 0 {
		conf.Retry = DefaultRetry
	}
	if conf.Transport == nil {
		conf.Transport = DefaultTransport
	}
	if conf.ClientTimeout == 0 {
		conf.ClientTimeout = DefaultClientTimeout
	}
	if conf.Client == nil {
		conf.Client = &http.Client{
			Transport: conf.Transport,
			Timeout:   conf.ClientTimeout,
		}
	}
}

// LBConfig for init LBClient config
type LBConfig struct {
	// Strategy request schedule policy
	// Default StrategyRoundRobin
	Strategy LoadBalancingStrategy
	// Retry request retry if return err, will retry all servers if retry set 1
	// Most retries: len(servers) * Retry
	// Default 1
	Retry int
	// ResponseParser response parser
	// Will auto parse response if set, and must use JPGet, JPPost...
	ResponseParser ResponseParser
	// Client for http request
	// If not set, will new with Transport and ClientTimeout
	Client *http.Client
	// Transport for http client
	// Default DefaultTransport
	Transport *http.Transport
	// ClientTimeout for `http.Client.Timeout`
	// Default 10s
	ClientTimeout time.Duration
}

// LBClient ...
type LBClient struct {
	*R
}

// New ...
func New(addr string, config ...*LBConfig) (*LBClient, error) {
	if addr == "" {
		return nil, ErrInvalidAddr
	}

	// set default config
	conf := &LBConfig{}
	if len(config) > 0 {
		conf = config[0]
	}
	setDefaultConf(conf)

	// validation addr param
	addrs := AddSchemeSlice(RemoveDuplicateElement(TrimStringSlice(
		strings.Split(addr, DefaultAddrsSeparator))))
	if len(addrs) == 0 {
		return nil, ErrInvalidAddr
	}
	log.Println("addrs:", addrs)

	strategy, serverWeighteds, err := determineStrategy(addrs)
	if err != nil {
		return nil, err
	}
	conf.Strategy = strategy

	// new R
	r := newR(addrs, serverWeighteds, conf)

	return &LBClient{r}, nil
}

func determineStrategy(addrs []string) (LoadBalancingStrategy, []ServerItem, error) {
	if len(addrs) == 0 {
		return 0, nil, ErrInvalidAddr
	}

	if existWeightedAddr(addrs) {
		weightedAddrs := make([]ServerItem, 0)
		for _, addr := range addrs {
			serverItem := ServerItem{}
			if strings.Contains(addr, DefaultAddrWeightedSeparator) {
				addrWeighteds := strings.Split(addr, DefaultAddrWeightedSeparator)
				serverItem.Server = addrWeighteds[0]
				weighted, err := strconv.Atoi(addrWeighteds[1])
				if err != nil {
					return 0, nil, fmt.Errorf("Error:%s:[%s]", ErrInvalidAddrWeighted, addr)
				}
				serverItem.Weighted = weighted
			} else {
				serverItem.Server = addr
				serverItem.Weighted = 1
			}
			weightedAddrs = append(weightedAddrs, serverItem)
		}
		return StrategyWeightedRoundRobin, weightedAddrs, nil
	}

	return DefaultStrategy, nil, nil
}

func existWeightedAddr(addrs []string) bool {
	for _, addr := range addrs {
		if strings.Contains(addr, DefaultAddrWeightedSeparator) {
			return true
		}
	}
	return false
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
func (lbc *LBClient) Get(path string, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.get(http.MethodGet, path, params, headers)
}

// Post method request
func (lbc *LBClient) Post(path string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.post(http.MethodPost, path, params, headers, body)
}

// Delete method request
func (lbc *LBClient) Delete(path string, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.delete(http.MethodDelete, path, params, headers)
}

// Put method request
func (lbc *LBClient) Put(path string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.put(http.MethodPut, path, params, headers, body)
}

// Patch method request
func (lbc *LBClient) Patch(path string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.patch(http.MethodPatch, path, params, headers, body)
}

// PGet get method request and parse response
func (lbc *LBClient) PGet(path string, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.parseGet(http.MethodGet, path, params, headers)
}

// PPost post method request and parse response
func (lbc *LBClient) PPost(path string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.parsePost(http.MethodPost, path, params, headers, body)
}

// PDelete delete method request and parse response
func (lbc *LBClient) PDelete(path string, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.parseDelete(http.MethodDelete, path, params, headers)
}

// PPut put method request and parse response
func (lbc *LBClient) PPut(path string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.parsePut(http.MethodPut, path, params, headers, body)
}

// PPatch patch method request and parse response
func (lbc *LBClient) PPatch(path string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.parseParamsHeaders(paramsHeaders)
	return lbc.R.parsePatch(http.MethodPatch, path, params, headers, body)
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
func (lbc *LBClient) JGet(path string, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.get(http.MethodGet, path, params, headers)
}

// JPost json post method request
func (lbc *LBClient) JPost(path string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.post(http.MethodPost, path, params, headers, body)
}

// JDelete json delete method request
func (lbc *LBClient) JDelete(path string, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.delete(http.MethodDelete, path, params, headers)
}

// JPut json put method request
func (lbc *LBClient) JPut(path string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.put(http.MethodPut, path, params, headers, body)
}

// JPatch json patch method request
func (lbc *LBClient) JPatch(path string, body []byte, paramsHeaders ...map[string]string) (resp *http.Response, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.patch(http.MethodPatch, path, params, headers, body)
}

// JPGet json get method request and parse response
func (lbc *LBClient) JPGet(path string, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.parseGet(http.MethodGet, path, params, headers)
}

// JPPost json post method request and parse response
func (lbc *LBClient) JPPost(path string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.parsePost(http.MethodPost, path, params, headers, body)
}

// JPDelete json delete method request and parse response
func (lbc *LBClient) JPDelete(path string, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.parseDelete(http.MethodDelete, path, params, headers)
}

// JPPut json put method request and parse response
func (lbc *LBClient) JPPut(path string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.parsePut(http.MethodPut, path, params, headers, body)
}

// JPPatch json patch method request and parse response
func (lbc *LBClient) JPPatch(path string, body []byte, paramsHeaders ...map[string]string) (statusCode int, data []byte, err error) {
	params, headers := lbc.jsonParamsHeaders(paramsHeaders)
	return lbc.R.parsePatch(http.MethodPatch, path, params, headers, body)
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
