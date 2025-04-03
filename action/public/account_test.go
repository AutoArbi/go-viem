package public

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetBalance(t *testing.T) {
	// "0xDE0B6B3A7640000" 表示 1 ETH (1e18 wei)
	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method == "getBalance" {
				return json.RawMessage("\"0xDE0B6B3A7640000\""), nil
			}
			return nil, errors.New("unexpected method: " + method)
		},
	}
	pc := &Client{client: mock}

	balance, err := pc.GetBalance(context.Background(), common.HexToAddress("0x123"), "latest")
	if err != nil {
		t.Fatalf("GetBalance error: %v", err)
	}

	expected := new(big.Int)
	expected.SetString("de0b6b3a7640000", 16) // 1e18 in hex
	if balance.Cmp(expected) != 0 {
		t.Errorf("expected balance %s, got %s", expected.String(), balance.String())
	}
}

func TestGetTransactionCount(t *testing.T) {
	// "0x10" 表示 nonce 为 16
	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method == "getTransactionCount" {
				return json.RawMessage("\"0x10\""), nil
			}
			return nil, errors.New("unexpected method: " + method)
		},
	}
	pc := &Client{client: mock}

	nonce, err := pc.GetTransactionCount(context.Background(), common.HexToAddress("0x123"), "pending")
	if err != nil {
		t.Fatalf("GetTransactionCount error: %v", err)
	}
	if nonce != 16 {
		t.Errorf("expected nonce 16, got %d", nonce)
	}
}

func TestCreateAccessList(t *testing.T) {
	mockResponse := `{"accessList": [{"address": "0x0000000000000000000000000000000000000001", "storageKeys": []}]}`
	mock := &mockClient{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			if method == "createAccessList" {
				return json.RawMessage(mockResponse), nil
			}
			return nil, errors.New("unexpected method: " + method)
		},
	}
	pc := &Client{client: mock}

	txParams := map[string]any{
		"from": "0x123",
		"to":   "0x456",
		"data": "0x",
	}
	al, err := pc.CreateAccessList(context.Background(), txParams)
	if err != nil {
		t.Fatalf("CreateAccessList error: %v", err)
	}

	if len(al) != 1 {
		t.Fatalf("expected access list length 1, got %d", len(al))
	}
	expectedAddr := common.HexToAddress("0x0000000000000000000000000000000000000001")
	if !reflect.DeepEqual(al[0].Address, expectedAddr) {
		t.Errorf("expected address %s, got %s", expectedAddr.Hex(), al[0].Address.Hex())
	}
	if len(al[0].StorageKeys) != 0 {
		t.Errorf("expected empty storageKeys, got %v", al[0].StorageKeys)
	}
}
