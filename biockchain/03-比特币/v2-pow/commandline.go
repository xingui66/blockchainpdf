package main

import "fmt"

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
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
		fmt.Printf("data: %s\n", block.Data)

		//判断是否已经是创世块
		if block.PrevHash == nil {
			break
		}
	}
	fmt.Println("区块链遍历结束!")
}
