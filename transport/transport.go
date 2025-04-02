package transport

import (
	"context"
	"encoding/json"
)

type Transport interface {
	Request(ctx context.Context, method string, params ...any) (json.RawMessage, error)
}

type Client struct {
	transports []Transport // 支持 fallback
}
