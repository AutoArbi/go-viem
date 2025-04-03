package public

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type AccessListBuilder struct {
	list map[common.Address]map[common.Hash]struct{}
}

func NewAccessListBuilder() *AccessListBuilder {
	return &AccessListBuilder{
		list: make(map[common.Address]map[common.Hash]struct{}),
	}
}

func (b *AccessListBuilder) Add(address common.Address, storageKey common.Hash) {
	if b.list[address] == nil {
		b.list[address] = make(map[common.Hash]struct{})
	}
	b.list[address][storageKey] = struct{}{}
}

func (b *AccessListBuilder) AddAddressOnly(address common.Address) {
	if b.list[address] == nil {
		b.list[address] = make(map[common.Hash]struct{})
	}
}

func (b *AccessListBuilder) Build() types.AccessList {
	result := types.AccessList{}
	for addr, slots := range b.list {
		var keys []common.Hash
		for k := range slots {
			keys = append(keys, k)
		}
		result = append(result, types.AccessTuple{
			Address:     addr,
			StorageKeys: keys,
		})
	}
	return result
}
