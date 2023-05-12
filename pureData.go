package pureData

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

var (
	blocksChan chan *BlockInfo
	blockDb    *leveldb.DB

	mergeChan chan *BlockInfo
	mergeDb   *leveldb.DB
)

const (
	blocksChanSize = 8192

	leveldbPath = "blocks"
	mergePath   = "merge"

	minCache   = 2048
	minHandles = 2048
)

func openLeveldb(path string, readOnly bool) (*leveldb.DB, bool) {
	db, err := leveldb.OpenFile(path, &opt.Options{
		OpenFilesCacheCapacity: minHandles,
		BlockCacheCapacity:     minCache / 2 * opt.MiB,
		WriteBuffer:            minCache / 4 * opt.MiB, // Two of these are used internally
		ReadOnly:               readOnly,
	})
	if err != nil {
		fmt.Println("Failed open leveldb")
		return nil, false
	}
	return db, true
}

// OpenInsertBlock 打开线程和相应的通道
func OpenInsertBlock() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(2)
	blocksChan = make(chan *BlockInfo, blocksChanSize)
	mergeChan = make(chan *BlockInfo, blocksChanSize)

	blockDb, _ = openLeveldb(leveldbPath, false)
	mergeDb, _ = openLeveldb(mergePath, false)

	go insertBlocks(&wg)
	go insertMerge(&wg)

	return &wg
}

func CloseInsertBlock(wg *sync.WaitGroup) {
	close(blocksChan)
	close(mergeChan)

	blockDb.Close()
	mergeDb.Close()

	wg.Wait()
}

// 插入数据库
func insertBlocks(wg *sync.WaitGroup) {
	defer wg.Done()
	for block := range blocksChan {
		str, err := blockDb.Get([]byte(block.number.String()), nil)
		fmt.Println(string(str), "\n")
		//err := blockDb.Put([]byte(block.number.String()), []byte(block.String()), nil)
		if err != nil {
			fmt.Println("Failed to store block, number is "+block.number.String()+", error is", err)
		}
		mergeChan <- block
	}
}

// 插入数据库
func insertMerge(wg *sync.WaitGroup) {
	defer wg.Done()
	for block := range mergeChan {
		str, err := mergeDb.Get([]byte(block.number.String()), nil)
		fmt.Println(string(str))
		//err := mergeDb.Put([]byte(block.number.String()), []byte(block.SimplisticString()), nil)
		if err != nil {
			fmt.Println("Failed to store block, number is "+block.number.String()+", error is", err)
		}
	}
}
