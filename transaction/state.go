package transaction

import (
	"math/big"
)

// 以下是导出交易的相关字段需求 (注: 结构体里面的数据要给json包访问 -> 需要首字母大写)

type stateTransition struct {
	Label uint8 // 0: state, 1: storage
	// type 类型(1: 转账; 2: 手续费扣除, 只有From字段; 3: 手续费添加给矿工, 只有To字段 ; 4: 合约销毁; 5: 矿工奖励, 只有To字段)
	// 类型 5(矿工奖励) 每个区块只有一个记录
	Type_ uint8 // (手续费扣除 2!=3 给矿工的手续费)

	From  *balance
	To    *balance
	Value *big.Int
}

func newStateTransition(sender, recipient common.Address, sBefore, rBefore, value *big.Int) *stateTransition {
	return &stateTransition{
		label: 0,

		type_: 1,
		value: new(big.Int).Set(value),
		from:  newFullBalance(sender, sBefore),
		to:    newFullBalance(recipient, rBefore),
	}
}

func (t *stateTransition) GetLabel() uint8 {
	return t.label
}

type balance struct {
	address  string
	beforeTx *big.Int
}

func newFullBalance(addr string, before *big.Int) *balance {
	return &balance{
		address:  addr,
		beforeTx: new(big.Int).Set(before),
	}
}
