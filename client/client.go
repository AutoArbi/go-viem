package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/AutoArbi/go-viem/transport"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Client support fallback
type Client struct {
	transports []transport.Transport
	privateKey *ecdsa.PrivateKey
	from       common.Address
}

func NewFallbackClient(transports ...transport.Transport) *Client {
	return &Client{transports: transports}
}

func (c *Client) Request(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	var lastErr error
	for _, ts := range c.transports {
		res, err := ts.Request(ctx, method, params...)
		if err == nil {
			return res, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("all transports failed: %w", lastErr)
}

// ================= AccessList ===================

// CreateAccessList performs eth_createAccessList RPC call using the first working transport.
func (c *Client) CreateAccessList(ctx context.Context, tx map[string]any) (types.AccessList, error) {
	var result struct {
		AccessList types.AccessList `json:"accessList"`
	}
	res, err := c.Request(ctx, "eth_createAccessList", tx)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode access list result: %w", err)
	}
	return result.AccessList, nil
}
