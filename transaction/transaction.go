package transaction

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strconv"
	"strings"
)

var big0 = new(big.Int).SetInt64(0)

type Transfer interface { // 转账
	GetLabel() uint8 // 0: 普通转账(state), 1: ERC20类转账(storage), 2: KECCAK256
	String() string
}

type Transaction struct {
	BlockNumber *big.Int
	Hash        *common.Hash
	From        *common.Address
	To          *common.Address
	Index       *big.Int
	Value       *big.Int

	Transfers []Transfer
}

func newTransaction(number *big.Int, hash *common.Hash, from, to *common.Address, index, value *big.Int) *Transaction {
	tx := &Transaction{}
	tx.BlockNumber = new(big.Int).Set(number)
	if hash.Hex() != (common.Hash{}).Hex() { // 非最后的区块奖励部分
		tx.initTransaction(hash, from, to, index, value)
	}
	return tx
}

func (tx *Transaction) initTransaction(hash *common.Hash, from, to *common.Address, index, value *big.Int) {
	tx.Hash = new(common.Hash)
	tx.From = new(common.Address)
	tx.To = new(common.Address)

	tx.Hash.SetBytes(hash.Bytes())
	tx.From.SetBytes(from.Bytes())
	if to != nil {
		tx.To.SetBytes(to.Bytes())
	} else {
		tx.To = nil
	}
	tx.Index = new(big.Int).Set(index)
	tx.Value = new(big.Int).Set(value)
	tx.Transfers = make([]Transfer, 0)
}

func UnmarshalTransaction(txStr string) *Transaction {
	infoStrs := strings.Split(txStr, "|")
	var (
		tmp = big.NewInt(0)

		number *big.Int
		hash   = new(common.Hash)
		from   = new(common.Address)
		to     = new(common.Address)
		index  *big.Int
		value  *big.Int
	)

	//for _, infoStr := range infoStrs {
	//	fmt.Println(infoStr)
	//}

	tmp.SetString(infoStrs[0], 10)
	number = new(big.Int).Set(tmp)

	hash.SetBytes(common.FromHex(infoStrs[1]))
	from.SetBytes(common.FromHex(infoStrs[2]))
	if len(infoStrs[3]) != 0 {
		to.SetBytes(common.FromHex(infoStrs[3]))
	} else {
		to = nil
	}

	tmp.SetString(infoStrs[4], 10)
	index = new(big.Int).Set(tmp)

	tmp.SetString(infoStrs[5], 10)
	value = new(big.Int).Set(tmp)

	transferStrs := strings.Split(infoStrs[6], " ")

	tx := newTransaction(number, hash, from, to, index, value)
	for _, transferStr := range transferStrs {
		if transferStr[0] == '0' {
			tx.Transfers = append(tx.Transfers, Transfer(unmarshalStateTransition(transferStr)))
		} else {
			tx.Transfers = append(tx.Transfers, Transfer(unmarshalStorageTransition(transferStr)))
		}
		//fmt.Println(transferStr)
		//fmt.Println(tx.Transfers[i])
		//fmt.Println()
	}

	return tx
}

func unmarshalStorageTransition(stStr string) *StorageTransition {
	stStrs := strings.Split(stStr, ",")

	var (
		contract common.Address
		slot     common.Hash
		preValue common.Hash
		newValue *common.Hash = new(common.Hash)
	)
	contract.SetBytes(common.FromHex(stStrs[1]))
	slot.SetBytes(common.FromHex(stStrs[2]))
	preValue.SetBytes(common.FromHex(stStrs[3]))
	if stStrs[4] != "" {
		newValue.SetBytes(common.FromHex(stStrs[4]))
	} else {
		newValue = nil
	}

	return newStorageTransition(contract, slot, &preValue, newValue)
}

func unmarshalStateTransition(stStr string) *StateTransition {
	stStrs := strings.Split(stStr, ",")

	var (
		tmp   = big.NewInt(0)
		value *big.Int
	)
	type_, _ := strconv.Atoi(stStrs[1])

	tmp.SetString(stStrs[4], 10)
	value = new(big.Int).Set(tmp)

	return newStateTransition(type_, unmarshalBalance(stStrs[2]), unmarshalBalance(stStrs[3]), value)
}

func unmarshalBalance(bStr string) *Balance {
	if bStr == "" {
		return nil
	}
	bStrs := strings.Split(bStr, "~")

	var (
		tmp     = big.NewInt(0)
		addr    common.Address
		balance *big.Int
	)
	addr.SetBytes(common.FromHex(bStrs[0]))

	tmp.SetString(bStrs[1], 10)
	balance = new(big.Int).Set(tmp)

	return newBalance(addr, balance)
}

func (tx Transaction) String() string {
	var (
		from  = ""
		to    = ""
		index = ""
		value = ""

		trs string = ""
	)

	for _, tr := range tx.Transfers { // 空格分离
		trs += tr.String() + " "
	}

	if tx.Hash == nil {
		return fmt.Sprintf("||||||%v", trs[:len(trs)-1])
	}

	if tx.From != nil {
		from = tx.From.Hex()
	}
	if tx.To != nil {
		to = tx.To.Hex()
	}
	if tx.Index != nil {
		index = tx.Index.String()
	}
	if tx.Value != nil {
		value = tx.Value.String()
	}

	if trs != "" {
		trs = trs[0 : len(trs)-1]
	}
	//fmt.Println(number.String(), tx.Hash.Hex(), from, to, index, value)
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", tx.BlockNumber.String(), tx.Hash.Hex(), from, to, index, value, trs)
}
