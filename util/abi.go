package util

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

// BuildCalldata generates ABI-encoded calldata for a method and args
func BuildCalldata(abiJSON, method string, args ...interface{}) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}
	return parsedABI.Pack(method, args...)
}
