package gohttplb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ResponseParser custom response parser
type ResponseParser interface {
	Parse(*http.Response) (int, []byte, error)
}

// DefaultResponseParser for default data struct
type DefaultResponseParser struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Parse implemente interface `ResponseParser.Parse`
func (parser *DefaultResponseParser) Parse(resp *http.Response) (statusCode int, data []byte, err error) {
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bodyBytes, parser); err != nil {
		return
	}

	data, err = json.Marshal(parser.Data)
	if err != nil {
		return
	}

	if !parser.Success {
		err = fmt.Errorf("Unsuccess: %d-%s", parser.Code, parser.Message)
		return
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		err = fmt.Errorf("StatusCode not ok: %d-%s", parser.Code, parser.Message)
		return
	}

	return
}
