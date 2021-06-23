package main

import (
	"fmt"
	"net"
	"os"
)
func startTCPServer()(){
	listener,err := net.ListenTCP("tcp",&net.TCPAddr{
		IP:[]byte{0,0,0,0},
		Port: 8002,
	})
	if err != nil{
		fmt.Println("监听失败:",err)
		os.Exit(1)
	}

	conn,err := listener.Accept()

	if err != nil{
		fmt.Println("建立连接异常:",err)
		os.Exit(1)
	}
	defer conn.Close()

	var data []byte = make([]byte,1024)
	_,err = conn.Read(data)
	if nil != err{
		fmt.Println("读取数据异常：",err)
	}
	fmt.Println("收到:",string(data))


}
