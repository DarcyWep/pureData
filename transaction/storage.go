package transaction

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type storageTransition struct {
	label uint8 // 0: state, 1: storage

	contractAddress common.Address
	slot            common.Hash // 智能合约的存储槽
	preValue        common.Hash
	newValue        *common.Hash // newValue = nil 则是 SLOAD, 否则为 SSTORE
}

func newStorageTransition(contractAddress common.Address, slot common.Hash) *storageTransition {
	return &storageTransition{
		label: 1,

		contractAddress: contractAddress,
		slot:            slot,
		newValue:        new(common.Hash),
	}
}

func (s *storageTransition) addStorageTransition(preValue, newValue *common.Hash) {
	s.preValue = *preValue
	if newValue == nil {
		s.newValue = nil
	} else {
		*s.newValue = *newValue
	}
}

func (s *storageTransition) GetLabel() uint8 {
	return s.label
}

func (s *storageTransition) String() string {
	var str string = ""
	if s.newValue == nil {
		str = fmt.Sprintf("%d,%v,%v,%v,", s.label, s.contractAddress.Hex(), s.slot.Hex(), s.preValue.Hex())
	} else {
		str = fmt.Sprintf("%d,%v,%v,%v,%v", s.label, s.contractAddress.Hex(), s.slot.Hex(), s.preValue.Hex(), s.newValue.Hex())
	}
	//fmt.Println(str)
	return str
}
