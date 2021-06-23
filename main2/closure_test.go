package main

import (
	"fmt"
	"testing"
)

func TestClosure(t *testing.T){
	f := closure2()
	fmt.Println(f()) //11
	fmt.Println(f()) //12

}

func closure2()(func()(int)){
	var a int = 10
	//使用闭包阻止局部变量内存释放
	return func()(int){
		a++
		return a
	}
}

func TestArrOrSlice(t *testing.T){

	var arr = []int{1,2,3,4,5,6,7}
	fmt.Println("arr:",arr)
	var s1 []int = arr[0:1]
	fmt.Println("s1:",s1)
	s1[0] = 9
	fmt.Println("arr:",arr)
	var s2 []int = s1[1:4]
	fmt.Println("s2:",s2)
	//切片或数组指向的是第一个元素的地址
	fmt.Printf("%p,%p,%p",s1,s2,arr)
	
}
