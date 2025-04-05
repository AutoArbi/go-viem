package transfer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"reflect"
	"strings"
)

type RPCResponseTransfer struct {
}

func NewRPCResponseTransfer() *RPCResponseTransfer {
	return &RPCResponseTransfer{}
}

// TransferBigInt parses a JSON-RPC response containing a hex string into a big.Int.
func (p *RPCResponseTransfer) TransferBigInt(response json.RawMessage) (*big.Int, error) {
	var hexStr string
	if err := json.Unmarshal(response, &hexStr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	val := new(big.Int)
	_, ok := val.SetString(strings.TrimPrefix(hexStr, "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return val, nil
}

// TransferUint64 parses a JSON-RPC response containing a hex string into an uint64.
func (p *RPCResponseTransfer) TransferUint64(response json.RawMessage) (uint64, error) {
	var hexStr string
	if err := json.Unmarshal(response, &hexStr); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	val := new(big.Int)
	_, ok := val.SetString(strings.TrimPrefix(hexStr, "0x"), 16)
	if !ok {
		return 0, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return val.Uint64(), nil
}

// TransferString parses a JSON-RPC response containing a string.
// TransferString 解析包含字符串类型的JSON-RPC响应
func (p *RPCResponseTransfer) TransferString(response json.RawMessage) (string, error) {
	if len(response) == 0 {
		return "", errors.New("empty response data")
	}

	if bytes.Equal(response, []byte("null")) {
		return "", errors.New("received null response")
	}

	var result string
	if err := json.Unmarshal(response, &result); err != nil {
		return "", fmt.Errorf("failed to parse string response [raw: %s]: %w",
			string(response),
			err)
	}

	if strings.HasPrefix(result, "0x") && len(result) < 3 {
		return "", fmt.Errorf("invalid hex string: %s", result)
	}

	return result, nil
}

// TransferBool parses a JSON-RPC response containing a boolean value.
func (p *RPCResponseTransfer) TransferBool(response json.RawMessage) (bool, error) {
	var result bool
	if err := json.Unmarshal(response, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal bool response: %w", err)
	}
	return result, nil
}

// TransferAddress parses a JSON-RPC response containing a hex string into an Ethereum address.
func (p *RPCResponseTransfer) TransferAddress(response json.RawMessage) (common.Address, error) {
	var addressHex string
	if err := json.Unmarshal(response, &addressHex); err != nil {
		return common.Address{}, fmt.Errorf("failed to unmarshal address response: %w", err)
	}

	address := common.HexToAddress(addressHex)
	return address, nil
}

// TransferTransactionReceipt parses a JSON-RPC response into a transaction.
func (p *RPCResponseTransfer) TransferTransactionReceipt(response json.RawMessage) (*types.Receipt, error) {
	var receipt types.Receipt
	if err := json.Unmarshal(response, &receipt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction receipt response: %w", err)
	}
	return &receipt, nil
}

// TransferBlock parses a JSON-RPC response into a block.
func (p *RPCResponseTransfer) TransferBlock(response json.RawMessage) (*types.Block, error) {
	var block types.Block
	if err := json.Unmarshal(response, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block response: %w", err)
	}
	return &block, nil
}

// TransferStruct parses a JSON-RPC response into a specified struct.
func (p *RPCResponseTransfer) TransferStruct(response json.RawMessage, result interface{}) error {
	if err := json.Unmarshal(response, result); err == nil {
		return nil
	}

	return p.transferWithCustomTypes(response, result)
}

// transferWithCustomTypes parses a JSON-RPC response into a struct with custom types.
func (p *RPCResponseTransfer) transferWithCustomTypes(data json.RawMessage, out interface{}) error {
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return fmt.Errorf("failed to unmarshal raw response: %w", err)
	}

	targetValue := reflect.ValueOf(out).Elem()
	targetType := targetValue.Type()

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldValue := targetValue.Field(i)

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		tagParts := strings.Split(tag, ",")
		jsonKey := tagParts[0]

		rawValue, exists := rawMap[jsonKey]
		if !exists {
			continue
		}

		// 根据字段类型处理
		switch fieldValue.Interface().(type) {
		case common.Hash:
			var s string
			if err := json.Unmarshal(rawValue, &s); err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(common.HexToHash(s)))

		case common.Address:
			var s string
			if err := json.Unmarshal(rawValue, &s); err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(common.HexToAddress(s)))

		case *big.Int:
			bi, err := p.TransferBigInt(rawValue)
			if err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(bi))

		case types.Bloom:
			var s string
			if err := json.Unmarshal(rawValue, &s); err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(types.BytesToBloom(common.FromHex(s))))

		default:
			if err := json.Unmarshal(rawValue, fieldValue.Addr().Interface()); err != nil {
				return fmt.Errorf("failed to parse field %s: %w", jsonKey, err)
			}
		}
	}
	return nil
}
