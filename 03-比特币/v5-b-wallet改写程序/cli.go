package main

import (
	"fmt"
	"os"
	"strconv"
)

//定义一个CLI结构
type CLI struct {
	bc *BlockChain
}

const Usage = `
	./blockchain addBlock <data>    "区块数据"
	./blockchain print    "打印区块"
	./blockchain getBalance <地址>  "获取某个地址的余额"
	./blockchain send <FROM> <TO> <AMOUNT> <MINER> <DATA> "转账"
	./blockchain createWallet "创建钱包"
	./blockchain listAddress "打印当前钱包内所有的地址"
`

//持续解析命令的方法
func (cli *CLI) Run() {
	fmt.Println("CLI Run called!")

	cmds := os.Args

	if len(cmds) < 2 {
		fmt.Println("输入参数无效，请检查!")
		fmt.Println(Usage)
		return
	}

	//解析命令
	switch cmds[1] {
	case "addBlock":
		fmt.Println("addBlock called!")
		if len(cmds) != 3 {
			fmt.Println("参数无效!")
			fmt.Println(Usage)
			return
		}

		data := cmds[2]
		cli.addBlock(data)
	case "send":
		//./blockchain send <FROM> <TO> <AMOUNT> <MINER> <DATA> "转账"
		fmt.Println("send called!")
		if len(cmds) != 7 {
			fmt.Println("参数无效!")
			fmt.Println(Usage)
			return
		}

		from := cmds[2]
		to := cmds[3]
		amountStr := cmds[4]
		amount, _ := strconv.ParseFloat(amountStr, 64)
		miner := cmds[5]
		data := cmds[6]

		cli.send(from, to, amount, miner, data)

	case "getBalance":
		fmt.Println("getBalance called!")
		if len(cmds) != 3 {
			fmt.Println("参数无效!")
			fmt.Println(Usage)
			return
		}

		address := cmds[2]
		cli.getBalance(address)

	case "createWallet":
		fmt.Println("createWallet called!")
		cli.createWallet()

	case "print":
		fmt.Println("print called!")
		cli.printBlock()
	case "listAddress":
		fmt.Println("listAddress called!")
		cli.listAddress()

	default:
		fmt.Println("不存在的命令:", cmds[1])
		fmt.Println(Usage)
	}
}
