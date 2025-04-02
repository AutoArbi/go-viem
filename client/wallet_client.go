package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AutoArbi/go-viem/transport"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strings"
)

func NewWalletClient(privateKeyHex string, transports ...transport.Transport) (*Client, error) {
	pk, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	from := crypto.PubkeyToAddress(pk.PublicKey)
	return &Client{
		transports: transports,
		privateKey: pk,
		from:       from,
	}, nil
}

// ============== Wallet Client ================

func (c *Client) SendETH(ctx context.Context, to common.Address, amount *big.Int, gasLimit uint64,
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
		Data:      nil, // 可选：填入 ABI 编码后的合约调用
	}

	signedTx, err := types.SignNewTx(c.privateKey, types.NewLondonSigner(chainID), tx)
	if err != nil {
		return common.Hash{}, err
	}

	var txHash common.Hash
	err = c.sendRawTransaction(ctx, signedTx, &txHash)
	return txHash, err
}

func (c *Client) SendETH1559(ctx context.Context, to common.Address, amount, maxFeePerGas, maxPriorityFeePerGas *big.Int, gasLimit uint64, accessList types.AccessList) (common.Hash, error) {
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

// SignTypedData signs EIP-712 structured data (typedDataJSON is JSON string per EIP-712)
func (c *Client) SignTypedData(typedDataJSON string) ([]byte, error) {
	if c.privateKey == nil {
		return nil, errors.New("wallet not initialized")
	}
	hash, err := TypedDataHash(typedDataJSON)
	if err != nil {
		return nil, err
	}
	return crypto.Sign(hash.Bytes(), c.privateKey)
}

// TypedDataHash implements EIP-712 hash(domain separator + message struct hash)
func TypedDataHash(typedDataJSON string) (common.Hash, error) {
	var typedData apitypes.TypedData

	if err := json.Unmarshal([]byte(typedDataJSON), &typedData); err != nil {
		return common.Hash{}, fmt.Errorf("invalid typed data json: %w", err)
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to hash domain: %w", err)
	}

	messageHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to hash message: %w", err)
	}

	// Final EIP-712 hash = keccak256(\"\\x19\\x01\" || domainSeparator || messageHash)
	var digestBytes []byte
	digestBytes = append(digestBytes, []byte{0x19, 0x01}...)
	digestBytes = append(digestBytes, domainSeparator...)
	digestBytes = append(digestBytes, messageHash...)

	final := sha3.NewLegacyKeccak256()
	final.Write(digestBytes)
	return common.BytesToHash(final.Sum(nil)), nil
}

// RevertReason decodes revert reason string from eth_call result
func RevertReason(hexData string) (string, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(hexData, "0x"))
	if err != nil {
		return "", err
	}
	if len(b) < 4 || !bytes.HasPrefix(b, []byte{0x08, 0xc3, 0x79, 0xa0}) {
		return "", fmt.Errorf("not revert reason format")
	}
	// offset 4 is selector, next 32 bytes is data offset, next 32 bytes is string length, then data
	if len(b) < 4+32+32 {
		return "", fmt.Errorf("invalid revert reason")
	}
	strlen := new(big.Int).SetBytes(b[4+32 : 4+64]).Int64()
	if int64(len(b)) < 4+64+strlen {
		return "", fmt.Errorf("invalid revert reason length")
	}
	return string(b[4+64 : 4+64+strlen]), nil
}

// BuildCalldata generates ABI-encoded calldata for a method and args
func BuildCalldata(abiJSON, method string, args ...interface{}) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}
	return parsedABI.Pack(method, args...)
}

func (c *Client) SignMessage(msg []byte) ([]byte, error) {
	if c.privateKey == nil {
		return nil, errors.New("wallet not initialized")
	}
	msgHash := crypto.Keccak256Hash(
		[]byte("\x19Ethereum Signed Message:\n" + fmt.Sprint(len(msg)) + string(msg)),
	)
	return crypto.Sign(msgHash.Bytes(), c.privateKey)
}

func (c *Client) SimulateCall(ctx context.Context, call map[string]any, blockTag string) (string, error) {
	if blockTag == "" {
		blockTag = "latest"
	}
	res, err := c.Request(ctx, "eth_call", call, blockTag)
	if err != nil {
		return "", err
	}
	var hexResult string
	if err := json.Unmarshal(res, &hexResult); err != nil {
		return "", err
	}
	return hexResult, nil
}

func (c *Client) EstimateGas(ctx context.Context, call map[string]any) (uint64, error) {
	res, err := c.Request(ctx, "eth_estimateGas", call)
	if err != nil {
		return 0, err
	}
	var hexGas string
	if err := json.Unmarshal(res, &hexGas); err != nil {
		return 0, err
	}
	return parseHexUint64(hexGas)
}

func (c *Client) getNonce(ctx context.Context) (uint64, error) {
	res, err := c.Request(ctx, "eth_getTransactionCount", c.from.Hex(), "pending")
	if err != nil {
		return 0, err
	}
	var hexNonce string
	if err := json.Unmarshal(res, &hexNonce); err != nil {
		return 0, err
	}
	return parseHexUint64(hexNonce)
}

func (c *Client) getChainID(ctx context.Context) (*big.Int, error) {
	res, err := c.Request(ctx, "eth_chainId")
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

func (c *Client) sendRawTransaction(ctx context.Context, tx *types.Transaction, out *common.Hash) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	res, err := c.Request(ctx, "eth_sendRawTransaction", fmt.Sprintf("0x%x", data))
	if err != nil {
		return err
	}
	return json.Unmarshal(res, out)
}

func parseHexUint64(hexStr string) (uint64, error) {
	val := new(big.Int)
	val.SetString(hexStr[2:], 16)
	return val.Uint64(), nil
}
