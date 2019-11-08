package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

//交易结构
type Transaction struct {
	//交易ID
	Txid []byte

	//多个交易输入
	TxInputs []TXInput

	//多个交易输出
	TXOutputs []TXOutput

	//时间戳
	TimeStamp int64
}

//交易输入
type TXInput struct {
	//1. 所引用的output所在的交易id
	TXID []byte
	//2. 所引用的output的索引值
	Index int64
	//3. 解锁脚本：
	ScriptSig string //先使用string代替，后续会改成签名
	//1. 私钥签名
	//2. 公钥
}

//交易输出
type TXOutput struct {
	//1. 锁定脚本
	LockScript string
	//2. 转账金额
	Value float64
}

//设置当前交易的id，使用交易本身的哈希值作为自己交易id
func (tx *Transaction) SetTxId() {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		fmt.Println("设置交易id失败, err:", err)
		return
	}

	hash := sha256.Sum256(buff.Bytes())

	tx.Txid = hash[:]
}

const reward = 12.5

//挖矿交易
//没有引用的输入, 只有输出，只有一个output
//对于挖矿交易，不需要签名，所以可以由矿工任意填写一个数据
func NewCoinbaseTx(miner string, data string) *Transaction {
	intputs := []TXInput{{
		TXID:      nil,
		Index:     -1,
		ScriptSig: data,
	}}

	outputs := []TXOutput{{
		LockScript: miner,
		Value:      reward,
	}}

	tx := &Transaction{
		TxInputs:  intputs,
		TXOutputs: outputs,
		TimeStamp: time.Now().Unix(),
	}

	//设置交易id
	tx.SetTxId()

	return tx
}

//普通交易
func NewTransaction() *Transaction {
	//TODO
	return nil
}
