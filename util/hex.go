package util

import "math/big"

func ParseHexUint64(hexStr string) (uint64, error) {
	val := new(big.Int)
	val.SetString(hexStr[2:], 16)
	return val.Uint64(), nil
}
