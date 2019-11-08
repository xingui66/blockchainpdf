package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivKey *ecdsa.PrivateKey
	PubKey  []byte //这不是原生的公钥，而是X, Y两个点的字节流拼成而成的
}

//创建一个秘钥对
func NewWallet() *Wallet {
	//私钥
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("创建秘钥对失败, err:", err)
		return nil
	}
	//公钥
	pubKeyRaw := priKey.PublicKey
	x := pubKeyRaw.X
	y := pubKeyRaw.Y
	pubKey := append(x.Bytes(), y.Bytes()...)

	return &Wallet{
		PrivKey: priKey,
		PubKey:  pubKey,
	}
}

func (w *Wallet) getAddress() string {
	//一、第一次哈希
	firstHash := sha256.Sum256(w.PubKey)
	//第二次哈希
	hasher := ripemd160.New()
	hasher.Write(firstHash[:])
	pubKeyHash := hasher.Sum(nil)

	//二、在前面添加1个字节的版本号
	payload := append([]byte{byte(00)}, pubKeyHash...)

	//三、做两次哈希运算，截取前四个字节，作为checksum，
	f1 := sha256.Sum256(payload)
	second := sha256.Sum256(f1[:])

	//checksum := second[:] //作闭右开
	checksum := second[:4] //作闭右开

	//四、拼接25字节数据
	payload = append(payload, checksum...)

	//五、base58处理，得到地址
	address := base58.Encode(payload)
	return address
}
