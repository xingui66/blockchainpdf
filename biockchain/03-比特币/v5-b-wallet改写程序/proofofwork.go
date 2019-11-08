package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const difficult = 16 // 三个前导0
//const difficult = 20 // 四个前导0

//定义一个结构：ProofOfWork，包含两个字段：
type ProofOfWork struct {
	//1. 区块数据block
	block *Block

	//2. 难度值target
	target *big.Int
	//big.Int可以接受[]byte, string等数据类型，将其转换为big.Int类型
	//big.Int可以通过运算方法，Add,Cmp,Abs
}

//提供一个创建方法
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//方式一：写成固定的
	////target系统自动调节出的一个数值，可以进一步推出一个哈希值
	//targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	//bigTmp := big.Int{}
	//bigTmp.SetString(targetStr, 16)
	//pow.target = &bigTmp

	//方式二：自动推导
	//系统难度值是一个数字，可以推导出哈希值
	//1. 0000000000000000000000000000000000000000000000000000000000000001
	tmp := big.NewInt(1)
	//2. 向左移动256位
	//  10001000000000000000000000000000000000000000000000000000000000000
	tmp.Lsh(tmp, 256-difficult)

	////3. 再次向右移动4个16进制位 (1个16进制数字，代表4个二进制位)
	//tmp.Rsh(tmp, 16)
	pow.target = tmp

	return &pow
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	b := pow.block

	tmp := [][]byte{
		[]byte(b.Version),
		b.PrevHash,
		b.MerkleRoot,
		digi2byte(b.TimeStamp),
		digi2byte(b.Bits),
		digi2byte(nonce), //<<<<===== 注意修改，要使用传递进来的nonce
	}

	data := bytes.Join(tmp, []byte(""))

	return data
}

//核心函数，不断改变Nonce，求出满足条件的哈希值
func (pow *ProofOfWork) Run() (int64, []byte) {

	var nonce int64
	var hash [32]byte

	for {
		fmt.Printf("挖矿中:%x\r", hash)
		//拼好的数据
		data := pow.prepareData(nonce)

		//计算哈希值
		hash = sha256.Sum256(data)

		//转换成big.Int
		tmpInt := big.Int{}
		tmpInt.SetBytes(hash[:])

		//比较
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//func (x *Int) Cmp(y *Int) (r int)
		if tmpInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功，当前哈希值为:%x, nonce: %d\n", hash, nonce)
			break
		} else {
			//fmt.Printf("当前哈希为:%x ,nonce:%d\n", hash, nonce)
			nonce++
		}
	}

	return nonce, hash[:]
}

//提供校验区块Nonce有效性函数
func (pow *ProofOfWork) isValid() bool {
	//1. 拼接数据
	data := pow.prepareData(pow.block.Nonce)
	//2. 计算哈希
	currentHash := sha256.Sum256(data)
	tmpInt := big.Int{}
	tmpInt.SetBytes(currentHash[:])

	////3. 比较结果（仅进行一次计算)
	//if tmpInt.Cmp(pow.target) == -1 {
	//	return true
	//}
	//return false

	//或者
	return tmpInt.Cmp(pow.target) == -1
}
