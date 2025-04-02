package transport

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
)

type IPCTransport struct {
	client *rpc.Client
}

func NewIPCTransport(path string) (*IPCTransport, error) {
	c, err := rpc.DialIPC(context.Background(), path)
	if err != nil {
		return nil, err
	}
	return &IPCTransport{client: c}, nil
}

func (t *IPCTransport) Request(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	var result json.RawMessage
	err := t.client.CallContext(ctx, &result, method, params...)
	return result, err
}
