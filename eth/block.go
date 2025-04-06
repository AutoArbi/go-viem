package eth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AutoArbi/go-viem/transfer"
	"github.com/AutoArbi/go-viem/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// GetBlockNumber get the latest block number
// method: eth_blockNumber
func (c *Client) GetBlockNumber(ctx context.Context) (*big.Int, error) {
	res, err := c.Client.Request(ctx, types.GetBlockNumber)
	if err != nil {
		return nil, err
	}

	blockNumber := transfer.NewRPCResponseTransfer()
	return blockNumber.TransferBigInt(res)
}

// GetBlockByNumber get block information by block number
// method: eth_getBlockByNumber
func (c *Client) GetBlockByNumber(ctx context.Context, blockNumber *big.Int, fullTx bool) (json.RawMessage, error) {
	blockNumHex := fmt.Sprintf("0x%x", blockNumber)
	res, err := c.Client.Request(ctx, types.GetBlockByNumber, blockNumHex, fullTx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetBlockByHash get block information by block hash
// method: eth_getBlockByHash
func (c *Client) GetBlockByHash(ctx context.Context, blockHash common.Hash, fullTx bool) (json.RawMessage, error) {
	res, err := c.Client.Request(ctx, types.GetBlockByHash, blockHash.Hex(), fullTx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetBlockTransactionCountByNumber get transaction count by block number
// method: eth_getBlockTransactionCountByNumber
func (c *Client) GetBlockTransactionCountByNumber(ctx context.Context, blockNumber *big.Int) (uint64, error) {
	blockNumHex := fmt.Sprintf("0x%x", blockNumber)
	res, err := c.Client.Request(ctx, types.GetBlockTransactionCountByNumber, blockNumHex)
	if err != nil {
		return 0, err
	}

	countHex := transfer.NewRPCResponseTransfer()
	return countHex.TransferUint64(res)
}

// GetBlockTransactionCountByHash get transaction count by block hash
// method: eth_getBlockTransactionCountByHash
func (c *Client) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) (uint64, error) {
	res, err := c.Client.Request(ctx, types.GetBlockTransactionCountByHash, blockHash.Hex())
	if err != nil {
		return 0, err
	}
	countHex := transfer.NewRPCResponseTransfer()
	return countHex.TransferUint64(res)
}

// SimulateBlocks simulate blocks
// method: eth_simulateBlocks
func (c *Client) SimulateBlocks(ctx context.Context, blockCount int) (json.RawMessage, error) {
	res, err := c.Client.Request(ctx, types.SimulateBlocks, blockCount)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// WatchBlockNumber watch block number
// method: eth_watchBlockNumber
func (c *Client) WatchBlockNumber(ctx context.Context) (json.RawMessage, error) {
	res, err := c.Client.Request(ctx, types.WatchBlockNumber)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// WatchBlocks watch blocks
// method: eth_watchBlocks
func (c *Client) WatchBlocks(ctx context.Context, blockCount int) (json.RawMessage, error) {
	res, err := c.Client.Request(ctx, types.WatchBlocks, blockCount)
	if err != nil {
		return nil, err
	}
	return res, nil
}
