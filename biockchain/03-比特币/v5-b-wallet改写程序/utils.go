package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

//数字转成[]byte
func digi2byte(num int64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, &num)
	if err != nil {
		fmt.Println("binary.Write err:", err)
		return nil
	}

	return buf.Bytes()
}

//判断文件是否存在
func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	//通过err错误码，判断文件是否存在
	if os.IsNotExist(err) {
		//文件不存在
		return false
	}

	return true
}
