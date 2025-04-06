package eth

import (
	"context"
	"github.com/AutoArbi/go-viem/transfer"
	"github.com/AutoArbi/go-viem/types"
	"math/big"
)

// GetChainID gets the chain ID
// method: eth_getChainId
func (c *Client) GetChainID(ctx context.Context) (*big.Int, error) {
	res, err := c.Client.Request(ctx, types.GetChainID)
	if err != nil {
		return nil, err
	}
	hexID := transfer.NewRPCResponseTransfer()
	return hexID.TransferBigInt(res)
}
