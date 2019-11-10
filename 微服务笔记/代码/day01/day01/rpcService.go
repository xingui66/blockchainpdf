package main

import (
	"fmt"
	"net/rpc"
	"net"
	"net/rpc/jsonrpc"
)

type HelloWorld struct {}

//要求 函数第一个参数为传入参数,第二个参数为传出参数  引用传递  ,返回值有且只能有一个,而且类型必须是error
func (this*HelloWorld)CallFunc(req int,resp*int)error{
	*resp = req + 1
	return nil
}

func main(){
	//注册服务  本质  在内部维护了一张hash表
	//rpc.RegisterName(ServiceName,new(HelloWorld))
	RegisterService(new(HelloWorld))

	//能够监听链接
	listener,err := net.Listen("tcp","localhost:1234")
	if err != nil {
		fmt.Println("设置监听失败",err)
		return
	}
	defer listener.Close()

	fmt.Println("开启监听...")
	//建立链接
	conn,err := listener.Accept()
	if err != nil {
		fmt.Println("建立链接失败",err)
		return
	}
	defer conn.Close()

	//用rpc链接
	rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
}
