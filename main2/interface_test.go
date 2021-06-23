package main

import (
	"fmt"
	"testing"
)

func TestInterface(t *testing.T){
	var che1 interface{} = Benz{"benz",12}
	var che2 interface{} = BMW{"bmw"}
	var myf interface{} = getMyFunc()
	switch che1.(type) {
	case int:fmt.Println("int")
	case Vertical:fmt.Println("vertical")
	default:fmt.Println("none")
	}
	switch che2.(type) {
	case int:fmt.Println("int")
	case Vertical:fmt.Println("vertical")
	default:fmt.Println("none")
	}

	//函数实现了Vertical接口 所以属于Vertical接口的子类型
	switch myf.(type) {
	case int:fmt.Println("int")
	case Vertical:fmt.Println("vertical")
	default:fmt.Println("none")
	}

	var vers []Vertical = []Vertical{Benz{name:"benz",age:20},BMW{name:"bmw"},getMyFunc()}
	for idx,instance := range vers{
		fmt.Println(idx,instance.run(2))
	}

}

type Vertical interface{
	run(int)([]string)
}

type Benz struct{
	name string
	age int32
}

type myFunc func()()

func (m myFunc)run(int)([]string){
	fmt.Println("func run")
	return nil
}

func getMyFunc()(myFunc){
	return func()(){}
}
func (b Benz)run(int)([]string){
	fmt.Println("benz run")
	return nil
}

type BMW struct{
	name string
}

func (b BMW)run(int)([]string){
	fmt.Println("bmw run")
	return nil
}

//含有没有变量名的形参 无法引用它
func test(int,a string)(){
	fmt.Print()
}