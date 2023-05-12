package transaction

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"morphDataGeted/core/types"
)

var big0 = new(big.Int).SetInt64(0)

type transfer interface { // 转账
	GetLabel() uint8 // 0: 普通转账(state), 1: ERC20类转账(storage), 2: KECCAK256
	String() string
}

type Transaction struct {
	// 当前交易保存的临时变量
	stateTransition   *stateTransition
	storageTransition *storageTransition

	BlockNumber *big.Int
	Hash        *common.Hash
	From        *common.Address
	To          *common.Address
	Index       *big.Int
	Value       *big.Int

	Transfers []transfer
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

func NewTransaction(tx *types.Transaction, from *common.Address, blockNumber *big.Int, index int) *Transaction {
	t := &Transaction{}
	t.init(tx, from, blockNumber, index)
	return t
}

func (tx *Transaction) init(transaction *types.Transaction, from *common.Address, blockNumber *big.Int, index int) {
	if transaction == nil { // 叔父区块奖励 + 挖矿奖励
		tx.Hash = nil
		tx.BlockNumber = new(big.Int).Set(blockNumber)
	} else {
		tx.initTransaction(transaction, from, blockNumber, index)
	}
}

func (tx *Transaction) initTransaction(transaction *types.Transaction, from *common.Address, blockNumber *big.Int, index int) {
	tx.BlockNumber = new(big.Int).Set(blockNumber)

	h := transaction.Hash()
	tx.Hash = &h

	tx.From = from
	tx.To = transaction.To()
	tx.Index = new(big.Int).SetInt64(int64(index))
	tx.Value = new(big.Int).Set(transaction.Value())

	tx.Transfers = make([]transfer, 0)
}

func (tx *Transaction) InitStateTransition(type_ uint8, sender, recipient *common.Address, value *big.Int) {
	tx.stateTransition = newStateTransition(type_, sender, recipient, value)
}

func (tx *Transaction) SetStateTransitionBalance(isFrom bool, beforeTx *big.Int) {
	tx.stateTransition.setBalance(isFrom, beforeTx)
}

func (tx *Transaction) AddStateTransition() {
	tx.Transfers = append(tx.Transfers, transfer(tx.stateTransition))
}

func (tx *Transaction) InitStorageTransition(contractAddress common.Address, slot common.Hash) {
	tx.storageTransition = newStorageTransition(contractAddress, slot)
}

func (tx *Transaction) AddStorageTransition(preValue, newValue *common.Hash) {
	tx.storageTransition.addStorageTransition(preValue, newValue)
	tx.Transfers = append(tx.Transfers, transfer(tx.storageTransition))
}

func (tx *Transaction) Snapshot() int {
	return len(tx.Transfers)
}

func (tx *Transaction) RevertToSnapshot(revid int) {
	if revid != 0 {
		tx.Transfers = tx.Transfers[:revid]
	}
}

func (tx Transaction) SimplisticString() string {
	var (
		from  = ""
		to    = ""
		index = ""
		value = ""

		trs string = ""
	)

	if tx.Hash == nil { // 矿工奖励
		for _, tr := range tx.Transfers { // 空格分离
			trs += tr.String() + " "
		}
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

	trs = tx.simplisticTransfersString()
	//fmt.Println(tx.BlockNumber.String(), tx.Hash.Hex(), from, to, index, value)
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", tx.BlockNumber.String(), tx.Hash.Hex(), from, to, index, value, trs)
}

func (tx Transaction) simplisticTransfersString() string {
	var trs4_5, trs string = "", ""

	var fee, reward *stateTransition
	stateTransfers := make([]*stateTransition, 0)
	for _, tr := range tx.Transfers { // 空格分离
		if tr.GetLabel() == 0 {
			stateTr := tr.(*stateTransition)
			switch stateTr.type_ {
			case 1:
				stateTransfers = append(stateTransfers, stateTr)
			case 2:
				fee = stateTr
			case 3:
				reward = stateTr
			case 4:
				tmp := tr.String()
				fmt.Println(tmp)
				trs4_5 += tmp + " "
			case 5:
				tmp := tr.String()
				fmt.Println(tmp)
				trs4_5 += tmp + " "
			}
		}
	}

	switch len(stateTransfers) {
	case 0:
		trs = feeAndRewardString(fee, reward) + " "
	case 1:
		tr1 := stateTransfers[0]
		if tr1.from.address == fee.from.address {
			tr1.from.beforeTx.Sub(tr1.from.beforeTx, fee.value)
		}
		tmp := tr1.String()
		//fmt.Println(tmp)
		trs = feeAndRewardString(fee, reward) + " " + tmp + " "
	default:
		sTrs := mergeStateTransfer(fee, reward, stateTransfers)
		trs = feeAndRewardString(fee, reward) + " "
		for _, tr := range sTrs {
			tmp := tr.String()
			//fmt.Println(tmp)
			trs += tmp + " "
		}
	}
	if trs4_5 != "" {
		trs += trs4_5 + " "
	}
	return trs[:len(trs)-1]
}

func feeAndRewardString(fee, reward *stateTransition) string {
	var trs string = ""
	fee.value.Sub(fee.value, reward.value)
	reward.from = fee.from.deepCopy()
	if fee.value.Cmp(big0) != 0 { // 手续费有燃烧
		tmp := fee.String()
		//fmt.Println(tmp)
		trs = tmp + " "
	}
	tmp := reward.String()
	//fmt.Println(tmp)
	trs += tmp
	//json.Unmarshal()
	return trs
}
func mergeStateTransfer(fee, reward *stateTransition, transfers []*stateTransition) []*stateTransition {
	var valueAmend = new(big.Int).Set(fee.from.beforeTx)
	valueAmend.Sub(valueAmend, fee.value) // 修正余额
	for _, tr := range transfers {
		if tr.from.address == fee.from.address {
			tr.from.beforeTx.Set(valueAmend)
			valueAmend.Sub(valueAmend, tr.value) // 修正余额
		}
		if tr.to.address == fee.from.address {
			tr.to.beforeTx.Set(valueAmend)
			valueAmend.Add(valueAmend, tr.value) // 修正余额
		}
	}

	var addrTransfers map[string]*addressRelatedTransfers = make(map[string]*addressRelatedTransfers, 0)
	for i, tr := range transfers {
		var nodef, nodet *addressRelatedTransfers
		var okf, okt bool
		nodef, okf = addrTransfers[tr.from.address]
		nodet, okt = addrTransfers[tr.to.address]
		if !okf {
			nodef = newAddressRelatedTransfers(tr.from.address)
			nodef.before.Set(tr.from.beforeTx)
			addrTransfers[tr.from.address] = nodef
		}
		if !okt {
			nodet = newAddressRelatedTransfers(tr.to.address)
			nodet.before.Set(tr.to.beforeTx)
			addrTransfers[tr.to.address] = nodet
		}

		nodef.txIndexes = append(nodef.txIndexes, i)
		nodef.outValue.Add(nodef.outValue, tr.value) // from节点支出
		nodef.outTrs = append(nodef.outTrs, tr)

		nodet.txIndexes = append(nodet.txIndexes, i)
		nodet.inValue.Add(nodet.inValue, tr.value) // to节点接收
		nodet.inTrs = append(nodet.inTrs, tr)
	}

	var inNodes []*addressRelatedTransfers = make([]*addressRelatedTransfers, 0)
	var outNodes []*addressRelatedTransfers = make([]*addressRelatedTransfers, 0)
	for _, node := range addrTransfers {
		//fmt.Println("test:", addr, node.outValue, node.inValue, node.outValue.Cmp(node.inValue))
		if node.inValue.Cmp(node.outValue) < 0 { // 收入少于支出，支出账户
			node.outValue.Sub(node.outValue, node.inValue)
			outNodes = append(outNodes, node)
			//fmt.Println("outNode:", addr, node.outValue, node.inValue, node.outValue.Cmp(node.inValue))
			//fmt.Println()
			continue
		}
		if node.inValue.Cmp(node.outValue) > 0 { // 收入多于支出，收入账户
			node.inValue.Sub(node.inValue, node.outValue)
			inNodes = append(inNodes, node)
			//fmt.Println("inNode:", addr, node.outValue, node.inValue, node.outValue.Cmp(node.inValue))
			//fmt.Println()
		}
	}

	// 重新构建交易
	stateTransitions := make([]*stateTransition, 0)
	var inIndex, inLen, outIndex = 0, len(inNodes), 0
	var out *addressRelatedTransfers
	if len(outNodes) == 0 || inLen == 0 {
		fmt.Println("循环花费？？？")
		return stateTransitions
	}
	for outIndex, out = range outNodes {
		for {
			if out.outValue.Cmp(big0) == 0 { // 当前账户完成支出
				break
			}
			if out.outValue.Cmp(inNodes[inIndex].inValue) >= 0 { // out 账户足够支出当前花费
				//fmt.Println("construct:")
				stateTransitions = append(stateTransitions, newFullStateTransition(out.address, inNodes[inIndex].address, out.before, inNodes[inIndex].before, inNodes[inIndex].inValue))
				out.before.Sub(out.before, inNodes[inIndex].inValue)     // 修正状态
				out.outValue.Sub(out.outValue, inNodes[inIndex].inValue) // 修正状态
				inIndex += 1
			} else { // out 账户不足以支出当前花费
				stateTransitions = append(stateTransitions, newFullStateTransition(out.address, inNodes[inIndex].address, out.before, inNodes[inIndex].before, out.outValue))
				inNodes[inIndex].before.Add(inNodes[inIndex].before, out.outValue) // 修正状态
				out.outValue.Sub(out.outValue, out.outValue)                       // 修正状态
			}
		}
	}
	if outIndex != len(outNodes)-1 || inLen != inIndex {
		fmt.Println("支出与收入不匹配？？？")
		return make([]*stateTransition, 0)
	}
	return stateTransitions
}

type addressRelatedTransfers struct {
	address   string
	txIndexes []int
	before    *big.Int
	inValue   *big.Int // 收入
	outValue  *big.Int // 支出
	inTrs     []*stateTransition
	outTrs    []*stateTransition
}

func newAddressRelatedTransfers(addr string) *addressRelatedTransfers {
	return &addressRelatedTransfers{
		address:   addr,
		txIndexes: make([]int, 0),
		before:    new(big.Int).SetInt64(0),
		inValue:   new(big.Int).SetInt64(0),
		outValue:  new(big.Int).SetInt64(0),
		inTrs:     make([]*stateTransition, 0),
		outTrs:    make([]*stateTransition, 0),
	}
}
