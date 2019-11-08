package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Person struct {
	Age  int
	Name string
}

func main() {
	p1 := Person{
		Age:  19,
		Name: "Lily",
	}

	//编码
	//1. 创建一个编码器
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	//2. 编码
	err := encoder.Encode(&p1)
	if err != nil {
		fmt.Println("编码失败,err:", err)
		return
	}

	fmt.Println("编码后的数据:", buf.Bytes())

	//.....

	//解码
	p2 := Person{}

	//1. 创建一个解码器
	decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))

	//2. 解码
	err = decoder.Decode(&p2)
	if err != nil {
		fmt.Println("解码失败, err:", err)
		return
	}

	fmt.Println("解码后的数据:", p2)
}
