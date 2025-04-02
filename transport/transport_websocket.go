package transport

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
)

type WebSocketTransport struct {
	client *rpc.Client
}

func NewWebSocketTransport(endpoint string) (*WebSocketTransport, error) {
	c, err := rpc.DialWebsocket(context.Background(), endpoint, "")
	if err != nil {
		return nil, err
	}
	return &WebSocketTransport{client: c}, nil
}

func (t *WebSocketTransport) Request(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	var result json.RawMessage
	err := t.client.CallContext(ctx, &result, method, params...)
	return result, err
}
