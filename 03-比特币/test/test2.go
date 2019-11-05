package main

func main() {
	//使用new关键字，也不一定分配到堆上
	//由于在main之外没有使用过str这个变量，所有分配到栈上
	str := new(string)
	*str = "hello"

	//fmt.Printf("str: %s", str)
}
