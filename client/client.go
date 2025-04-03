package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AutoArbi/go-viem/transport"
	"github.com/AutoArbi/go-viem/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

// Client is a JSON-RPC client interface
type Client interface {
	Request(ctx context.Context, method string, params ...any) (json.RawMessage, error)
}

// client is a JSON-RPC client that supports fallback
type client struct {
	transports      []transport.Transport
	privateKey      *ecdsa.PrivateKey
	from            common.Address
	timeout         time.Duration
	pollingInterval time.Duration
	retryCount      int
}

type config struct {
	transports      []transport.Transport
	privateKey      *ecdsa.PrivateKey
	from            common.Address
	timeout         time.Duration
	pollingInterval time.Duration
	retryCount      int
}

// NewClient creates a client and applies all options
func NewClient(opts ...Option) (*client, error) {
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

	if len(cfg.transports) == 0 {
		return nil, errors.New("at least one transport required")
	}
	if cfg.timeout <= 0 {
		return nil, errors.New("timeout must be positive")
	}
	if cfg.retryCount < minRetryCount {
		return nil, fmt.Errorf("retry count must be >= %d", minRetryCount)
	}

	return &client{
		transports:      cfg.transports,
		privateKey:      cfg.privateKey,
		from:            cfg.from,
		timeout:         cfg.timeout,
		pollingInterval: cfg.pollingInterval,
		retryCount:      cfg.retryCount,
	}, nil
}

// WithTransport adds Transport
func WithTransport(t ...transport.Transport) Option {
	return func(c *config) error {
		if len(t) == 0 {
			return errors.New("transports cannot be empty")
		}
		c.transports = t
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
func (c *client) Request(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
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
			for _, t := range c.transports {
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
func (c *client) SendETH(ctx context.Context, to common.Address, amount *big.Int, gasLimit uint64,
	maxFeePerGas *big.Int, maxPriorityFeePerGas *big.Int) (common.Hash, error) {
	if c.privateKey == nil {
		return common.Hash{}, fmt.Errorf("private key is required for wallet operations")
	}

	nonce, err := c.getNonce(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	chainID, err := c.getChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &to,
		Value:     amount,
		Gas:       gasLimit,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Data:      nil, // Optional: fill in ABI-encoded contract call
	}

	signedTx, err := types.SignNewTx(c.privateKey, types.NewLondonSigner(chainID), tx)
	if err != nil {
		return common.Hash{}, err
	}

	var txHash common.Hash
	err = c.sendRawTransaction(ctx, signedTx, &txHash)
	return txHash, err
}

// SendETH1559 sends an EIP-1559 transaction
func (c *client) SendETH1559(ctx context.Context, to common.Address,
	amount, maxFeePerGas, maxPriorityFeePerGas *big.Int, gasLimit uint64,
	accessList types.AccessList) (common.Hash, error) {

	if c.privateKey == nil {
		return common.Hash{}, fmt.Errorf("private key is required for wallet operations")
	}
	nonce, err := c.getNonce(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	chainID, err := c.getChainID(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	tx := &types.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      nonce,
		To:         &to,
		Value:      amount,
		GasTipCap:  maxPriorityFeePerGas,
		GasFeeCap:  maxFeePerGas,
		Gas:        gasLimit,
		AccessList: accessList,
	}
	signedTx, err := types.SignNewTx(c.privateKey, types.NewLondonSigner(chainID), tx)
	if err != nil {
		return common.Hash{}, err
	}
	var txHash common.Hash
	err = c.sendRawTransaction(ctx, signedTx, &txHash)
	return txHash, err
}

// SignTypedData signs EIP-712 structured data
func (c *client) SignTypedData(typedDataJSON string) ([]byte, error) {
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
func (c *client) SignMessage(msg []byte) ([]byte, error) {
	if c.privateKey == nil {
		return nil, errors.New("wallet not initialized")
	}
	msgHash := crypto.Keccak256Hash(
		[]byte("\x19Ethereum Signed Message:\n" + fmt.Sprint(len(msg)) + string(msg)),
	)
	return crypto.Sign(msgHash.Bytes(), c.privateKey)
}

// SimulateCall simulates an eth_call
func (c *client) SimulateCall(ctx context.Context, call map[string]any, blockTag string) (string, error) {
	if blockTag == "" {
		blockTag = "latest"
	}
	res, err := c.Request(ctx, "simulateCall", call, blockTag)
	if err != nil {
		return "", err
	}
	var hexResult string
	if err := json.Unmarshal(res, &hexResult); err != nil {
		return "", err
	}
	return hexResult, nil
}

// EstimateGas estimates the gas
func (c *client) EstimateGas(ctx context.Context, call map[string]any) (uint64, error) {
	res, err := c.Request(ctx, "estimateGas", call)
	if err != nil {
		return 0, err
	}
	var hexGas string
	if err := json.Unmarshal(res, &hexGas); err != nil {
		return 0, err
	}
	return util.ParseHexUint64(hexGas)
}

func (c *client) getNonce(ctx context.Context) (uint64, error) {
	res, err := c.Request(ctx, "getNonce", c.from.Hex(), "pending")
	if err != nil {
		return 0, err
	}
	var hexNonce string
	if err := json.Unmarshal(res, &hexNonce); err != nil {
		return 0, err
	}
	return util.ParseHexUint64(hexNonce)
}

func (c *client) getChainID(ctx context.Context) (*big.Int, error) {
	res, err := c.Request(ctx, "getChainId")
	if err != nil {
		return nil, err
	}
	var hexID string
	if err := json.Unmarshal(res, &hexID); err != nil {
		return nil, err
	}
	id := new(big.Int)
	id.SetString(hexID[2:], 16)
	return id, nil
}

func (c *client) sendRawTransaction(ctx context.Context, tx *types.Transaction, out *common.Hash) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	res, err := c.Request(ctx, "sendRawTransaction", fmt.Sprintf("0x%x", data))
	if err != nil {
		return err
	}
	return json.Unmarshal(res, out)
}
