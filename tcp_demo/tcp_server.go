package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main()(){
	listener,err := net.ListenTCP("tcp",&net.TCPAddr{IP:  net.IPv4(192,168,1,6),Port:9000})
	ch := make(chan int)
	if err != nil{
		fmt.Println("listen 9000 error")
		return
	}
	defer listener.Close()

	conn,err := listener.Accept()

	defer conn.Close()
	if err !=nil{
		fmt.Println("accept client error")
		return
	}else{
		fmt.Println("accept one client")
	}

	go func(c net.Conn)(){
		var b []byte = make([]byte,1024)
		for ;;{
			n,_ := conn.Read(b)
			fmt.Println("对方:"+string(b[0:n]))
			conn.Write(b[0:n])
			if strings.Contains(string(b[0:n]),"exit") {
				ch <- 1
				break
			}
		}
	}(conn)

	flag := <-ch
	if flag == 1{
		os.Exit(0)
	}
}
