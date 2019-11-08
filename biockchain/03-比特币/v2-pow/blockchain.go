package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type BlockChain struct {
	db   *bolt.DB //句柄，数据库对象handler
	tail []byte   //存储最后一个区块哈希值
}

const genesisInfo = "hello world"
const blockChainFileName = "blockchain.db"
const blockBucket = "blockBucket"
const lastBlockHashKey = "lastBlockHashKey"

//创建blockChain，同时添加一个创世块
func NewBlockChain() *BlockChain {
	var lastHash []byte

	//创建区块链时，向里面写入一个创世快
	//包含两个功能：
	//1. 如果区块链不存在，则创建，写入创世快
	db, err := bolt.Open(blockChainFileName, 0600, nil)
	if err != nil {
		fmt.Println("创建区块链失败, err:", err)
		return nil
	}

	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		if b == nil {
			//创建bucket
			b, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				fmt.Println("创建bucket失败, err:", err)
				return err
			}

			//写入创世块
			genesisBlock := NewBlock(genesisInfo, nil)
			//第一次：写入区块的数据
			_ = b.Put(genesisBlock.Hash, genesisBlock.Serialize() /*区块转换成字节流*/)
			//第二次：写入最后一个区块哈希
			_ = b.Put([]byte(lastBlockHashKey), genesisBlock.Hash)

			hash := b.Get([]byte(lastBlockHashKey))
			fmt.Printf("lastHash : %x\n", hash)

			//获取数据库中block序列化之后的数据
			blockInfo := b.Get(genesisBlock.Hash)
			block := Deserialize(blockInfo)
			fmt.Printf("block from db :%v\n", block)

			lastHash = genesisBlock.Hash
		} else {
			//bucket已经存在, 直接读取最后一区块的哈希值
			lastHash = b.Get([]byte(lastBlockHashKey))
			fmt.Printf("lastHash : %x\n", lastHash)
		}

		return nil
	})
	//2. 如果区块链存在，获取最后一个区块的哈希值

	//拼接成BlockChain实例返回

	return &BlockChain{db: db, tail: lastHash}
}

//1 <- 2 <-3
//添加区块的方法
func (bc *BlockChain) AddBlock(data string) {
	fmt.Println("AddBlock called!")
	//最后一个区块的哈希值
	lastHash := bc.tail

	//1. 创建新的区块
	newBlock := NewBlock(data, lastHash)

	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b == nil {
			log.Fatal("AddBlock是，bucket不应为空!")
		}

		//写入区块
		_ = b.Put(newBlock.Hash, newBlock.Serialize())
		_ = b.Put([]byte(lastBlockHashKey), newBlock.Hash)

		//更新tail的值
		bc.tail = newBlock.Hash
		return nil
	})
}


