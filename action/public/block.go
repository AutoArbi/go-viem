package public

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AutoArbi/go-viem/util"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// GetBlockByNumber get block information by block number
// method: getBlockByNumber
func (c *Client) GetBlockByNumber(ctx context.Context, blockNumber *big.Int, fullTx bool) (json.RawMessage, error) {
	blockNumHex := fmt.Sprintf("0x%x", blockNumber)
	res, err := c.client.Request(ctx, "getBlockByNumber", blockNumHex, fullTx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetBlockByHash get block information by block hash
// method: getBlockByHash
func (c *Client) GetBlockByHash(ctx context.Context, blockHash common.Hash, fullTx bool) (json.RawMessage, error) {
	res, err := c.client.Request(ctx, "getBlockByHash", blockHash.Hex(), fullTx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetBlockTransactionCountByNumber get transaction count by block number
// method: getBlockTransactionCountByNumber
func (c *Client) GetBlockTransactionCountByNumber(ctx context.Context, blockNumber *big.Int) (uint64, error) {
	blockNumHex := fmt.Sprintf("0x%x", blockNumber)
	res, err := c.client.Request(ctx, "getBlockTransactionCountByNumber", blockNumHex)
	if err != nil {
		return 0, err
	}
	var countHex string
	if err := json.Unmarshal(res, &countHex); err != nil {
		return 0, err
	}
	return util.ParseHexUint64(countHex)
}

// GetBlockTransactionCountByHash get transaction count by block hash
// method: getBlockTransactionCountByHash
func (c *Client) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) (uint64, error) {
	res, err := c.client.Request(ctx, "getBlockTransactionCountByHash", blockHash.Hex())
	if err != nil {
		return 0, err
	}
	var countHex string
	if err := json.Unmarshal(res, &countHex); err != nil {
		return 0, err
	}
	return util.ParseHexUint64(countHex)
}

// SimulateBlocks simulate blocks
// method: simulateBlocks
func (c *Client) SimulateBlocks(ctx context.Context, blockCount int) (json.RawMessage, error) {
	res, err := c.client.Request(ctx, "simulateBlocks", blockCount)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// WatchBlockNumber watch block number
// method: watchBlockNumber
func (c *Client) WatchBlockNumber(ctx context.Context) (json.RawMessage, error) {
	res, err := c.client.Request(ctx, "watchBlockNumber")
	if err != nil {
		return nil, err
	}
	return res, nil
}

// WatchBlocks watch blocks
// method: watchBlocks
func (c *Client) WatchBlocks(ctx context.Context, blockCount int) (json.RawMessage, error) {
	res, err := c.client.Request(ctx, "watchBlocks", blockCount)
	if err != nil {
		return nil, err
	}
	return res, nil
}
