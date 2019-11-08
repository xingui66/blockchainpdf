package main

type BlockChain struct {
	blocks []*Block
}

const genesisInfo = "天气转凉，注意御寒!"

//创建blockchain，同时添加一个创世块
func NewBlockChain() *BlockChain {
	genesisBlock := NewBlock(genesisInfo, nil)
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}

//1 <- 2 <-3
//添加区块的方法
func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevHash := lastBlock.hash

	//1. 创建新的区块
	newBlock := NewBlock(data, prevHash)

	//2. 添加到bc的blocks
	bc.blocks = append(bc.blocks, newBlock)
}
