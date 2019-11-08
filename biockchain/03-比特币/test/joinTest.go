package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	str := []string{"hello", "world", "!"}
	res := strings.Join(str, "=")
	fmt.Println("res:", res)

	bytes1 := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("!"),
	}

	b1 := bytes.Join(bytes1, []byte(""))
	fmt.Printf("b1 :%s\n", b1)
}
