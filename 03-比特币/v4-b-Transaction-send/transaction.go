package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
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

//判断一个交易是否为挖矿交易
func (tx *Transaction) isCoinbaseTx() bool {
	input := tx.TxInputs[0]
	if len(tx.TxInputs) == 1 && input.TXID == nil && input.Index == -1 {
		return true
	}

	return false
}

//普通交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) (*Transaction, error) {

	//1. 1. 找到付款人能够支配的合理的钱，返回金额和utxoinfo
	utxoinfos, value := bc.FindNeedUtxoInfo(from, amount)

	//2. 判断返回金额是否满足转账条件，如果不满足，创建交易失败。
	if value < amount {
		return nil, errors.New("付款人金额不足!")
	}

	//3. 拼接一个新的交易
	var inputs []TXInput
	var outputs []TXOutput

	//1. 拼装inputs
	for _, utxoinfo := range utxoinfos {
		input := TXInput{
			TXID:      utxoinfo.txid,
			Index:     utxoinfo.index,
			ScriptSig: from,
		}

		inputs = append(inputs, input)
	}
	//1. 遍历返回的utxonifo切片，逐个转成input结构
	//2. 拼装outputs
	//1. 拼装一个属于收款人的output
	output := TXOutput{
		LockScript: to,
		Value:      amount,
	}
	outputs = append(outputs, output)

	//2. 判断一下是否需要找零，如果有，拼装一个属于付款方output
	if value > amount {
		//找零
		output1 := TXOutput{
			LockScript: from,
			Value:      value - amount,
		}

		outputs = append(outputs, output1)
	}

	tx := Transaction{
		TxInputs:  inputs,
		TXOutputs: outputs,
		TimeStamp: time.Now().Unix(),
	}

	//3. 设置交易id
	tx.SetTxId()

	//4. 返回
	return &tx, nil
}
