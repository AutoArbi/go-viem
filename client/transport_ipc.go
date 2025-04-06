package client

import (
	"context"
	"encoding/json"
	"github.com/AutoArbi/go-viem/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// IPCTransport struct
type IPCTransport struct {
	client *rpc.Client
}

// NewIPCTransport create a new IPCTransport instance
func NewIPCTransport(path string) (*IPCTransport, error) {
	c, err := rpc.DialIPC(context.Background(), path)
	if err != nil {
		return nil, err
	}
	return &IPCTransport{client: c}, nil
}

// Request implements the Transport interface's Request method
func (t *IPCTransport) Request(ctx context.Context, method types.RPCMethod, params ...any) (json.RawMessage, error) {
	var result json.RawMessage
	err := t.client.CallContext(ctx, &result, string(method), params...)
	return result, err
}
