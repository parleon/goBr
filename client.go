package gobr

import (
	"net/http"
)


type IClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	runtime *Runtime
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.runtime.do(req)
}
