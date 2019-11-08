package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
