package main

import "fmt"

func (cli *CLI) addBlock(data string) {
	//cli.bc.AddBlock(data) //TODO
}

func (cli *CLI) printBlock() {
	it := NewIterator(cli.bc)

	for {
		block := it.Next()

		fmt.Println("++++++++++++++++++++++")
		fmt.Printf("version::%s\n", block.Version)
		fmt.Printf("prevHash:%x\n", block.PrevHash)
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("merkleRoot:%x\n", block.MerkleRoot)
		fmt.Printf("timeStamp:%d\n", block.TimeStamp)
		fmt.Printf("bits:%d\n", block.Bits)
		fmt.Printf("nonce:%d\n", block.Nonce)
		pow := NewProofOfWork(block)
		fmt.Printf("isValid: %v\n", pow.isValid())
		fmt.Printf("data: %s\n", block.Transactions[0].TxInputs[0].ScriptSig)

		//判断是否已经是创世块
		if block.PrevHash == nil {
			break
		}
	}
	fmt.Println("区块链遍历结束!")
}

func (cli *CLI) getBalance(address string) {
	//utxos := cli.bc.FindMyUtxo(address)
	utxoinfos := cli.bc.FindMyUtxo(address)
	var total float64

	for _, utxoinfo := range utxoinfos {
		total += utxoinfo.output.Value
	}

	fmt.Printf("'%s'的比特币余额为:%f\n", address, total)
}

func (cli *CLI) send(from, to string, amount float64, miner, data string) {
	fmt.Printf("'%s'向'%s转账:'%f', miner:%s, data:%s\n", from, to, amount, miner, data)

	//输入数据的有效性会进行校验
	//TODO

	//创建挖矿交易
	coninbaseTx := NewCoinbaseTx(miner, data)
	txs := []*Transaction{coninbaseTx}

	//一个区块只添加一笔有效的普通交易
	tx, err := NewTransaction(from, to, amount, cli.bc)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("发现有效的交易，准备添加到区块, txid:%x\n", tx.Txid)
		txs = append(txs, tx)
	}

	//创建区块，添加到区块链
	cli.bc.AddBlock(txs)
}
