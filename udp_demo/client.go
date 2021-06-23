package udp_demo

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func startUDPClient()(){
	fmt.Println("正在进入聊天室..")
	socket,err := net.DialUDP("udp",nil,&net.UDPAddr{
		IP:net.IPv4(10,2,124,209),
		Port:8001,
	})
	if err != nil{
		fmt.Println("连接udp端口失败:",err)
		return
	}
	defer socket.Close()
	fmt.Println("成功进入聊天室,开始聊天吧(按回车键发送)")
	for{
		//send routine
		go func()(){
			for{
				in := bufio.NewReader(os.Stdin)
				str,_ ,err := in.ReadLine()
				if err != nil{
					continue;
				}
				var sendData []byte = []byte(str)
				_,err = socket.Write(sendData)
				if err != nil{
					fmt.Println("发送数据失败:",err)
					continue
				}
			}
		}()
		//receive routine
		go func()(){
			for{

				data := make([]byte,4096)
				_,_,_ = socket.ReadFromUDP(data)

				if err != nil{
					fmt.Println("接受数据失败:",err)
					continue
				}
				fmt.Println("对方:"+string(data))
			}
		}()


		time.Sleep(1000000000*60)
	}




}
