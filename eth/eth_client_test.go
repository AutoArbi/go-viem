package eth

import (
	"context"
	"encoding/json"
	"errors"
)

type mockClient struct {
	requestFunc func(ctx context.Context, method string, params ...any) (json.RawMessage, error)
}

func (m *mockClient) Request(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	if m.requestFunc != nil {
		return m.requestFunc(ctx, method, params...)
	}
	return nil, errors.New("not implemented")
}
