package eth

import (
	"context"
	"fmt"
	"github.com/AutoArbi/go-viem/transfer"
	"github.com/AutoArbi/go-viem/types"
	"strings"
)

// EstimateGas estimates the gas
// method: eth_estimateGas
func (c *Client) EstimateGas(ctx context.Context, call map[string]any) (uint64, error) {
	res, err := c.Client.Request(ctx, types.EstimateGas, call)
	if err != nil {
		return 0, err
	}
	hexGas := transfer.NewRPCResponseTransfer()
	return hexGas.TransferUint64(res)
}

// SimulateCall simulates an eth_call
// method: eth_simulateCall
func (c *Client) SimulateCall(ctx context.Context, call map[string]any, blockTag types.BlockTag) (string, error) {
	if blockTag == "" {
		blockTag = types.LATEST
	}
	res, err := c.Client.Request(ctx, types.SimulateCall, call, string(blockTag))
	if err != nil {
		return "", err
	}
	hexResult := transfer.NewRPCResponseTransfer()
	result, err := hexResult.TransferString(res)
	if err != nil {
		return "", fmt.Errorf("failed to parse simulation result: %w", err)
	}

	if !strings.HasPrefix(result, "0x") {
		return "", fmt.Errorf("invalid hex format result: %s", result)
	}
	return result, nil
}
