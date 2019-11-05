package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func main() {
	//创建一条椭圆曲线
	curve := elliptic.P256()

	//1. 创建秘钥对
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("生成私钥失败, err:", err)
		return
	}

	data := "hello world"

	hash := sha256.Sum256([]byte(data))

	//2. 使用私钥签名
	//r和s是数据签名
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	fmt.Println("r len:", len(r.Bytes()))
	fmt.Println("s len:", len(s.Bytes()))

	//r与s的长度是相同的，我们将两者拼接到一起，进行传输
	//到对端，从中间分割，还原成big.Int类型即可
	signature := append(r.Bytes(), s.Bytes()...)

	//进行数据传输。。。。

	pubKey := privKey.PublicKey

	//3. 公钥验证签名

	//还原r，s
	var r1, s1 big.Int
	r1.SetBytes(signature[:len(signature)/2])
	s1.SetBytes(signature[len(signature)/2:])

	//func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool {
	//res := ecdsa.Verify(&pubKey, hash[:], &r1, &s1) ///正确的
	res := ecdsa.Verify(&pubKey, hash[:], &r1, &r1) //错误的
	fmt.Println("res :", res)
}
