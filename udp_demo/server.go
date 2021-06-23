package udp_demo

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func sum(s []int,c chan int)(){
	sum := 0
	for _,v := range s{
		sum += v
	}
	c <- sum
	fmt.Println("go routine finished")
}
func test1()(){
	 var s []int= []int{1,2,3,4,5,6}

	 var c chan int = make(chan int)
	 go sum(s[:len(s)/2],c)
	 go sum(s[len(s)/2:],c)
	time.Sleep(1000*1000*2)
	 x,y := <-c,<-c
	time.Sleep(1000*1000*2)
	fmt.Println("handle data finished!")
	 fmt.Println(x+y)
}

func test2()(){
	c1 := make(chan int)
	go func(c chan int)(){
		<- c1
		fmt.Println("go routine end!")
	}(c1)
	c1 <- 1
	//c1 <- 1
	//c1 <- 1

	fmt.Println("udp_demo end!")
	time.Sleep(1000000000*5)

}

func test3()(){
	ch := make(chan int,1)

	ch <- 1
}

func iterateChan()(){
	//declare a chan int variable with 3 buffers
	ch := make(chan int,3)

	//start go routine to send channel some data
	go func()(){
		ch <- 1
		ch <- 2
		ch <- 3
		fmt.Println("channel data send finished!")
		close(ch)
	}()

	//wait send channel finish
	time.Sleep(1000000000*5)

	//go func()(){
		//iterate channel
		for i := range ch{
			fmt.Println(i)
		}
	//}()
	time.Sleep(1000000000*5)

}

func sliceTest()(){
	//declare a slice
	var slice1 []int
	//declare a array with length is 5
	var arr1 [5]int
	fmt.Println(slice1,arr1)

	//use make way to declare a slice
	//5 init length 10 capacity
	slice2 := make([]int,5,10)
	fmt.Println(slice2)

	//init slice.len=cap=3
	slice3 := []int{1,2,3}
	fmt.Println(cap(slice3),len(slice3))

	//to reference to slice3
	slice4 := slice3[2:3]
	fmt.Println(slice4)

	slice5 := slice3[:]
	fmt.Println(slice5)

	slice6 := slice3[:1]
	fmt.Println(slice6)

	slice7 := make([]int,2,3)
	fmt.Println(len(slice7),cap(slice7))
	slice7 = append(slice7, 2,3,4,5,6,7)
	fmt.Println(len(slice7),cap(slice7))
	fmt.Println(slice7)

	var p * []int = &slice7
	fmt.Printf("%p\n",p)

	var a int
	fmt.Println(a)

	var s []int32
	s = append(s,1)
	fmt.Println(s)

	var numbers []int
	numbers = append(numbers,0)
	numbers = append(numbers,1)
	numbers = append(numbers,2,3,4)
	fmt.Println("numbers:",numbers)
	numbers1 := make([]int,len(numbers),cap(numbers))
	copy(numbers1,numbers)
	fmt.Println("numbers1:",numbers)
	fmt.Printf("numbers:%p,numbers1:%p",&numbers,&numbers1)


}

func startUDPServer()(){
	socket,err := net.ListenUDP("udp",&net.UDPAddr{
		IP:   net.IPv4(0,0,0,0),
		Port: 8001,
	})

	if err != nil{
		fmt.Println("监听失败,err:",err)
		os.Exit(1)
	}
	fmt.Println("正在监听...")
	defer socket.Close()
	var clientConnect int32 = 0
	for{
		//receive routine
		go func()(){
			for{
				var data [4096]byte
				n,addr,err := socket.ReadFromUDP(data[:])
				if err != nil{
					continue
				}
				fmt.Printf("来自%v:%v\n",addr,
					string(data[:n]))
				clientConnect++
				//send routine
				if clientConnect == 1{

					go func()(){
						for{
							in := bufio.NewReader(os.Stdin)
							str,_ ,err := in.ReadLine()
							if err != nil{
								continue;
							}
							var sendData []byte = []byte(str)
							_,err =socket.WriteToUDP(sendData[:],addr)
							if nil != err{
								fmt.Print("发送失败:",err)
								continue
							}
						}
					}()
				}
			}
		}()


		time.Sleep(1000000000*30)
	}
}



func main()(){
	startUDPClient()
}
// client -> server ok
//server -> client *