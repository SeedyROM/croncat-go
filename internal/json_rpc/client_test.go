package json_rpc

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

var ENDPOINT_URL = "https://uni-rpc.reece.sh"

func TestValidClientCall(t *testing.T) {
	client := &Client{Url: ENDPOINT_URL}

	resp, err := client.Call("block", nil, 1)

	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, resp)
	assert.Nil(t, resp.Error)
	assert.Equal(t, "2.0", resp.JsonRpc)
	assert.Equal(t, 1, resp.Id)
	assert.NotEmpty(t, resp.Result)
}

func TestInvalidClientCallMarhsal(t *testing.T) {
	client := &Client{Url: ENDPOINT_URL}

	// Test invalid request
	_, err := client.Call("block", []interface{}{1, 4, make(chan int)}, 1)

	assert.NotNil(t, err)
}

// TODO: Mock this.
func TestInvalidClientMethod(t *testing.T) {
	client := &Client{Url: ENDPOINT_URL}

	// Test invalid method
	_, err := client.Call("invalid", []interface{}{1, 4, "oops"}, 1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Method not found")
}

// TODO: Mock this.
func TestInvalidClientCallParams(t *testing.T) {
	client := &Client{Url: ENDPOINT_URL}

	_, err := client.Call("block", []interface{}{1, 4, "oops"}, 1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invalid params")

}

// TODO: Mock this.
func TestInvalidClientCallUrl(t *testing.T) {
	client := &Client{Url: "https://invalid.reece.sh"}

	_, err := client.Call("block", nil, 1)
	assert.NotNil(t, err)
}

func TestInvalidClientCallHttpResponse(t *testing.T) {
	defer gock.Off()

	// Mock the response from the server
	gock.New(ENDPOINT_URL).
		Post("/").
		Reply(500).
		BodyString("xjijijw")

	client := &Client{Url: ENDPOINT_URL}

	_, err := client.Call("block", nil, 1)
	assert.NotNil(t, err)
}
