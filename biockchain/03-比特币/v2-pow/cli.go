package main

import (
	"fmt"
	"os"
)

//定义一个CLI结构
type CLI struct {
	bc *BlockChain
}

const Usage = `
	./blockchain addBlock <data>    "区块数据"
	./blockchain print    "打印区块"
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
	case "print":
		fmt.Println("print called!")
		cli.printBlock()
	default:
		fmt.Println("不存在的命令:", cmds[1])
		fmt.Println(Usage)
	}
}
