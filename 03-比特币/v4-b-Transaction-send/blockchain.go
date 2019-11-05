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
			//创建一个挖矿交易，里面写入创世语
			coinbaseTx := NewCoinbaseTx("中本聪", genesisInfo)
			genesisBlock := NewBlock([]*Transaction{coinbaseTx}, nil)

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
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	fmt.Println("AddBlock called!")
	//最后一个区块的哈希值
	lastHash := bc.tail

	//1. 创建新的区块
	newBlock := NewBlock(txs, lastHash)

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

//定义一个结构，同时包含output和它的位置信息
type UtxoInfo struct {
	output TXOutput
	index  int64
	txid   []byte
}

//遍历账本，查询指定地址所有的utxo
func (bc *BlockChain) FindMyUtxo(address string) []UtxoInfo {
	fmt.Println("FindMyUtxo called, address:", address)

	//var outputs []TXOutput
	var utxoinfos []UtxoInfo

	//定义一个map，用于存储已经消耗过的output
	//key ==> 交易id， value：在这个交易中的索引的切片
	spentOutput := make(map[string][]int64)
	//map[0x2222] = {0}
	//map[0x3333] = {0, 1}

	//1. 遍历区块
	it := NewIterator(bc)

	for {
		block := it.Next()
		//2. 遍历交易
		for _, tx := range block.Transactions {
		LABEL1:
			//3. 遍历output
			for outputIndex, output := range tx.TXOutputs {
				//判断当前的output是否是目标地址锁定的
				if output.LockScript == address {
					//再添加之前进行过滤 ，依据：spentOutput集合
					//1. 先查看当前交易（0x3333）是否已经存在于spentOutput容器中
					currTxId := string(tx.Txid)
					//{0, 1}
					indexArr := spentOutput[currTxId]

					if len(indexArr) != 0 {
						//说明容器中存在当前交易的output
						for _, spentIndex /*0, 1*/ := range indexArr {
							if outputIndex == int(spentIndex) {
								fmt.Println("当前的output已经被使用过了，无需统计, index:", outputIndex)
								continue LABEL1
							}
						}
					}
					//2. 如果不存在，则直接添加
					//3. 如果存在，进一步查看当前的output是否存在于这个容器//map[0x3333] = {0, 1}
					//	a. 获取这个交易id（0x333）对应的数组值indexArray : {0, 1}
					//  b. 判断当前索引是否属于{0, 1}

					fmt.Printf("找到了属于'%s'的output, index:%d, value:%f\n", address, outputIndex, output.Value)
					utxoinfo := UtxoInfo{
						output: output,
						index:  int64(outputIndex),
						txid:   tx.Txid,
					}

					utxoinfos = append(utxoinfos, utxoinfo)
					//outputs = append(outputs, output)
				}
			}

			//遍历inputs， 得到一个map
			if !tx.isCoinbaseTx() {
				//如果不是挖矿交易，才有必要遍历innputs
				for _, input := range tx.TxInputs {
					if input.ScriptSig == address {
						spentKey := string(input.TXID) //这个input的来源
						spentOutput[spentKey] = append(spentOutput[spentKey], input.Index)

						//下面是错误的，无法将数据写到spentOutput中
						//indexArray := spentOutput[spentKey]
						//indexArray = append(indexArray, input.Index)
					}
				}
			}
		}

		if block.PrevHash == nil {
			break
		}
	}

	return utxoinfos
}

func (bc *BlockChain) FindNeedUtxoInfo(address string, amount float64) ([]UtxoInfo, float64) {
	fmt.Printf("FindNeedUtxoInfo called, address :%s, amount:%f\n", address, amount)

	//1. 遍历账本，找到所有address（付款人）的utxo集合
	utxoinfos := bc.FindMyUtxo(address)

	//返还的utxoinfo里面包含金额
	var retValue float64
	var retUtxoInfo []UtxoInfo

	//2. 筛选出满足条件的数量即可，不要全部返还
	for _, utxoinfo := range utxoinfos {
		retUtxoInfo = append(retUtxoInfo, utxoinfo)
		retValue += utxoinfo.output.Value

		if retValue >= amount {
			//满足转账需求，直接返回
			break
		}
	}

	return retUtxoInfo, retValue
}
