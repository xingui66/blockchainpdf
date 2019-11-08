package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

//定义区块结构
//区块头
//区块体
type Block struct {
	//版本号
	Version string

	//前区块哈希
	PrevHash []byte

	//当前区块的哈希
	//在比特币定义的区块中，是没有这个当前区块哈希字段的，
	//为了方便处理，添加一个字段
	Hash []byte

	//merkle根, 根据当前区块的交易数据计算出来的
	MerkleRoot []byte

	//时间戳
	TimeStamp int64

	//难度值，系统提供的
	Bits int64

	//随机数
	Nonce int64

	//区块体，交易数据
	//Data []byte
	Transactions []*Transaction //多条交易
}

func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{
		Version:      "0",
		PrevHash:     prevHash,
		MerkleRoot:   nil,
		TimeStamp:    time.Now().Unix(),
		Bits:         0,
		Nonce:        0,
		Transactions: txs,
	}

	//挖矿过程暂且省略
	//block.setHash()
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash

	return block
}

//序列化区块
//结构=》编码 =》字节流
func (b *Block) Serialize() []byte {
	//创建编码器
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	//编码
	err := encoder.Encode(b)
	if err != nil {
		fmt.Println("区块编码失败, err:", err)
		return nil
	}
	return buff.Bytes()
}

//反序列化区块
//字节流=》解码=》结构
func Deserialize(data []byte) *Block {
	var block Block
	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))

	//解码
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("区块解析字节流失败,err:", err)
		return nil
	}
	return &block
}
