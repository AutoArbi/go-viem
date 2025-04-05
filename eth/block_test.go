package eth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AutoArbi/go-viem/client"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var _ client.Transport = (*mockClient)(nil)

func TestGetBlockByNumber(t *testing.T) {
	expectedBlockNumber := big.NewInt(100)
	expectedParam := fmt.Sprintf("0x%x", expectedBlockNumber) // "0x64"
	fullTx := true

	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method != "getBlockByNumber" {
				return nil, fmt.Errorf("unexpected method: %s", method)
			}
			if len(params) != 2 {
				return nil, fmt.Errorf("expected 2 params, got %d", len(params))
			}
			if params[0] != expectedParam {
				return nil, fmt.Errorf("expected block number param %s, got %v", expectedParam, params[0])
			}
			if params[1] != fullTx {
				return nil, fmt.Errorf("expected fullTx %v, got %v", fullTx, params[1])
			}
			return json.RawMessage("{\"block\": \"data\"}"), nil
		},
	}
	pc := &Client{Client: mock}
	res, err := pc.GetBlockByNumber(context.Background(), expectedBlockNumber, true)
	if err != nil {
		t.Fatalf("GetBlockByNumber error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(res, &obj); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if obj["block"] != "data" {
		t.Errorf("expected block data 'data', got %v", obj["block"])
	}
}

func TestGetBlockByHash(t *testing.T) {
	expectedBlockHash := common.HexToHash("0xabc123")
	fullTx := false

	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method != "getBlockByHash" {
				return nil, fmt.Errorf("unexpected method: %s", method)
			}
			if len(params) != 2 {
				return nil, fmt.Errorf("expected 2 params, got %d", len(params))
			}
			if params[0] != expectedBlockHash.Hex() {
				return nil, fmt.Errorf("expected block hash param %s, got %v", expectedBlockHash.Hex(), params[0])
			}
			if params[1] != fullTx {
				return nil, fmt.Errorf("expected fullTx %v, got %v", fullTx, params[1])
			}
			return json.RawMessage("{\"block\": \"hash data\"}"), nil
		},
	}
	pc := &Client{Client: mock}
	res, err := pc.GetBlockByHash(context.Background(), expectedBlockHash, fullTx)
	if err != nil {
		t.Fatalf("GetBlockByHash error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(res, &obj); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if obj["block"] != "hash data" {
		t.Errorf("expected block data 'hash data', got %v", obj["block"])
	}
}

func TestGetBlockTransactionCountByNumber(t *testing.T) {
	expectedBlockNumber := big.NewInt(200)
	expectedParam := fmt.Sprintf("0x%x", expectedBlockNumber)
	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method != "getBlockTransactionCountByNumber" {
				return nil, fmt.Errorf("unexpected method: %s", method)
			}
			if len(params) != 1 {
				return nil, fmt.Errorf("expected 1 param, got %d", len(params))
			}
			if params[0] != expectedParam {
				return nil, fmt.Errorf("expected param %s, got %v", expectedParam, params[0])
			}
			return json.RawMessage("\"0x5\""), nil
		},
	}
	pc := &Client{Client: mock}
	count, err := pc.GetBlockTransactionCountByNumber(context.Background(), expectedBlockNumber)
	if err != nil {
		t.Fatalf("GetBlockTransactionCountByNumber error: %v", err)
	}
	if count != 5 {
		t.Errorf("expected transaction count 5, got %d", count)
	}
}

func TestGetBlockTransactionCountByHash(t *testing.T) {
	expectedBlockHash := common.HexToHash("0xdef456")
	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method != "getBlockTransactionCountByHash" {
				return nil, fmt.Errorf("unexpected method: %s", method)
			}
			if len(params) != 1 {
				return nil, fmt.Errorf("expected 1 param, got %d", len(params))
			}
			if params[0] != expectedBlockHash.Hex() {
				return nil, fmt.Errorf("expected param %s, got %v", expectedBlockHash.Hex(), params[0])
			}
			return json.RawMessage("\"0x3\""), nil
		},
	}
	pc := &Client{Client: mock}
	count, err := pc.GetBlockTransactionCountByHash(context.Background(), expectedBlockHash)
	if err != nil {
		t.Fatalf("GetBlockTransactionCountByHash error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected transaction count 3, got %d", count)
	}
}

func TestSimulateBlocks(t *testing.T) {
	blockCount := 10
	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method != "simulateBlocks" {
				return nil, fmt.Errorf("unexpected method: %s", method)
			}
			if len(params) != 1 {
				return nil, fmt.Errorf("expected 1 param, got %d", len(params))
			}
			if bc, ok := params[0].(int); !ok || bc != blockCount {
				return nil, fmt.Errorf("expected blockCount %d, got %v", blockCount, params[0])
			}
			return json.RawMessage("{\"simulated\": true}"), nil
		},
	}
	pc := &Client{Client: mock}
	res, err := pc.SimulateBlocks(context.Background(), blockCount)
	if err != nil {
		t.Fatalf("SimulateBlocks error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(res, &obj); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if simulated, ok := obj["simulated"].(bool); !ok || !simulated {
		t.Errorf("expected simulated true, got %v", obj["simulated"])
	}
}

func TestWatchBlockNumber(t *testing.T) {
	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method != "watchBlockNumber" {
				return nil, fmt.Errorf("unexpected method: %s", method)
			}
			return json.RawMessage("{\"watch\": \"blockNumber\"}"), nil
		},
	}
	pc := &Client{Client: mock}
	res, err := pc.WatchBlockNumber(context.Background())
	if err != nil {
		t.Fatalf("WatchBlockNumber error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(res, &obj); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if obj["watch"] != "blockNumber" {
		t.Errorf("expected watch value 'blockNumber', got %v", obj["watch"])
	}
}

func TestWatchBlocks(t *testing.T) {
	blockCount := 5
	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method != "watchBlocks" {
				return nil, fmt.Errorf("unexpected method: %s", method)
			}
			if len(params) != 1 {
				return nil, fmt.Errorf("expected 1 param, got %d", len(params))
			}
			if bc, ok := params[0].(int); !ok || bc != blockCount {
				return nil, fmt.Errorf("expected blockCount %d, got %v", blockCount, params[0])
			}
			return json.RawMessage("{\"watch\": \"blocks\"}"), nil
		},
	}
	pc := &Client{Client: mock}
	res, err := pc.WatchBlocks(context.Background(), blockCount)
	if err != nil {
		t.Fatalf("WatchBlocks error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(res, &obj); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if obj["watch"] != "blocks" {
		t.Errorf("expected watch value 'blocks', got %v", obj["watch"])
	}
}
