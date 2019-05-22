package gohttplb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ResponseParser custom response parser
type ResponseParser interface {
	parse(*http.Response) (int, []byte, error)
}

// ResponseD define a ResponseD struct
type ResponseD struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// DefaultResponseParser for default data struct
type DefaultResponseParser struct {
	*ResponseD
}

func (parser *DefaultResponseParser) parse(resp *http.Response) (statusCode int, data []byte, err error) {
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bodyBytes, parser.Data); err != nil {
		return
	}

	data, err = json.Marshal(parser.ResponseD.Data)
	if err != nil {
		return
	}

	if !parser.ResponseD.Success {
		err = fmt.Errorf("Unsuccess: %s", parser.ResponseD.Message)
		return
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		err = fmt.Errorf("StatusCode not ok: %s", parser.ResponseD.Message)
		return
	}

	return
}
