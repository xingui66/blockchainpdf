package main

import (
	"fmt"
	"os"
)

func main() {
	cmds := os.Args

	for i, v := range cmds {
		fmt.Println("i:", i, ", v:", v)
	}
}
