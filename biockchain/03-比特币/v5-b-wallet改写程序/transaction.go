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
	//ScriptSig string //先使用string代替，后续会改成签名

	//1. 私钥签名
	ScriptSig []byte
	//2. 公钥
	PubKey []byte
}

//交易输出
type TXOutput struct {
	//1. 锁定脚本
	//LockScript string

	//1收款人的公钥哈希
	PubKeyHash []byte
	//2. 转账金额
	Value float64
}

//收款人给付款人地址，锁定的时候不是使用地址锁定的，而是使用公钥哈希锁定的
//提供一个生成output的方法
func NewTXOutput(value float64, address string) TXOutput {
	//计算公钥哈希
	pubKeyHash := getPubKeyHashFromAddress(address)

	output := TXOutput{
		PubKeyHash: pubKeyHash,
		Value:      value,
	}

	return output
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
		ScriptSig: []byte(data),
		PubKey:    nil,
	}}

	output := NewTXOutput(reward, miner)
	outputs := []TXOutput{output}

	//outputs := []TXOutput{{
	//	LockScript: miner,
	//	Value:      reward,
	//}}

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

	//1. 打开钱包
	wm := NewWalletManager()
	if wm == nil {
		return nil, errors.New("打开钱包失败!")
	}

	//2. 找到付款方对应的私钥和公钥
	w, ok := wm.Wallets[from]
	if !ok {
		return nil, fmt.Errorf("没有找到：'%s'对应的钱包!", from)
	}
	//创建input的时候需要私钥签名和公钥
	//priKey := w.PrivKey //TODO 汇通在使用
	pubKey := w.PubKey

	//付款人的公钥哈希
	pubKeyHash := getPubKeyHashFromPubKey(pubKey)

	//1. 1. 找到付款人能够支配的合理的钱，返回金额和utxoinfo
	utxoinfos, value := bc.FindNeedUtxoInfo(pubKeyHash, amount)

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
			ScriptSig: nil,    //钱包再交易创建的最后再处理 TODO
			PubKey:    pubKey, //付款人的公钥
		}

		inputs = append(inputs, input)
	}
	//1. 遍历返回的utxonifo切片，逐个转成input结构
	//2. 拼装outputs
	//1. 拼装一个属于收款人的output
	output := NewTXOutput(amount, to)

	//output := TXOutput{
	//	LockScript: to,
	//	Value:      amount,
	//}
	outputs = append(outputs, output)

	//2. 判断一下是否需要找零，如果有，拼装一个属于付款方output
	if value > amount {
		//找零
		//output1 := TXOutput{
		//	LockScript: from,
		//	Value:      value - amount,
		//}
		output1 := NewTXOutput(value-amount, from)
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
