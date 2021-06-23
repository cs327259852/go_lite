package main

import "fmt"
/*
没有发生宕机时，调用recover没有任何效果
 */
func handler1()(){
	//没有发生宕机 函数正常执行
	defer func()(){
		err := recover()
		fmt.Println("error of recover:",err)
	}()
	doSth("handler1")
}
func handle2()(result int){
	//定义闭包 在函数发生宕机时设置返回值
	defer func()(){
		err := recover()
		if err != nil {
			//发生宕机 设置返回值
			result = -1
		}
	}()
	//手动触发宕机
	panic("test panic")
	//由于宕机地点在该行代码前面，所以宕机地点后的代码都不会执行，宕机后逆序执行宕机前的defer语句，并退出当前方法（如果执行了recover，继续执行，否则停止运行）
	doSth("handler2")
	return 0
}

func doSth(s string){
	println(s,"do sth...")
}

func PanicRecover2Execute()(){
	handler1()
	result := handle2()
	fmt.Println(result)
}