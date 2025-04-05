package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestAccessListBuilderEmpty(t *testing.T) {
	builder := NewAccessListBuilder()
	accessList := builder.Build()
	if len(accessList) != 0 {
		t.Errorf("expected empty access list, got %d entries", len(accessList))
	}
}

func TestAccessListBuilderAdd(t *testing.T) {
	builder := NewAccessListBuilder()
	addr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	storageKey := common.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	builder.Add(addr, storageKey)

	accessList := builder.Build()
	if len(accessList) != 1 {
		t.Fatalf("expected 1 entry in access list, got %d", len(accessList))
	}
	tuple := accessList[0]
	if tuple.Address != addr {
		t.Errorf("expected address %s, got %s", addr.Hex(), tuple.Address.Hex())
	}
	if len(tuple.StorageKeys) != 1 {
		t.Errorf("expected 1 storage key, got %d", len(tuple.StorageKeys))
	}
	if tuple.StorageKeys[0] != storageKey {
		t.Errorf("expected storage key %s, got %s", storageKey.Hex(), tuple.StorageKeys[0].Hex())
	}
}

func TestAccessListBuilderAddAddressOnly(t *testing.T) {
	builder := NewAccessListBuilder()
	addr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	builder.AddAddressOnly(addr)

	accessList := builder.Build()
	if len(accessList) != 1 {
		t.Fatalf("expected 1 entry in access list, got %d", len(accessList))
	}
	tuple := accessList[0]
	if tuple.Address != addr {
		t.Errorf("expected address %s, got %s", addr.Hex(), tuple.Address.Hex())
	}
	if len(tuple.StorageKeys) != 0 {
		t.Errorf("expected 0 storage keys, got %d", len(tuple.StorageKeys))
	}
}

func TestAccessListBuilderMultipleKeys(t *testing.T) {
	builder := NewAccessListBuilder()
	addr := common.HexToAddress("0x3333333333333333333333333333333333333333")
	key1 := common.HexToHash("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	key2 := common.HexToHash("0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")
	builder.Add(addr, key1)
	builder.Add(addr, key2)

	accessList := builder.Build()
	if len(accessList) != 1 {
		t.Fatalf("expected 1 entry in access list, got %d", len(accessList))
	}
	tuple := accessList[0]
	if tuple.Address != addr {
		t.Errorf("expected address %s, got %s", addr.Hex(), tuple.Address.Hex())
	}
	if len(tuple.StorageKeys) != 2 {
		t.Fatalf("expected 2 storage keys, got %d", len(tuple.StorageKeys))
	}
	// 检查两个 storage key 是否都存在（顺序不定）
	foundKey1, foundKey2 := false, false
	for _, k := range tuple.StorageKeys {
		if k == key1 {
			foundKey1 = true
		}
		if k == key2 {
			foundKey2 = true
		}
	}
	if !foundKey1 || !foundKey2 {
		t.Errorf("expected both storage keys to be present, got %v", tuple.StorageKeys)
	}
}

func TestAccessListBuilderMultipleAddresses(t *testing.T) {
	builder := NewAccessListBuilder()
	// 对地址1添加一个 storage key
	addr1 := common.HexToAddress("0x4444444444444444444444444444444444444444")
	key1 := common.HexToHash("0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	builder.Add(addr1, key1)
	// 对地址2只添加地址
	addr2 := common.HexToAddress("0x5555555555555555555555555555555555555555")
	builder.AddAddressOnly(addr2)

	accessList := builder.Build()
	if len(accessList) != 2 {
		t.Fatalf("expected 2 entries in access list, got %d", len(accessList))
	}

	// 检查每个地址对应的数据
	for _, tuple := range accessList {
		switch tuple.Address.Hex() {
		case addr1.Hex():
			if len(tuple.StorageKeys) != 1 || tuple.StorageKeys[0] != key1 {
				t.Errorf("for addr1 expected storage key %s, got %v", key1.Hex(), tuple.StorageKeys)
			}
		case addr2.Hex():
			if len(tuple.StorageKeys) != 0 {
				t.Errorf("for addr2 expected no storage keys, got %v", tuple.StorageKeys)
			}
		default:
			t.Errorf("unexpected address in access list: %s", tuple.Address.Hex())
		}
	}
}
