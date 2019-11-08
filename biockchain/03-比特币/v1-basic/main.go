package main

import (
	"fmt"
)

func main() {
	bc := NewBlockChain()
	bc.AddBlock("helloworld")
	bc.AddBlock("你好")

	for _, block := range bc.blocks {
		fmt.Println("++++++++++++++++++++++")
		fmt.Printf("version::%s\n", block.version)
		fmt.Printf("prevHash:%x\n", block.prevHash)
		fmt.Printf("hash:%x\n", block.hash)
		fmt.Printf("merkleRoot:%x\n", block.merkleRoot)
		fmt.Printf("timeStamp:%d\n", block.timeStamp)
		fmt.Printf("bits:%d\n", block.bits)
		fmt.Printf("nonce:%d\n", block.nonce)
		fmt.Printf("data: %s\n", block.data)
	}

}
