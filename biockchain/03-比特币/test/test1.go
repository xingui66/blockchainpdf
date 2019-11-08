package main

type People struct {
	Name string
	Age  int
}

func getInfo() *People {

	//perple在函数内部创建，属于局部变量，但是在getInfo函数之外引用了
	//此时，编译器会将people变量分配到堆上
	return &People{
		Name: "Lily",
		Age:  20,
	}
}

func main() {
	_ = getInfo()
}
