package transaction

import (
	"github.com/ethereum/go-ethereum/common"
)

type AccessSlot struct {
	IsRead  bool // 未读就是false
	IsWrite bool // 未写为false
}

func (a *AccessSlot) String() string {
	return Bool2String(a.IsRead) + Bool2String(a.IsWrite)
}

type AccessSlotMap map[common.Hash]*AccessSlot

func (a *AccessSlotMap) String() string {
	var str string = ""
	for slotKey, accessSlot := range *a {
		str = str + slotKey.Hex() + "~" + accessSlot.String() + "," // ~: 隔离访问状态, ,: 隔离slot
	}
	if len(str) > 0 {
		str = str[:len(str)-1]
	}
	return str
}

func NewAccessSlot() *AccessSlot {
	return &AccessSlot{
		IsRead:  false,
		IsWrite: false,
	}
}
func newAccessSlotMap() *AccessSlotMap {
	m := make(AccessSlotMap)
	return &m
}

type AccessAddress struct {
	Slots *AccessSlotMap
	//Slots   AccessSlotMap
	IsRead  bool // 未读就是false
	IsWrite bool // 未写为false

	// 粗粒度的读写，slot的读写记录到这里
	CoarseRead  bool
	CoarseWrite bool
}

func (a *AccessAddress) String() string {
	var str string = ""
	str = a.Slots.String() + "#" + Bool2String(a.IsRead) + Bool2String(a.IsWrite) + Bool2String(a.CoarseRead) + Bool2String(a.CoarseWrite)
	return str
}

type AccessAddressMap map[common.Address]*AccessAddress

func (a *AccessAddressMap) String() string {
	var str string = ""
	for addr, accessAddr := range *a {
		str = str + addr.Hex() + ":" + accessAddr.String() + " " // ~: 隔离地址和访问状态,  : 隔离每个地址
	}
	if len(str) > 0 {
		str = str[:len(str)-1]
	}
	return str
}

func NewAccessAddress() *AccessAddress {
	return &AccessAddress{
		Slots: newAccessSlotMap(),
		//Slots:   make(AccessSlotMap),
		IsRead:  false,
		IsWrite: false,

		CoarseRead:  false,
		CoarseWrite: false,
	}
}

func NewAccessAddressMap() *AccessAddressMap {
	m := make(AccessAddressMap)
	return &m
}

func Bool2String(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
