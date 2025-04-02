package transport

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
)

type HTTPTransport struct {
	endpoint string
	client   *rpc.Client
}

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

func (t *HTTPTransport) Request(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	var result json.RawMessage
	err := t.client.CallContext(ctx, &result, method, params...)
	return result, err
}
