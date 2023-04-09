package json_rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestMarshal(t *testing.T) {
	assert := assert.New(t)

	req := NewRequest("2.0", "block", []interface{}{"latest"}, 1)

	str, err := req.ToJson()

	assert.Nil(err)
	assert.Equal([]byte(`{"jsonrpc":"2.0","method":"block","params":["latest"],"id":1}`), str)
}

func TestRequestUnmarshal(t *testing.T) {
	assert := assert.New(t)

	req := &Request{}

	err := req.FromJson([]byte(`{"jsonrpc":"2.0","method":"block","params":["latest"],"id":1}`))

	assert.Nil(err)
	assert.Equal("2.0", req.JsonRpc)
	assert.Equal("block", req.Method)
	assert.Equal([]interface{}{"latest"}, req.Params)
	assert.Equal(1, req.Id)
}

func TestResponseUnmarshal(t *testing.T) {
	assert := assert.New(t)

	req := &Response{}

	err := req.FromJson([]byte(`{"jsonrpc":"2.0","result":"0x1b4","id":1}`))

	assert.Nil(err)
	assert.Equal("2.0", req.JsonRpc)
	assert.Equal("0x1b4", req.Result)
	assert.Equal(1, req.Id)
}
