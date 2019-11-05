package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
)

const walletFileName = "wallet.dat"

//- 定义结构
type WalletManager struct {
	//1. 定义一个map来管理所有的钱包
	//2. key：地址
	//3. value：wallet
	Wallets map[string]*Wallet
}

//- 创建结构
func NewWalletManager() *WalletManager {
	//return &WalletManager{
	//	wallets: make(map[string]*Wallet),
	//}

	var wm WalletManager
	wm.Wallets = make(map[string]*Wallet)

	//加载已经存在钱包，从wallet.dat
	err := wm.loadFromFile()
	if err != nil {
		fmt.Println("loadFromFile err:", err)
		return nil
	}

	return &wm
}

func (wm *WalletManager) createWallet() (string, error) {
	//调用wallet结构的创建方法
	w := NewWallet()

	if w == nil {
		return "", errors.New("创建钱包失败!")
	}

	address := w.getAddress()

	//填充自己的wallets结构
	wm.Wallets[address] = w

	err := wm.saveToFile()
	if err != nil {
		fmt.Println("存储钱包失败,err:", err)
		return "", err
	}

	return address, nil
}

//1. 将wm结构写入到磁盘=》向map中添加数据
//2. 使用gob对wm进行编码后写入文件
func (wm *WalletManager) saveToFile() error {
	var buff bytes.Buffer

	//对interface数据进行注册
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(wm)
	if err != nil {
		fmt.Println("saveToFile encode err:", err)
		return err
	}

	//写入磁盘
	err = ioutil.WriteFile(walletFileName, buff.Bytes(), 0600)
	if err != nil {
		fmt.Println("saveToFile writeFile err:", err)
		return err
	}

	return nil
}

//加载钱包里面的秘钥对
func (wm *WalletManager) loadFromFile() error {
	//1. 判断文件是否存在，如果不存在，则不需要加载
	if !isFileExist(walletFileName) {
		//这个是第一执行创建钱包时会进入的逻辑, 不属于错误
		fmt.Println("钱包不存在，准备创建!")
		return nil
	}

	//2. 文件存在，读取文件
	fmt.Println("钱包存在，准备读取...")
	data, err := ioutil.ReadFile(walletFileName)
	if err != nil {
		fmt.Println("loadFromFile readFile err:", err)
		return err
	}

	//3.gob进行解码
	//注册接口函数
	gob.Register(elliptic.P256())

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(wm)
	if err != nil {
		fmt.Println("loadFromFile Decode err:", err)
		return err
	}

	return nil
}

func (wm *WalletManager) listAddress() (addresses []string) {
	//1. 遍历map，获取所有的key值
	for address, _ := range wm.Wallets {
		//2. 拼装成切片返回
		addresses = append(addresses, address)
	}

	//3. 将地址数组排序后返回
	//默认是升序排列
	//sort.Strings(addresses)

	return
}
