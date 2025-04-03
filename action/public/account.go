package public

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AutoArbi/go-viem/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
)

// GetBalance gets the balance of the specified address at a certain block, default blockTag is "latest"
func (c *PublicClient) GetBalance(ctx context.Context, address common.Address, blockTag string) (*big.Int, error) {
	if blockTag == "" {
		blockTag = "latest"
	}
	res, err := c.client.Request(ctx, "getBalance", address.Hex(), blockTag)
	if err != nil {
		return nil, err
	}
	var hexBalance string
	if err := json.Unmarshal(res, &hexBalance); err != nil {
		return nil, err
	}
	val := new(big.Int)
	val.SetString(strings.TrimPrefix(hexBalance, "0x"), 16)
	return val, nil
}

// GetTransactionCount gets the transaction count of the specified address
func (c *PublicClient) GetTransactionCount(ctx context.Context, address common.Address, blockTag string) (uint64, error) {
	if blockTag == "" {
		blockTag = "latest"
	}
	res, err := c.client.Request(ctx, "getTransactionCount", address.Hex(), blockTag)
	if err != nil {
		return 0, err
	}
	var hexNonce string
	if err := json.Unmarshal(res, &hexNonce); err != nil {
		return 0, err
	}
	return util.ParseHexUint64(hexNonce)
}

// CreateAccessList creates an access list for the specified address
func (c *PublicClient) CreateAccessList(ctx context.Context, tx map[string]any) (types.AccessList, error) {
	var result struct {
		AccessList types.AccessList `json:"accessList"`
	}
	res, err := c.client.Request(ctx, "createAccessList", tx)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode access list result: %w", err)
	}
	return result.AccessList, nil
}
