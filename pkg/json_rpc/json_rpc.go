package json_rpc

import "encoding/json"

type Request struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

func NewRequest(jsonrpc string, method string, params interface{}, id int) *Request {
	return &Request{
		JsonRpc: jsonrpc,
		Method:  method,
		Params:  params,
		Id:      id,
	}
}

func (r *Request) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Request) FromJson(data []byte) error {
	return json.Unmarshal(data, r)
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type Response struct {
	JsonRpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error"`
	Id      int         `json:"id"`
}

func (r *Response) FromJson(data []byte) error {
	return json.Unmarshal(data, r)
}
