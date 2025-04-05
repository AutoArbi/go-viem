package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type mockTransport struct {
	requestFunc func(ctx context.Context, method string, params ...any) (json.RawMessage, error)
}

func (m *mockTransport) Request(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	if m.requestFunc != nil {
		return m.requestFunc(ctx, method, params...)
	}
	return nil, errors.New("not implemented")
}

func TestNewClient_NoTransport(t *testing.T) {
	_, err := NewClient()
	if err == nil {
		t.Fatal("expected error when no transport provided")
	}
}

func TestNewClient_WithTransport(t *testing.T) {
	mt := &mockTransport{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			return json.RawMessage("\"OK\""), nil
		},
	}
	cl, err := NewClient(WithTransport(mt))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if cl.timeout != defaultTimeout {
		t.Errorf("expected default timeout %v, got %v", defaultTimeout, cl.timeout)
	}
	if cl.retryCount != defaultRetryCount {
		t.Errorf("expected default retry count %d, got %d", defaultRetryCount, cl.retryCount)
	}
}

func TestRequest_SuccessAfterRetry(t *testing.T) {
	attempt := 0
	mt := &mockTransport{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			attempt++
			if attempt < 2 {
				return nil, errors.New("simulated failure")
			}
			return json.RawMessage("\"Success\""), nil
		},
	}
	cl, err := NewClient(WithTransport(mt), WithRetryCount(2), WithPollingInterval(10*time.Millisecond))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	res, err := cl.Request(context.Background(), "test_method")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	var str string
	if err := json.Unmarshal(res, &str); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}
	if str != "Success" {
		t.Errorf("expected 'Success', got '%s'", str)
	}
}

func TestGetNonceAndChainID(t *testing.T) {
	mt := &mockTransport{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			switch method {
			case "getNonce":
				return json.RawMessage("\"0x1\""), nil
			case "getChainId":
				return json.RawMessage("\"0x1\""), nil
			default:
				return nil, errors.New("unexpected method")
			}
		},
	}
	cl, err := NewClient(WithTransport(mt))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	cl.from = common.HexToAddress("0x123")
}

func TestSendETH_RequestError(t *testing.T) {
	mt := &mockTransport{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			switch method {
			case "eth_chainId":
				return json.Marshal("0x1")
			case "eth_getTransactionCount":
				return json.Marshal("0x0")
			case "eth_sendRawTransaction":
				return nil, errors.New("simulated send raw tx error")
			default:
				return nil, fmt.Errorf("unexpected method call: %s", method)
			}
		},
	}

	cl, err := NewClient(
		WithTransport(mt),
		WithPrivateKey("4c0883a69102937d6231471b5dbb6204fe512961708279a0d02c0b9a0d3d8a27"),
	)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	to := common.HexToAddress("0x123")
	amount := big.NewInt(1e18) // 1 ETH
	chainID := big.NewInt(1)
	gasLimit := uint64(21000)
	nonce := uint64(0)
	maxFee := big.NewInt(25e9)     // 25 Gwei
	maxPriority := big.NewInt(2e9) // 2 Gwei

	_, err = cl.SendETH(
		context.Background(),
		to,
		amount,
		chainID,
		gasLimit,
		nonce,
		maxFee,
		maxPriority,
	)

	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if !strings.Contains(err.Error(), "simulated send raw tx error") {
		t.Errorf("Unexpected error message: %v", err)
	}
}
