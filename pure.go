package pureData

import (
	"github.com/syndtr/goleveldb/leveldb"
	"math/big"
	"pureData/transaction"
	"strings"
)

func GetTransactionsByNumber(db *leveldb.DB, number *big.Int) ([]*transaction.Transaction, error) {
	block, err := db.Get([]byte(number.String()), nil)
	if err != nil {
		return nil, err
	}
	txStrs := strings.Split(string(block), ";")

	txs := make([]*transaction.Transaction, 0)
	for _, txStr := range txStrs {
		txs = append(txs, transaction.UnmarshalTransaction(txStr))
	}
	return txs, nil
}
