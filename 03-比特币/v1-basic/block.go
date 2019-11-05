package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

//定义区块结构
//区块头
//区块体
type Block struct {
	//版本号
	version string

	//前区块哈希
	prevHash []byte

	//当前区块的哈希
	//在比特币定义的区块中，是没有这个当前区块哈希字段的，
	//为了方便处理，添加一个字段
	hash []byte

	//merkle根, 根据当前区块的交易数据计算出来的
	merkleRoot []byte

	//时间戳
	timeStamp int64

	//难度值，系统提供的
	bits int64

	//随机数
	nonce int64

	//区块体，交易数据
	data []byte
}

func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{
		version:    "0",
		prevHash:   prevHash,
		merkleRoot: nil,
		timeStamp:  time.Now().Unix(),
		bits:       0,
		nonce:      0,
		data:       []byte(data),
	}

	//挖矿过程暂且省略
	block.setHash()

	return block
}

func (b *Block) setHash() {
	//将当前区块的数据拼接到一起
	//data1 := append([]byte(b.version), b.prevHash...)
	//data1 = append(data1, b.data...)

	tmp := [][]byte{
		[]byte(b.version),
		b.prevHash,
		b.merkleRoot,
		//[]byte(string(b.timeStamp)),
		//[]byte(string(b.bits)),
		//[]byte(string(b.nonce)),
		digi2byte(b.timeStamp),
		digi2byte(b.bits),
		digi2byte(b.nonce),
		b.data,
	}

	data1 := bytes.Join(tmp, []byte(""))
	fmt.Printf("tmp data:%x\n", data1)
	fmt.Printf("time stamp : %s\n", []byte(string(b.timeStamp)), )

	//哈希运算
	hash := sha256.Sum256(data1)
	b.hash = hash[:]
}
