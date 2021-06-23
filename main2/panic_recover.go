package main

import (
	"fmt"
	"runtime"
)

/*
宕机恢复 try/catch机制
recover是一个内置函数，用于将处于恐慌状态的goroutine恢复执行，
由于goroutine恐慌后只能执行defer延迟语句，所以recover函数只能写在defer后
recover返回宕机的错误信息
程序宕机后，执行宕机点前面的defer语句，并传递恐慌到调用者，直到goroutine停止，如果defer语句中有recover函数，从宕机点退出当前函数回到调用者继续执行，不会停止
 */

func PanicRecover()(){
	divideFunc()
	fmt.Println("after panic outside divideFunc")
	ProtectRun(func()(){
		fmt.Println("发生宕机前。。")
		panic(&panicContext{"手动触发宕机"})
		fmt.Println("发生宕机后")
	})

	ProtectRun(func()(){
		fmt.Println("赋值宕机前")
		var a *int
		//此处发生运作时错误，会发生宕机，即使使用defer recover之后恢复goroutine，后面的语句页不会执行，只是会跳到当前函数外继续执行
		*a = 1
		fmt.Println("赋值宕机后")
	})

	fmt.Println("udp_demo end ")
}

func divideFunc()(a int){
	defer  func()(){
		recover()
	}()
	fmt.Println("触发除0宕机")
	divides(1,0)
	println("after panic in divideFunc")
	a = 1
	return
}

func divides(a,b int)(){
	if b == 0{
		//除数为0时停止goroutine
		panic("divide 0")
	}
}

type panicContext struct{
	function string
}

func ProtectRun(entry func()())(){
	defer func()(){
		err := recover()
		switch err.(type){
		case runtime.Error:
			fmt.Println("runtime error:",err)
		default:fmt.Println("error:",err)
		}
	}()
	entry()
}


type Animal interface{
	Speak()()
}
type Cat struct {

}
func (*Cat)Speak()(){
	fmt.Println("cat speak()")
}


type Dog struct{

}

func (*Dog)Speak()(){
	fmt.Println("dog speak()")
}
type Bull struct{

}
func (*Bull)Speak()(){
	fmt.Println("bull Speak()")
}

func main()(){
	 //animals := []Animal{Cat{},Dog{},Bull{}}
	//animals := []Animal{Dog{}, Cat{}, Bull{}}
	//box1 := box{
	//	open:func()(){
	//		fmt.Println("inline open")
	//	},
	//}
	//
	//box1.open()
	//box1.Open()
	//(&box1).Open()
	//(&box1).open()
	//divideFunc()
	//fmt.Print(2)

	//func()(){
	//	defer func()(){ recover()}()
	//	fmt.Print(1)
	//	panic("")
	//	fmt.Print(3)
	//}()
	//fmt.Print(2)

	privateVarFunc := getPrivateVar()
	//此处打印函数中的局部变量地址和函数中打印的局部变量地址一样，说明使用闭包保留了局部变量，没有释放局部变量
	// * 加指针类型变量表示获取指针类型指向内存地址的值
	fmt.Println(*privateVarFunc())
	var f func()(int)
	for{
		var innerVar int = 10
		f = func()(int){
			innerVar++
			return innerVar
		}
		if true{
			break
		}
	}
	fmt.Println(f())
	fmt.Println(f())
	//fmt.Print(innerVar) //若想访问块中的局部变量，可以使用闭包形式阻止局部变量释放
}
 func closure()(func()()){
	 var a int = 10
	 f := func()(){
		 a++
		 fmt.Print(a)
	 }
	 fmt.Print(&f)
	 return f
 }

type box struct{
	open func ()()
}

func (* box)Open()(){
	fmt.Println("receiver open")
}
type myint int
type mymethod func()()

//使用闭包保存局部变量 并访问
 func getPrivateVar()(func()(*string)){
 	var name string = "hello"
 	fmt.Println(&name)
 	return func()(*string){
 		return &name
	}
 }