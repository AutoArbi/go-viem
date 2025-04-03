package util

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// RevertReason decodes the revert reason string from eth_call result
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
