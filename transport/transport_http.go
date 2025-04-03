package transport

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
)

// HTTPTransport struct
type HTTPTransport struct {
	endpoint string
	client   *rpc.Client
}

// NewHTTPTransport create a new HTTPTransport instance
func NewHTTPTransport(endpoint string) (*HTTPTransport, error) {
	c, err := rpc.DialHTTP(endpoint)
	if err != nil {
		return nil, err
	}
	return &HTTPTransport{
		endpoint: endpoint,
		client:   c,
	}, nil
}

// Request implements the Transport interface's Request method
func (t *HTTPTransport) Request(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	var result json.RawMessage
	err := t.client.CallContext(ctx, &result, method, params...)
	return result, err
}
