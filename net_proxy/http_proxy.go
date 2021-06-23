package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func main()(){
	var port string = ":10086"
	li,err := net.Listen("tcp",port)

	if err != nil{
		fmt.Println("tcp listen err",err)
		return
	}
	defer li.Close()
	clientData := make([]byte,40960)
	for {
		client, err := li.Accept()

		if err != nil {
			fmt.Println("连接错误", err)
			break
		}
		n, err := client.Read(clientData)

		if err != nil {
			fmt.Println("读取输入流错误", err)
			break
		}
		fmt.Print("client request data:", string(clientData[:n]))

		var method, host, address string
		fmt.Sscanf(string(clientData[:bytes.IndexByte(clientData[:], '\n')]), "%s%s", &method, &host)
		hostPortURL, err := url.Parse(host)
		if err != nil {
			log.Println(err)
			return
		}

		if hostPortURL.Opaque == "443" { //https访问
			address = hostPortURL.Scheme + ":443"
		} else {                                            //http访问
			if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
				address = hostPortURL.Host + ":80"
			} else {
				address = hostPortURL.Host
			}

		}

		//获得了请求的host和port，就开始拨号吧
		server, err := net.Dial("tcp", address)
		if err != nil {
			log.Println(err)
			return
		}

		if method == "CONNECT" {
			fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n")
		} else {
			server.Write(clientData[:n])
		}
		//进行转发
		go io.Copy(server, client)
		go io.Copy(client, server)
	}
}
