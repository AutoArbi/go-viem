package client

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
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

	nonce, err := cl.getNonce(context.Background())
	if err != nil {
		t.Fatalf("getNonce failed: %v", err)
	}
	if nonce != 1 {
		t.Errorf("expected nonce 1, got %d", nonce)
	}

	chainID, err := cl.getChainID(context.Background())
	if err != nil {
		t.Fatalf("getChainID failed: %v", err)
	}
	expected := big.NewInt(1)
	if chainID.Cmp(expected) != 0 {
		t.Errorf("expected chainID %s, got %s", expected.String(), chainID.String())
	}
}

func TestSendETH_RequestError(t *testing.T) {
	mt := &mockTransport{
		requestFunc: func(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
			return nil, errors.New("simulated request error")
		},
	}
	// 这里使用一个无效私钥，仅用于触发逻辑，不会实际签名
	cl, err := NewClient(WithTransport(mt), WithPrivateKey("4c0883a69102937d6231471b5dbb6204fe512961708279a0d02c0b9a0d3d8a27"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	_, err = cl.SendETH(context.Background(), common.HexToAddress("0x456"), big.NewInt(1e18), 21000, big.NewInt(100e9), big.NewInt(2e9))
	if err == nil {
		t.Error("expected error from SendETH due to simulated request error")
	}
}
