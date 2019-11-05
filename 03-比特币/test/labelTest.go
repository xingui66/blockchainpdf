package main

import "fmt"

func main() {
	//nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9,}

LABEL1:
	for {
		fmt.Println("xxxxxxxxx")

		for i := 0; i < 9; i++ {
			if i == 5 {
				//continue LABEL1
				//goto LABEL1
				break LABEL1
			}
			fmt.Println("i:", i)
		}

		fmt.Println("YYYYYYYYY")
	}
}
