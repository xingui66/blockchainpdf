package main

import (
	"fmt"
)


//使用客户端链接服务,调用服务端的函数   grpc框架的封装
func main(){
	//链接
	//conn,err := rpc.Dial("tcp",":1234")
	/*conn,err := net.Dial("tcp",":1234")
	if err != nil {
		fmt.Println("建立链接失败",err)
		return
	}
	defer conn.Close()

	//把链接绑定到rpc上,并且以json做序列化
	rpcConn := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	//调用远程服务
	req := 10
	var resp int
	err = rpcConn.Call(ServiceName+".CallFunc",req,&resp)*/

	//初始化客户端
	client,err:= InitClient("127.0.0.1:1234")
	if err != nil {
		fmt.Println("建立链接失败",err)
		return
	}

	//调用方法
	req := 10
	var resp int
	err = client.CallFunc(req,&resp)


	fmt.Println("获取的数据为",resp)
}
