package transport

import (
	"context"
	"encoding/json"
)

// Transport JSON-RPC transport interface
type Transport interface {
	Request(ctx context.Context, method string, params ...any) (json.RawMessage, error)
}

// Client JSON-RPC client
type Client struct {
	transports []Transport // 支持 fallback
}
