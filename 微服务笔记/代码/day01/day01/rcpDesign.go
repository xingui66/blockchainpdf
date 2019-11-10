package main

import (
	"net/rpc"
	"net"
	"net/rpc/jsonrpc"
)

const  ServiceName = "hello"

//当父类
type BaseService interface {
	CallFunc(int,*int)error
}

func RegisterService(service BaseService){
	rpc.RegisterName(ServiceName,service)
}

//像调用本地函数一样调用远程函数
type RpcClient struct {
	c *rpc.Client
}


func InitClient(addr string)(RpcClient,error){
	//给结构体中的client赋值
	conn,err := net.Dial("tcp",addr)
	if err != nil {
		return RpcClient{},err
	}

	rpcClent := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	return RpcClient{c:rpcClent},nil
}

func (this*RpcClient)CallFunc(req int,resp*int)error{
	return this.c.Call(ServiceName+".CallFunc",req,resp)
}




