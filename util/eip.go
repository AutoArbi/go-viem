package util

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"golang.org/x/crypto/sha3"
)

// TypedDataHash implements EIP-712 hash (domain separator + message struct hash)
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

	// Final EIP-712 hash = keccak256("\x19\x01" || domainSeparator || messageHash)
	var digestBytes []byte
	digestBytes = append(digestBytes, []byte{0x19, 0x01}...)
	digestBytes = append(digestBytes, domainSeparator...)
	digestBytes = append(digestBytes, messageHash...)

	final := sha3.NewLegacyKeccak256()
	final.Write(digestBytes)
	return common.BytesToHash(final.Sum(nil)), nil
}
