package pureData

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"math/big"
	"testing"
)

const (
	minCache     = 2048
	minHandles   = 2048
	nativeDbPath = "blocks"
	mergeDbPath  = "merge"
)

func openLeveldb(path string) (*leveldb.DB, error) {
	return leveldb.OpenFile(path, &opt.Options{
		OpenFilesCacheCapacity: minHandles,
		BlockCacheCapacity:     minCache / 2 * opt.MiB,
		WriteBuffer:            minCache / 4 * opt.MiB, // Two of these are used internally
		ReadOnly:               true,
	})
}

func Test(t *testing.T) {
	db, err := openLeveldb(nativeDbPath) // get native transaction or merge transaction
	//db, err := openLeveldb(mergeDbPath) // get native transaction or merge transaction
	defer db.Close()
	if err != nil {
		fmt.Println("open leveldb error,", err)
		return
	}
	number := new(big.Int).SetInt64(11090501)
	txs, _ := GetTransactionsByNumber(db, number)
	for _, tx := range txs {
		fmt.Println(tx.String())
	}
}
