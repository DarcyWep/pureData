package pureData

import (
	"github.com/DarcyWep/pureData/transaction"
	"github.com/syndtr/goleveldb/leveldb"
	"math/big"
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

//const (
//	minCache     = 2048
//	minHandles   = 2048
//	nativeDbPath = "/home/eth/Project/morph/ethereumdata/nativedb"
//	mergeDbPath  = "/home/eth/Project/morph/ethereumdata/mergedb"
//)
//
//func openLeveldb(path string) (*leveldb.DB, error) {
//	return leveldb.OpenFile(path, &opt.Options{
//		OpenFilesCacheCapacity: minHandles,
//		BlockCacheCapacity:     minCache / 2 * opt.MiB,
//		WriteBuffer:            minCache / 4 * opt.MiB, // Two of these are used internally
//		ReadOnly:               true,
//	})
//}
//
//func main() {
//	db, err := openLeveldb(nativeDbPath) // get native transaction or merge transaction
//	//db, err := openLeveldb(mergeDbPath) // get native transaction or merge transaction
//	defer db.Close()
//	if err != nil {
//		fmt.Println("open leveldb error,", err)
//		return
//	}
//	//number := new(big.Int).SetInt64(11090501)
//	//min, max, addSpan := big.NewInt(12000001), big.NewInt(12050000), big.NewInt(1)
//	min, max, addSpan := big.NewInt(14000001), big.NewInt(14000049), big.NewInt(1)
//	for i := min; i.Cmp(max) == -1; i = i.Add(i, addSpan) {
//		txs, err := GetTransactionsByNumber(db, i)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		for _, tx := range txs {
//			//if len(tx.Transfers) == 1 {
//			//	fmt.Println(i.String(), tx.Hash, tx.Transfers)
//			//}
//			fmt.Println(i.String(), tx.Hash, tx.AccessAddress.String())
//		}
//	}
//
//}
