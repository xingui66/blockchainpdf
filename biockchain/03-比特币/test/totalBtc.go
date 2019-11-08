package main

import "fmt"

func main() {
	total := 0.0          //总量
	intervalCount := 21.0 //万
	reward := 50.0 //最初一个区块的奖励

	for reward != 0 {
		//单个区间内挖矿的总量
		amount := intervalCount * reward

		//统计所有的比特币
		total += amount
		reward *= 0.5 //奖励减半
	}

	fmt.Println("total:", total)
}
