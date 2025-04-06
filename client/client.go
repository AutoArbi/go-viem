package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AutoArbi/go-viem/types"
	"github.com/AutoArbi/go-viem/util"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"time"
)

const (
	defaultTimeout         = 30 * time.Second
	defaultPollingInterval = 5 * time.Second
	defaultRetryCount      = 3
	minRetryCount          = 0
)

// Option config function type
type Option func(*config) error

// Transport is a JSON-RPC Client interface
type Transport interface {
	Request(ctx context.Context, method types.RPCMethod, params ...any) (json.RawMessage, error)
}

// Client is a JSON-RPC Client that supports fallback
type Client struct {
	transport       []Transport
	privateKey      *ecdsa.PrivateKey
	from            common.Address
	timeout         time.Duration
	pollingInterval time.Duration
	retryCount      int
}

type config struct {
	transport       []Transport
	privateKey      *ecdsa.PrivateKey
	from            common.Address
	timeout         time.Duration
	pollingInterval time.Duration
	retryCount      int
}

// NewClient creates a Client and applies all options
func NewClient(opts ...Option) (*Client, error) {
	cfg := &config{
		timeout:         defaultTimeout,
		pollingInterval: defaultPollingInterval,
		retryCount:      defaultRetryCount,
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("apply option failed: %w", err)
		}
	}

	if len(cfg.transport) == 0 {
		return nil, errors.New("at least one transport required")
	}
	if cfg.timeout <= 0 {
		return nil, errors.New("timeout must be positive")
	}
	if cfg.retryCount < minRetryCount {
		return nil, fmt.Errorf("retry count must be >= %d", minRetryCount)
	}

	return &Client{
		transport:       cfg.transport,
		privateKey:      cfg.privateKey,
		from:            cfg.from,
		timeout:         cfg.timeout,
		pollingInterval: cfg.pollingInterval,
		retryCount:      cfg.retryCount,
	}, nil
}

// WithTransport adds Transport
func WithTransport(t ...Transport) Option {
	return func(c *config) error {
		if len(t) == 0 {
			return errors.New("transport cannot be empty")
		}
		c.transport = t
		return nil
	}
}

// WithPrivateKey sets the private key
func WithPrivateKey(privateKeyHex string) Option {
	return func(c *config) error {
		pk, err := crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			return err
		}
		c.privateKey = pk
		c.from = crypto.PubkeyToAddress(pk.PublicKey)
		return nil
	}
}

// WithTimeout sets the request timeout
func WithTimeout(d time.Duration) Option {
	return func(c *config) error {
		if d <= 0 {
			return errors.New("timeout must be positive")
		}
		c.timeout = d
		return nil
	}
}

// WithPollingInterval sets the polling interval
func WithPollingInterval(d time.Duration) Option {
	return func(c *config) error {
		if d <= 0 {
			return errors.New("polling interval must be positive")
		}
		c.pollingInterval = d
		return nil
	}
}

// WithRetryCount sets the retry count
func WithRetryCount(count int) Option {
	return func(c *config) error {
		if count < minRetryCount {
			return fmt.Errorf("retry count must be >= %d", minRetryCount)
		}
		c.retryCount = count
		return nil
	}
}

// Request calls all Transports in sequence and returns the first successful result
func (c *Client) Request(ctx context.Context, method types.RPCMethod, params ...any) (json.RawMessage, error) {
	var (
		res     json.RawMessage
		lastErr error
	)

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	for attempt := 0; attempt <= c.retryCount; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			for _, t := range c.transport {
				res, lastErr = t.Request(ctx, method, params...)
				if lastErr == nil {
					return res, nil
				}
			}

			if attempt < c.retryCount {
				time.Sleep(c.pollingInterval)
			}
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.retryCount+1, lastErr)
}

// SendETH sends ETH
func (c *Client) SendETH(ctx context.Context, to common.Address, amount, chainID *big.Int, gasLimit, nonce uint64, maxFeePerGas, maxPriorityFeePerGas *big.Int) (common.Hash, error) {
	if c.privateKey == nil {
		return common.Hash{}, fmt.Errorf("private key is required for wallet operations")
	}

	tx := &ethTypes.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &to,
		Value:     amount,
		Gas:       gasLimit,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Data:      nil, // Optional: fill in ABI-encoded contract call
	}

	signedTx, err := ethTypes.SignNewTx(c.privateKey, ethTypes.NewLondonSigner(chainID), tx)
	if err != nil {
		return common.Hash{}, err
	}

	var txHash common.Hash
	err = c.SendRawTransaction(ctx, signedTx, &txHash)
	return txHash, err
}

// SignTypedData signs EIP-712 structured data
func (c *Client) SignTypedData(typedDataJSON string) ([]byte, error) {
	if c.privateKey == nil {
		return nil, errors.New("wallet not initialized")
	}
	hash, err := util.TypedDataHash(typedDataJSON)
	if err != nil {
		return nil, err
	}
	return crypto.Sign(hash.Bytes(), c.privateKey)
}

// SignMessage signs a message
func (c *Client) SignMessage(msg []byte) ([]byte, error) {
	if c.privateKey == nil {
		return nil, errors.New("wallet not initialized")
	}
	msgHash := crypto.Keccak256Hash(
		[]byte("\x19Ethereum Signed Message:\n" + fmt.Sprint(len(msg)) + string(msg)),
	)
	return crypto.Sign(msgHash.Bytes(), c.privateKey)
}

// SendETH1559 sends an EIP-1559 transaction
func (c *Client) SendETH1559(ctx context.Context, to common.Address,
	amount, maxFeePerGas, maxPriorityFeePerGas, chainID *big.Int, gasLimit, nonce uint64,
	accessList ethTypes.AccessList) (common.Hash, error) {

	if c.privateKey == nil {
		return common.Hash{}, fmt.Errorf("private key is required for wallet operations")
	}
	tx := &ethTypes.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      nonce,
		To:         &to,
		Value:      amount,
		GasTipCap:  maxPriorityFeePerGas,
		GasFeeCap:  maxFeePerGas,
		Gas:        gasLimit,
		AccessList: accessList,
	}
	signedTx, err := ethTypes.SignNewTx(c.privateKey, ethTypes.NewLondonSigner(chainID), tx)
	if err != nil {
		return common.Hash{}, err
	}
	var txHash common.Hash
	err = c.SendRawTransaction(ctx, signedTx, &txHash)
	return txHash, err
}

// SendRawTransaction sends a raw transaction
// method: eth_sendRawTransaction
func (c *Client) SendRawTransaction(ctx context.Context, tx *ethTypes.Transaction, out *common.Hash) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	res, err := c.Request(ctx, types.SendRawTransaction, fmt.Sprintf("0x%x", data))
	if err != nil {
		return err
	}
	return json.Unmarshal(res, out)
}
