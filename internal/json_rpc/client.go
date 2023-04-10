package json_rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Url string
}

func (c *Client) Call(method string, params interface{}, id int) (*Response, error) {
	req := NewRequest("2.0", method, params, id)

	// Serialize the request
	str, err := req.ToJson()
	if err != nil {
		return nil, err
	}

	// Send the request
	resp, err := http.Post(c.Url, "application/json", bytes.NewBuffer(str))
	if err != nil {
		return nil, err
	}

	// Get the body
	body, err := ioutil.ReadAll(resp.Body)
	// notest
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	res := &Response{}
	err = res.FromJson(body)
	if err != nil {
		return nil, err
	}

	// Check if the error field is not nil
	if res.Error != nil {
		return nil, fmt.Errorf("%s: %s", res.Error.Message, res.Error.Data)
	}

	return res, nil
}
