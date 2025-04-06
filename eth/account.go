package eth

import (
	"context"
	"fmt"
	"github.com/AutoArbi/go-viem/transfer"
	"github.com/AutoArbi/go-viem/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

// GetBalance gets the balance of the specified address at a certain block, default blockTag is "latest"
// method: eth_getBalance
func (c *Client) GetBalance(ctx context.Context, address common.Address, blockTag types.BlockTag) (*big.Int, error) {
	if blockTag == "" {
		blockTag = types.LATEST
	}
	res, err := c.Client.Request(ctx, types.GetBalance, address.Hex(), blockTag)
	if err != nil {
		return nil, err
	}
	balance := transfer.NewRPCResponseTransfer()
	return balance.TransferBigInt(res)
}

// GetTransactionCount gets the transaction count of the specified address
// method: eth_getTransactionCount
func (c *Client) GetTransactionCount(ctx context.Context, address common.Address, blockTag types.BlockTag) (uint64, error) {
	if blockTag == "" {
		blockTag = types.LATEST
	}
	res, err := c.Client.Request(ctx, types.GetTransactionCount, address.Hex(), blockTag)
	if err != nil {
		return 0, err
	}
	transactionCount := transfer.NewRPCResponseTransfer()
	return transactionCount.TransferUint64(res)
}

// CreateAccessList creates an access list for the specified address
// method: eth_createAccessList
func (c *Client) CreateAccessList(ctx context.Context, tx map[string]any) (ethTypes.AccessList, error) {

	var response struct {
		AccessList ethTypes.AccessList `json:"accessList"`
		Error      string              `json:"error,omitempty"`
		GasUsed    string              `json:"gasUsed"`
	}

	res, err := c.Client.Request(ctx, types.CreateAccessList, tx)
	if err != nil {
		return nil, err
	}
	accessList := transfer.NewRPCResponseTransfer()
	if err := accessList.TransferStruct(res, &response); err != nil {
		return nil, fmt.Errorf("failed to parse access list: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("RPC error: %s", response.Error)
	}

	return response.AccessList, nil
}
